# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.0.11](https://github.com/trustyou/kured-alert-silencer/tree/v0.0.11) - 2024-11-05

### Fixed

- silence: Compare silences with same precision

### Testing

- Add Kubernetes 1.30

## [v0.0.10](https://github.com/trustyou/kured-alert-silencer/tree/v0.0.10) - 2024-11-04

### Fixed

- silence: Remove newline from debug messages

### Build

- deps: Update module github.com/aws/smithy-go to v1.22.0
- deps: Update module github.com/prometheus/common to v0.60.0
- deps: Update kubernetes packages to v0.31.2
- deps: Update module github.com/prometheus/common to v0.60.1

## [v0.0.9](https://github.com/trustyou/kured-alert-silencer/tree/v0.0.9) - 2024-10-15

### Added

- Add traces for debugging not created silences

### Fixed

- silence: Detect when silence is expired

## [v0.0.8](https://github.com/trustyou/kured-alert-silencer/tree/v0.0.8) - 2024-10-07

### Added

- Add node name to silencer comment

### Fixed

- Makefile phony

### Build

- deps: Update module github.com/aws/smithy-go to v1.21.0
- deps: Update golang Docker tag to v1.23.2

### Refactor

- Start log messages with lowercase

## [v0.0.7](https://github.com/trustyou/kured-alert-silencer/tree/v0.0.7) - 2024-09-19

### Fixed

- Check silence expiration when checking current silences defined

## [v0.0.6](https://github.com/trustyou/kured-alert-silencer/tree/v0.0.6) - 2024-09-18

### Fixed

- Change alertmanager log messages

### Ci

- Remove v from VERSION in update-changelog

## [v0.0.5](https://github.com/trustyou/kured-alert-silencer/tree/v0.0.5) - 2024-09-17

### Fixed

- Info log messages format
- Keep main process running on watcher ends or error

### Build

- deps: Update module github.com/prometheus/common to v0.58.0
- deps: Update wagoid/commitlint-github-action action to v6.1.2
- deps: Update module github.com/prometheus/common to v0.59.0
- deps: Update module github.com/prometheus/common to v0.59.1
- deps: Update golang Docker tag to v1.23.1
- deps: Update alpine Docker tag to v3.20.3
- deps: Update kubernetes packages to v0.31.1

## [v0.0.4](https://github.com/trustyou/kured-alert-silencer/tree/v0.0.4) - 2024-08-22

### Documentation

- Add prefix v to last CHANGELOG version

### Ci

- Add script for generating version

## [v0.0.3](https://github.com/trustyou/kured-alert-silencer/tree/v0.0.3) - 2024-08-22 [YANKED]

### Fixed

- Use TTL from Kured to update silences

### Documentation

- Add resources and security context to deployment yaml

### Build

- deps: Update alpine Docker tag to v3.20.1
- deps: Update module github.com/prometheus/common to v0.55.0
- deps: Update module github.com/aws/smithy-go to v1.20.3
- deps: Update actions/setup-go action to v4
- deps: Update golang Docker tag to v1.22.5
- deps: Update kubernetes packages to v0.30.3
- deps: Update alpine Docker tag to v3.20.2
- deps: Update golang Docker tag to v1.22.6
- deps: Update module github.com/aws/smithy-go to v1.20.4

### Styling

- Fix go code with go fmt
- Fix YAML spaces on brackets

### Ci

- Fix auto-tag workflow when not tag needed
- Add and fix pre-commit checks to ensure formating
- Fix pre-commit gofmt checker
- Update cliff template to order groups

## [v0.0.2](https://github.com/trustyou/kured-alert-silencer/tree/v0.0.2) - 2024-06-19

### Build

- Fix version in docker binaries

### CI

- Fix re-run action
- Fix retry-workflow file name
- Add auto-tag.yaml github action
- Add manifest in release

### Styling

- Fix yaml indentation

## [v0.0.1](https://github.com/trustyou/kured-alert-silencer/tree/v0.0.1) - 2024-06-18

### Build

- Fix version in docker binaries

### CI

- Fix re-run action
- Fix retry-workflow file name

## [v0.0.0](https://github.com/trustyou/kured-alert-silencer/tree/v0.0.0) - 2024-06-17

### Added

- Add version flag to show binary information

### Build

- Update module github.com/spf13/viper to v1.19.0
- Update helm/kind-action action to v1.10.0
- Update module github.com/go-openapi/strfmt to v0.23.0
- Update golang Docker tag to v1.22.4
- Update kubernetes packages to v0.30.2
- Update module github.com/spf13/cobra to v1.8.1
- Update module github.com/prometheus/common to v0.54.0

### CI

- Add github images and doc
- Add renovate configuration
- Fix tests to trigger once per PR
- Re-run failed jobs if e2e fails
- Adjust timeout for test failures in 20min
