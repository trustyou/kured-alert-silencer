# kured-alert-silencer

`kured-alert-silencer` is a tool designed to automatically silence alerts in Alertmanager during the reboot process of Kubernetes nodes managed by Kured. This ensures that the alerts triggered by node reboots do not create unnecessary noise in your monitoring system. It provides a configurable silence duration for each reboot event, and allows the use of Go templates for dynamic silence matchers.

## Features

- Automatically silences alerts during Kured node reboots
- Configurable silence durations
- Templated silence matchers using `{{ .NodeName }}` and Go templates
- Seamless integration with Kubernetes and Alertmanager

## Installation

### Prerequisites

Before installing `kured-alert-silencer`, ensure that the following components are already installed in your Kubernetes cluster:

- Kured
- Alertmanager

### Installation Steps

1. **Apply RBAC Configuration**:
   Apply the necessary RBAC configuration to grant `kured-alert-silencer` the required permissions.

   ```bash
   kubectl apply -f install/kubernetes/rbac.yaml
   ```

2. **Deploy kured-alert-silencer**:
   Adjust the deployment parameters to match your cluster environment and apply the deployment configuration.

   ```bash
   kubectl apply -f install/kubernetes/deployment.yaml
   ```

## Compatibility

The following table lists the versions of `kured-alert-silencer`, Kured, and Kubernetes that have been tested and verified to work together:

| kured-alert-silencer | Kured  | Kubernetes |
| -------------------- | ------ | ---------- |
| 0.0.0                | 1.15.1 | 1.28, 1.29 |

## Configuration

To view the available configuration parameters and usage instructions, run the following command:

```bash
docker run --rm -i ghcr.io/trustyou/kured-alert-silencer:0.0.0 --help
```

## Future Enhancements

- Implement authentication for Alertmanager

## Contributing

Contributions are welcome! Please open an issue or submit a pull request on GitHub. For major changes, please open an issue first to discuss what you would like to change.
