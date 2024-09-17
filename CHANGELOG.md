# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
