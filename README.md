# kured-alert-silencer

An opinionated way of silencing alerts while Kured reboot k8s nodes.

## Specs

- check daemon set for any changes in lock
- check alerts based on locks
- create missing alerts with default TTL

we can reuse the kured service account

configuration parameters:

- alertmanager URL : for including basic auth or something similar
- alerts TTL: 10min in our case
- alert labels: default instance=nodeName
- support for alertamanager token?
