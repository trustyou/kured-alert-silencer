#!/usr/bin/env bash

# USE KUBECTL_CMD to pass context and/or namespaces.
KUBECTL_CMD="${KUBECTL_CMD:-kubectl}"
SENTINEL_FILE="${SENTINEL_FILE:-/var/run/reboot-required}"

echo "Creating reboot sentinel on all nodes"

# Restart all nodes except last worker used to run alertmanager
for nodename in $("$KUBECTL_CMD" get nodes -o name | head -n-1); do
    docker exec "${nodename/node\//}" hostname
    docker exec "${nodename/node\//}" touch "${SENTINEL_FILE}"
done
