# Amazon Web Services Switch User (`awsu`)

`awsu` is a client-tool for advanced STS session token and role handling.  

It has the following features:

* Support for long-term (IAM) and short-term credentials via internal, cross-account and external roles
* Configuration through environment variables or flags
* Console access for roles
* Native support for [Yubikeys](https://www.yubico.com/) as source for TOTP MFA tokens - currently `awsu` will always try to assume a session token (passing in the MFA from a Yubikey) regardless of the profiles used; it searches for the `mfa_serial` first in the role profile (if available) and then in the source profile
* Optional caching for temporary credentials
* Snap-in replacement for arbitrary tooling with `exec` mode and shell aliases

## Installation

Install it via [kreuzwerker/homebrew-taps](https://github.com/kreuzwerker/homebrew-taps) or from the release tab on Github.

```
brew install kreuzwerker/taps/awsu
```

### Requirements

* `ykman` for interacting with Yubikeys (https://github.com/Yubico/yubikey-manager)

## Configuration

Working with `awsu` requires a configured [shared credentials file](https://aws.amazon.com/blogs/security/a-new-and-standardized-way-to-manage-credentials-in-the-aws-sdks/).

With an existing file in place `awsu` can be configured to use a certain profile using the `-p` flag or environment variables.

* Profiles can be choosen with `-p` or `AWS_PROFILE` (SDK standard)
* Shared credentials location can be overridden from it's default (e.g. `~/.aws/credentials`) using `AWS_SHARED_CREDENTIALS_FILE` (SDK standard)
* Caching can be disabled with `-n` or `AWSU_NO_CACHE`
* Cache TTLs and the session length for session tokens and roles can be choosen with `AWSU_CACHE_SESSION_TOKEN_TTL` and `AWSU_CACHE_ROLE_TTL` - after half of that TTL has expired the cached files are considered invalid and will be refreshed; this is done to avoid issues during long running operations
* Logging can be enabled with `-v` or `AWSU_VERBOSE`

## Running `awsu`

`awsu` main mode can be run in two flavors: export and exec.

### Export

When `awsu` is invoked without additional arguments, the resulting credentials are exposed as shell exports. In this mode, `awsu` can be used with `eval` to actually set these variables like this:

```
eval $(awsu)
```

After using export mode, credentials can used until they expire.

### Exec

When `awsu` is invoked with a doubledash (`--`) as the last argument it will execute the application specified after the doubledash (including all arguments) with the resulting credentials being set into the environment of this application.

In exec mode it makes sense to use shell aliases to drive `awsu` like e.g. in zsh:

```
alias aws="awsu -v -- aws"
alias terraform="awsu -v -- terraform"
# etc ...
```

## Registering Yubikeys

1. Insert a Yubikey device
* Invoke `awsu register :iam-username`
* This will
  * register a new virtual MFA device that is named after the `:iam-username`
  * store the secret key of this device with an issuer / name combination derived from the virtual MFA's serial number (ARN) that is compatible with Yubikey
  * associate the virtual MFA device with the given `:iam-username`

Afterwards successful registration `awsu` will

1. print the serial number of the key to `stdout` - when entering this as the value to an `mfa_serial` key in an `~/.aws/config` profile, it will be picked up by `awsu` (and most AWS SDK using tools)
* encode the QR code to `stderr` - this code can be scanned for usage e.g. with the AWS Console

Note: the QR code has a slightly non-standard key-uri-format: `otpauth://totp/:username@:profile?secret=:secret&issuer=Amazon`. This makes certain authenticator apps understand which icon to pick and matches the IAM username directly to the used profile.

### Configure MFA in the shared credentials file

After the setting up the MFA device you can configure `aws-cli` to use MFA in the following way:

```
[my_iam_user]
aws_access_key_id = AKIABLAHBLAHBLAHBLAH
aws_secret_access_key = <blah>
region = us-east-1

[my_admin_role]
role_arn = arn:aws:iam::123456789123:role/my_admin_role
source_profile = my_iam_user
mfa_serial = arn:aws:iam::123456789123:mfa/my_iam_user
region = us-east-1
```
