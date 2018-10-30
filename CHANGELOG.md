# Changelog

All notable [changes](http://keepachangelog.com/en/1.0.0/) to this project will be documented in this file.

## [2.3.1] (unreleased)

### Added

- Added support for command to generate OTP token on Yubikey and return it to stdout (#32)

### Changed

- Fixed wrong default for config file (#31)

## [2.3.0]

### Added

- Added environment variables and config flags for all configuration mechanisms
- Added support for long-term credential console link generation

### Changed

- Fixed expires output bug and added expiry hints
- Added better looking error handling
- Massive update to the `README`

#### Internal

- Removed cache and session ttl mechanism and replaced it with duration and grace
- Created new client structure with explicit generators, sources and target  
- Moved logic for console out of the command and into a dedicated helper

### Removed

- Removed "list" command

## [2.2.1]

### Added

- The MFA serial can now additionally be specified using `-m`, `--mfa-serial` or `AWSU_MFA_SERIAL`
- A new generator (next to the default `yubikey`) called `manual` has been exposed using `-g`, `--generator` or `AWSU_TOKEN_GENERATOR` - this can be used to manually enter tokens for scenarios where roles are used in contexts where IAM policy [conditions](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_mfa_sample-policies.html#ExampleMFADenyNotRecent) don't prevent the usage of tokens that are older than e.g. 1 hour

### Changed

- Fixed #26 and #27 (Limit role credential duration to AWS default of 1 hour)

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
