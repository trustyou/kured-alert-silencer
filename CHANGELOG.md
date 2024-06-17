# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.0.0](https://github.com/trustyou/kured-alert-silencer/tree/v0.0.0) - 2024-06-17

### Added

* Add version flag to show binary information

### Build

* Update module github.com/spf13/viper to v1.19.0
* Update helm/kind-action action to v1.10.0
* Update module github.com/go-openapi/strfmt to v0.23.0
* Update golang Docker tag to v1.22.4
* Update kubernetes packages to v0.30.2
* Update module github.com/spf13/cobra to v1.8.1
* Update module github.com/prometheus/common to v0.54.0

### CI

* Add github images and doc
* Add renovate configuration
* Fix tests to trigger once per PR
* Re-run failed jobs if e2e fails
* Adjust timeout for test failures in 20min
