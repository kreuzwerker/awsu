# Changelog

All notable [changes](http://keepachangelog.com/en/1.0.0/) to this project will be documented in this file.

## [2.0.1]

### Added

- Added config-file-less mode using environment variables (e.g. in case Terraform is not used)

### Changed

- Use username for virtual device name (instead of random id) - this should make self-service policies possible again
- Added missing `export` prefix to export mode
