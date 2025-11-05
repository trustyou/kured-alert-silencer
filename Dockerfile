FROM golang:1.25.4 AS builder

ARG VERSION

COPY . .
RUN VERSION=${VERSION} make kured-alert-silencer

FROM alpine:3.22.2

RUN apk update --no-cache \
  && apk upgrade --no-cache \
  && apk add --no-cache \
    ca-certificates \
    tzdata
COPY --from=builder /go/dist/kured-alert-silencer /usr/bin/kured-alert-silencer

ENTRYPOINT ["/usr/bin/kured-alert-silencer"]
