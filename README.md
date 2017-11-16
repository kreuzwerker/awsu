# Amazon Web Services Switch User (`awsu`)

`awsu` is a binary that maps application environments to AWS profiles (e.g. from `~/.aws/config`).

It has the following additional features:

* Support for long-term (IAM) and short-term credentials via internal, cross-account and external roles
* Automapping of environments to / Terraform [Terraform](https://www.terraform.io/) [workspaces](https://www.terraform.io/docs/state/workspaces.html)
* Configuration through environment variables or project-specific configuration files
* Console access for external and federated roles
* Native support for [Yubikeys](https://www.yubico.com/) as source for TOTP MFA tokens (only when assuming roles)
* Snap-in replacement for arbitrary tooling with `exec` mode and shell aliases

## Configuring `awsu`

You can choose between two configuration modes with `awsu`: a configuration file mode or an environment variable mode. Both modes map `awsu` workspaces to AWS profiles. The currently used workspace is determined by

1. the `--workspace` flag passed to `awsu` or
2. the `TF_WORKSPACE` environment variable or
3. the `AWSU_WORKSPACE` environment variable or
4. the current Terraform workspace
5. `default`

### Configuration via environment

This mode works by setting `AWSU_PROFILE_:WORKSPACE` environment variables to values that correspond to an AWS profile.

Example:

```
export AWSU_PROFILE_DEFAULT=sandbox
export AWSU_PROFILE_PRODUCTION=live
```

### Configuration via configuration file

This mode works by creating a `.awsu` INI file config in each project directory where you want to call `awsu`. Inside this file

1. add a `[profiles]` section
2. add every workspace as key and the matching AWS profile as value

Example:

```
[profiles]

default     = sandbox
production  = production
```

## Running `awsu`

`awsu` can be run in two modes: export and exec.

### Export mode

When `awsu` is invoked without additional arguments, the resulting credentials are exposed as shell exports. In this mode, `awsu` can be used with `eval` to actually set these variables like this:

```
eval $(awsu)
```

After using export mode, credentials can used until they expire. `awsu` always tries to assume roles for the minimum amount of time (15 minutes).

### Exec mode

When `awsu` is invoked with a doubledash (`--`) as the last argument it will execute the application specified after the doubledash (including all arguments) with the resulting credentials being set into the environment of this application.

In exec mode it makes sense to use shell aliases to drive `awsu` like e.g. in zsh:

```
alias aws="awsu -v -- aws"
alias terraform="awsu -v -- terraform"
# etc ...
```

## Caching

The resulting session metadata (regardless of the mode) is cached in `~/.awsu/sessions/:id/:environment.json`. The `:id` is the result of a SHA1 hash over the absolute path to the `.awsu` file. The whole path is logged to `stderr` when invoking `awsu` with the `verbose` parameter.

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

### Configure MFA as part of aws-cli

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
