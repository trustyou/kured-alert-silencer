#!/bin/bash

# Last worker is used to run alertmanager and should not be rebooted
NODECOUNT=${NODECOUNT:-4}
KUBECTL_CMD="${KUBECTL_CMD:-kubectl}"

# Define the command to be run in the Kubernetes pod, split across multiple lines for readability
COMMAND=$(cat <<EOF
apt update > /dev/null 2>&1 &&
apt install -y curl jq > /dev/null 2>&1 &&
curl -s http://alertmanager.default:9093/api/v2/silences |
jq '[.[] | .matchers[].value] | unique | length'
EOF
)

echo "Checking the number of silences in Alertmanager"
output=$("$KUBECTL_CMD" run --rm -i --quiet --image debian:10.9-slim check-silences -- bash -c "$COMMAND")

output_trimmed=$(echo "$output" | xargs)

if [[ "$output_trimmed" == "$NODECOUNT" ]]; then
    echo "Success: The number of silences is $NODECOUNT."
else
    echo "Failure: The number of silences is not $NODECOUNT. It is $output_trimmed."
    exit 1
fi
