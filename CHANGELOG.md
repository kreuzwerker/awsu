# Changelog

All notable [changes](http://keepachangelog.com/en/1.0.0/) to this project will be documented in this file.

## [2.2.0]

### Changed

- Fixed #23 (Increase role and session token duration)
- Fixed bug with missing grace period in session token

## [2.1.1]

### Changed

- Fixed #21 (Don’t get session token w/o MFA)

## [2.1.0]

### Changed

- Abandonded SDK internal logic for assuming roles with tokens and read shared configs directly instead
- Dropped workspace support for the time being - directly select profiles with `-p` or `AWS_PROFILE` instead
- Always get session token (including MFA) before doing anything else - this allows assuming the role in another tool e.g. Terraform while still having a MFA in the mix
- Changed cache location of the sessions to just use the name of the profile

### Added

- Added `list` command to show all configured profiles
- Added `no-cache` option to prevent caching

## [2.0.2]

### Changed

- Put instructions on how to delete the MFA into registration error message

## [2.0.1]

### Added

- Added config-file-less mode using environment variables (e.g. in case Terraform is not used)
- Trigger verbose mode via `AWSU_VERBOSE`

### Changed

- Use username for virtual device name (instead of random id) - this should make self-service policies possible again
- Added missing `export` prefix to export mode
- Log code when fetching from Yubikey
- Always assume roles for 1h
