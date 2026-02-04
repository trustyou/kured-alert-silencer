FROM golang:1.25.7 AS builder

ARG VERSION

COPY . .
RUN VERSION=${VERSION} make kured-alert-silencer

FROM alpine:3.23.3

RUN apk update --no-cache \
  && apk upgrade --no-cache \
  && apk add --no-cache \
    ca-certificates \
    tzdata
COPY --from=builder /go/dist/kured-alert-silencer /usr/bin/kured-alert-silencer

ENTRYPOINT ["/usr/bin/kured-alert-silencer"]
