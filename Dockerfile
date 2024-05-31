FROM alpine:3.20.0

RUN apk update --no-cache \
  && apk upgrade --no-cache \
  && apk add --no-cache \
    ca-certificates \
    tzdata

COPY ./dist/kured-alert-silencer /usr/bin/kured-alert-silencer

ENTRYPOINT ["/usr/bin/kured-alert-silencer"]
