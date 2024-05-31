.DEFAULT: all
.PHONY: all clean image minikube-publish manifest test kured-alert-silencer-all

TEMPDIR=./.tmp
GORELEASER_CMD=$(TEMPDIR)/goreleaser
DH_ORG=trustyou
VERSION=$(shell git rev-parse --short HEAD)
SUDO=$(shell docker info >/dev/null 2>&1 || echo "sudo -E")
KUBERNETES_VERSION=1.28
KIND_CLUSTER_NAME=chart-testing

all: image

clean:
	rm -rf ./dist

kured-alert-silencer:
	CGO_ENABLED=0 go build -o dist/kured-alert-silencer cmd/kured-alert-silencer/main.go

image: kured-alert-silencer
	$(SUDO) docker buildx build --load -t ghcr.io/$(DH_ORG)/kured-alert-silencer:$(VERSION) .

manifest:
	sed -i "s#image: ghcr.io/.*kured-alert-silencer.*#image: ghcr.io/$(DH_ORG)/kured-alert-silencer:$(VERSION)#g" \
		install/kubernetes/deployment.yaml
	sed -i 's|#\(.*\)--silence-duration=10m|\1--silence-duration=1h|g' \
		install/kubernetes/deployment.yaml
	sed -i 's|#\(.*\)--alertmanager-url=http://localhost:9093|\1--alertmanager-url=http://alertmanager.default:9093|g' \
		install/kubernetes/deployment.yaml

lint:
	echo "Running linter on pkg"
	go vet ./pkg/...
	echo "Running linter on cmd"
	go vet ./cmd/...

test: lint
	echo "Running go tests"
	go test ./...

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

delete-kind:
	kind delete cluster --name $(KIND_CLUSTER_NAME)
