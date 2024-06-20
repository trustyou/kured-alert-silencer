.DEFAULT: all

TEMPDIR = ./.tmp
DH_ORG = trustyou
VERSION ?= $(shell git rev-parse --short HEAD)
SUDO = $(shell docker info >/dev/null 2>&1 || echo "sudo -E")
KUBERNETES_VERSION = 1.28
KIND_CLUSTER_NAME = chart-testing

.PHONY: all
all: image

.PHONY: clean
clean:
	rm -rf ./dist

.PHONY: kured-alert-silencer
kured-alert-silencer:
	CGO_ENABLED=0 go build -ldflags \
		"-X github.com/prometheus/common/version.Version=$(VERSION) \
		-X github.com/prometheus/common/version.Branch=$(shell git rev-parse --abbrev-ref HEAD) \
		-X github.com/prometheus/common/version.Revision=$(shell git rev-parse --short HEAD) \
		-X github.com/prometheus/common/version.BuildUser=$(shell whoami)@$(shell hostname) \
		-X github.com/prometheus/common/version.BuildDate=$(shell date +%Y-%m-%dT%H:%M:%SZ)" \
		-o dist/kured-alert-silencer cmd/kured-alert-silencer/main.go

.PHONY: image
image:
	$(SUDO) docker buildx build $(DOCKER_EXTRA_ARGS) \
		--build-arg VERSION=$(VERSION) \
		--load -t ghcr.io/$(DH_ORG)/kured-alert-silencer:$(VERSION) .

.PHONY: push-images
push-images: DOCKER_EXTRA_ARGS ?= --platform linux/amd64,linux/arm64
push-images:
	$(SUDO) docker buildx build $(DOCKER_EXTRA_ARGS) \
		--build-arg VERSION=$(VERSION) \
		--push -t ghcr.io/$(DH_ORG)/kured-alert-silencer:$(VERSION) .

.PHONY: manifest
manifest:
	sed -i "s#image: ghcr.io/.*kured-alert-silencer.*#image: ghcr.io/$(DH_ORG)/kured-alert-silencer:$(VERSION)#g" \
		install/kubernetes/deployment.yaml
	sed -i 's|#\(.*\)--silence-duration=10m|\1--silence-duration=1h|g' \
		install/kubernetes/deployment.yaml
	sed -i 's|#\(.*\)--alertmanager-url=http://localhost:9093|\1--alertmanager-url=http://alertmanager.default:9093|g' \
		install/kubernetes/deployment.yaml

.PHONY: format
format:
	echo "Format source code"
	go fmt ./...

.PHONY: lint
lint: format
	echo "Running linter on pkg"
	go vet ./pkg/...
	echo "Running linter on cmd"
	go vet ./cmd/...

.PHONY: test
test: lint
	echo "Running go tests"
	go test ./...

.PHONY: kind
e2e: image manifest
	kind create cluster --name $(KIND_CLUSTER_NAME) --config .github/kind-cluster-$(KUBERNETES_VERSION).yaml
	helm repo add kubereboot https://kubereboot.github.io/charts
	kubectl apply -f test/alertmanager.yaml
	helm install -n kube-system kured kubereboot/kured \
		--set configuration.period=10s \
		--set configuration.rebootCommand="rm /var/run/reboot-required && pkill containerd"
	kind load --name $(KIND_CLUSTER_NAME) docker-image ghcr.io/$(DH_ORG)/kured-alert-silencer:$(VERSION)
	for i in {1..20}; do \
		if kubectl get ds -n kube-system kured | grep -E 'kured.*5.*5.*5.*5.*5'; then \
			echo "Kured daemonset is ready"; \
			break; \
		else \
			echo "Retrying in 10 seconds..."; \
			sleep 10; \
		fi \
	done
	kubectl apply -f install/kubernetes/rbac.yaml
	kubectl apply -f install/kubernetes/deployment.yaml
	./test/kind/create-reboot-sentinels.sh
	./test/kind/follow-coordinated-reboot.sh
	./test/kind/check-silences.sh

.PHONY: delete-kind
delete-kind:
	kind delete cluster --name $(KIND_CLUSTER_NAME)

.PHONY: update-changelog
update-changelog:
	git cliff -t $(VERSION) -u -p CHANGELOG.md
