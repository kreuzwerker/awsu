# Amazon Web Services Switch User (`awsu`)

`awsu` is a binary that uses dotfiles (`.awsu`) in project directories to reference AWS credentials (e.g. in `~/.aws/config`). It supports assuming cross-account and external roles, auto-mapping of profiles to [Terraform](https://www.terraform.io/) [workspaces](https://www.terraform.io/docs/state/workspaces.html), console access for federation tokens and [Yubikeys](https://www.yubico.com/) as source for TOTP MFA tokens (only when assuming roles).

## Workflow with `awsu`

1. add a new `.awsu` file into a project directory
* add a `[profiles]` section and add every environment / Terraform workspace as key / value pair; the current environment is either
  * deducted from the current Terraform workspace
  * overridden with the `TF_WORKSPACE` environment variable
  * set with the `--workspace` parameter when using `awsu` directly
  * of `default` if not specified otherwise

The resulting session metadata is cached in `~/.awsu/sessions/:id/:environment.json`. The `:id` is the result of a SHA1 hash over the absolute path to the `.awsu` file. The whole path is logged to `stderr` when invoking `awsu` with the `verbose` parameter.

`awsu` can now be used in two modes: export and exec.

### Export mode

When `awsu` is invoked without additional arguments, the resulting credentials are exposed as shell exports. In this mode, `awsu` can be used with `eval` to actually set these variables like this:

```
eval $(awsu)
```

After using export mode, credentials can used until they expire. `awsu` always tries to assume roles for the minimum amounf of time (15 minutes).

### Exec mode

When `awsu` is invoked with a doubledash (`--`) as the last argument it will execute the application specified after the doubledash (including all arguments) with the resulting credentials being set into the environment of this application.

In exec mode it makes sense to use shell aliases to drive `awsu` like e.g. in zsh:

```
alias aws="awsu -v -- aws"
alias terraform="awsu -v -- terraform"
# etc ...
```

### `.awsu` example

```
[profiles]

default     = sandbox
production  = production
staging     = staging
```

## Registering Yubikeys

1. Insert a Yubikey device
* Invoke `awsu register :iam-username`
* This will
  * register a new virtual MFA device with a name `awsu-:id` with `:id` being a 16-byte random ID
  * store the secret key of this device with an issuer / name combination derived from the virtual MFA's serial number (ARN) that is compatible with Yubikey
  * associate the virtual MFA device with the given `:iam-username`

Afterwards successful registration `awsu` will

1. print the serial number of the key to `stdout` - when entering this as the value to an `mfa_serial` key in an `~/.aws/config` profile, it will be picked up by `awsu` (and most AWS SDK using tools)
* encode the QR code to `stderr` - this code can be scanned for usage e.g. with the AWS Console

Note: the QR code has a slightly non-standard key-uri-format: `otpauth://totp/:username@:profile?secret=:secret&issuer=Amazon`. This makes certain authenticator apps understand which icon to pick and matches the IAM username directly to the used profile.
