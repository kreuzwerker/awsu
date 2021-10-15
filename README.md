# Amazon Web Services Switch User (`awsu`)

[![Release](https://img.shields.io/github/v/release/kreuzwerker/awsu)](https://github.com/kreuzwerker/awsu/releases)
[![Build Status](https://github.com/kreuzwerker/awsu/workflows/build/badge.svg)](https://github.com/kreuzwerker/awsu/actions)
[![Documentation](https://godoc.org/github.com/kreuzwerker/awsu?status.svg)](http://godoc.org/github.com/kreuzwerker/awsu) 
[![Go Report Card](https://goreportcard.com/badge/github.com/kreuzwerker/awsu)](https://goreportcard.com/report/github.com/kreuzwerker/awsu) 

`awsu` provides a convenient integration of [AWS](https://aws.amazon.com/) virtual [MFA devices](https://aws.amazon.com/iam/details/mfa/) into commandline based workflows. It does use [Yubikeys](https://www.yubico.com/) to provide the underlying [TOTP](https://tools.ietf.org/html/rfc6238) one-time passwords but does not rely on additional external infrastructure such as e.g. federation.

There is also a high-level video overview from [This Is My Architecture](https://amzn.to/2Tpiv1m) Munich:

[![Video overview](https://img.youtube.com/vi/4FUqak5E_CA/0.jpg)](https://www.youtube.com/watch?v=4FUqak5E_CA)

[ [Installation](#installation) | [Usage](#usage) | [Configuration](#configuration) | [Caching](#caching) | [Commands](#commands) | [General multifactor considerations](#general-multifactor-considerations) ]

# Installation

Production-ready Mac releases can be installed e.g.through `brew` via [kreuzwerker/homebrew-taps](https://github.com/kreuzwerker/homebrew-taps):

```
brew tap kreuzwerker/taps && brew install kreuzwerker/taps/awsu
```

Linux is only available for download from the release tab. No Windows builds are provided at the moment.

## Prequisites

`awsu` relies on [shared credentials files](https://aws.amazon.com/blogs/security/a-new-and-standardized-way-to-manage-credentials-in-the-aws-sdks/) (the same configuration files that other tools such as e.g. the [AWS commandline utilities](https://aws.amazon.com/cli/) are also using) being configured. The profiles are used to determine

1. which IAM long-term credentials ([access key pairs](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html#access-keys-and-secret-access-keys)) are going to be used
2. if a / which virtual MFA device is going to be used
3. if a / which [IAM role](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles.html) is going be used

In contrast to the official AWS CLI `awsu` also supports putting an `mfa_serial` key into a profile which contains long-term credentials (instead of a role). In this case a virtual MFA device is _always_ used when using the long-term credential profile in question.

# Usage

An abstract overview of the usage workflow of `awsu` looks like this:

1. You ensure you fulfill the prerequisites above
2. You start using `awsu` by using the `register` command: this will create a virtual MFA device on AWS, store the secret of that device on your Yubikey and enable the virtual MFA device by getting two valid one-time passwords from the Yubikey
3. Now you just run `awsu` and - after a double-dash `—` - you specify the program you want to run in a given profile (e.g. `admin`); depending on the profile `awsu` will
   1. determine the access-key pair to use
   2. optionally use TOTP tokens from your Yubikey device (to get a session token from AWS in order to add the MFA context to the request) and
   3. optionally assume an IAM role
   4. use the credentials resulting from these operation(s), cache them and export them to the environment of the program specified after the double-dash

Given the following shared credentials config:

```  
[default]
aws_access_key_id     = AKIAIOSFODNN7EXAMPLE
aws_secret_access_key = wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY

[foo]
aws_access_key_id     = AKIAIFODNNOS7EXAMPLE
aws_secret_access_key = bPxRfiCYEXAMPLEKEY/K7MDENG/wJalrXUtnFEMI
mfa_serial            = arn:aws:iam::123456789012:mfa/user@example.com

[bar]
mfa_serial            = arn:aws:iam::123456789012:mfa/user@example.com
role_arn              = arn:aws:iam::123456789012:role/foo-cross-account
source_profile        = default

[wee]
external_id           = 1a03197b-3cb5-491b-bc06-84795afc95ef
mfa_serial            = arn:aws:iam::123456789012:mfa/user@example.com
role_arn              = arn:aws:iam::121234567890:role/bar-cross-account
source_profile        = default

[gee]
role_arn              = arn:aws:iam::121234567890:role/wee-cross-account
source_profile        = foo
```

`awsu` will produce the following results:

| Profile   | Credentials                                          | Cached? | MFA?                     |
| --------- | ---------------------------------------------------- | ------- | ------------------------ |
| `default` | Long-term*                                           | No*     | No*                      |
| `foo`     | Short-term session-token                             | Yes     | Yes, from "foo" itself** |
| `bar`     | Short-term session-token, then role                  | Yes     | Yes, from "bar" itself   |
| `wee`     | Short-term session-token, then role with external ID | Yes     | Yes, from "wee" itself   |
| `gee`     | Short-term session-token, then role                  | Yes     | Yes, from "foo"          |

*) unless a MFA is specified as parameter to `awsu` - then a short-term session-token (equivalent to `foo`) is produced

**) the form of using `mfa_serial` directly long-term credential profiles is not supported by the official AWS CLI (it will just ignore it) - please note that this form will _always_ produce short-term credentials which may not be useful in some circumstances e.g. when re-registering a previously unregistered virtual MFA device

## Configuration

In the following the global configuration parameters are described. Note that all parameters are implemented as flags with environment variable equivalents - when these start with `AWS_` (like the setting of profiles) they have equivalent semantics to e.g. other SDK using applications such as the AWS CLI.

### Profile and shared credential files

These options describe the currently selected profile and the locations of the shared credential files.

|                             | Long                      | Short | Environment                   | Default   |
| --------------------------- | ------------------------- | ----- | ----------------------------- | --------- |
| Currently used profile      | `profile`                 | `p`   | `AWS_PROFILE`                 | `default` |
| Shared config file location | `config-file`             | `c`   | `AWS_CONFIG_FILE`             | Platform  |
| Shared credentials file     | `shared-credentials-file` | `s`   | `AWS_SHARED_CREDENTIALS_FILE` | Platform  |

### Short-term credentials

These options describe caching options, session token and role durations and other aspects of short-term credentials.

| Description                                                  | Long        | Short | Environment      | Default                                                      |
| ------------------------------------------------------------ | ----------- | ----- | ---------------- | ------------------------------------------------------------ |
| Disable caching                                              | `no-cache`  | `n`   | `AWSU_NO_CACHE`  | `false`                                                      |
| Duration of session tokens & roles                           | `duration`  | `d`   | `AWSU_DURATION`  | 1 hour, maximum depends on config of the role in question (up to 12 hours) |
| Grace period until caches expire - in other words: the time a session will be guaranteed to be valid | `grace`     | `r`   | `AWSU_GRACE`     | 45 minutes                                                   |
| Source of OTP tokens                                         | `generator` | `g`   | `AWSU_GENERATOR` | `yubikey` - can be set to `manual` if you want to manually enter OTP passwords |
| MFA serial override                                          | `mfa`       | `m`   | `AWSU_SERIAL`    | None - can be used to set or override MFA serials            |

### Other

| Description     | Long      | Short | Environment    | Default |
| --------------- | --------- | ----- | -------------- | ------- |
| Verbose logging | `verbose` | `v`   | `AWSU_VERBOSE` | `false` |

## Caching

Caching is only used for cacheable (short-term) credentials. It can be disabled completely (even on a case-by-case basis) and is controlled by two primary factors: duration and grace.

The duration is always equivalent to the duration of the session token used which - in turn - is equivalent to the duration of the optionally assumed role (1 hour by default).

The grace period is the minimum safe distance to the duration before it's considered invalid (45 minutes by default). This is useful for dealing with long-running deployments that might be interrupted when e.g. a role becomes invalid mid-deployment.

Example: in the default setting with a duration of 1 hour and a grace period of 45 minutes `awsu` will consider cached credentials invalid after 15 minutes.

## Commands

### Default

The default command (invoking just `awsu`) has two modes: _export_ and _exec_. It supports just the global parameters described above.

#### Export mode

When `awsu` is invoked without additional arguments, the resulting credentials are exposed as shell exports. In this mode, `awsu` can be used with `eval` to actually set these variables like this:

```
eval $(awsu)
```

After using export mode, credentials can used until they expire.

#### Exec mode

When `awsu` is invoked with a doubledash (`--`) as the last argument it will execute the application specified after the doubledash (including all arguments) with the resulting credentials being set into the environment of this application.

In exec mode it makes sense to use shell aliases to drive `awsu` like e.g. in zsh:

```
alias aws="awsu -v -- aws"
alias terraform="awsu -v -- terraform"
# etc ...
```

Note that when using this alias style:

1. you can always reference the alias targets with absolute paths to temporarily escape, e.g. by referring to `aws` as e.g. `/usr/local/bin/aws`
2. you can still configure `awsu` by using the environment variable style of parameter passing

### Register

`awsu register :iam-username` will perform the following actions:

1. create a new virtual MFA device that is named after the `:iam-username`
2. store the secret key of this device with a name derived from the virtual MFA's serial number (ARN) that is compatible with Yubikey
3. enable the virtual MFA device with the given `:iam-username`

After successful registration `awsu` will log the serial of the MFA for usage as `mfa_serial` in your profiles.

#### Parameters

The following parameters are exclusive to `register`.

| Description                                                  | Long     | Short | Environment   | Default  |
| ------------------------------------------------------------ | -------- | ----- | ------------- | -------- |
| Generates a QR code of the MFA secret that can be used for backup purposes on smartphones | `qr`     | `q`   | `AWSU_QR`     | `true`   |
| Sets the "issuer" part of the QR code URL - depending on the smartphone app used this may add stylistic information to the one-time passwords (e.g. an icon) | `issuer` | `i`   | `AWSU_ISSUER` | `Amazon` |

### Unregister

`awsu unregister :iam-username` will perform the following actions:

1. deactivate the virtual MFA device that is named after the `:iam-username`
2. remove the matching TOTP secret from the Yubikey
3. delete the virtual MFA device that is named after the `:iam-username`

### Console

`awsu console` will open (in a browser) the link to the AWS console for a profile. It supports:

1. long-term credentials
2. short-term credentials for
   1. cross-account ("internal") roles
   2. external ("federated") roles (role with an `external_id` field in their profile)

#### Parameters

The following parameters are exclusive to `console`.

| Description                                                  | Long   | Short | Environment | Default |
| ------------------------------------------------------------ | ------ | ----- | ----------- | ------- |
| Opens the resulting link in a browser (as opposed to just returning it) | `open` | `o`   | `-`         | `true`  |

### Token

`awsu token` will generate a fresh one-time password from an inserted Yubikey. In order to determine the correct secret it will

1. use the MFA directly configured on the currently used `profile` or
2. use the the MFA configured through the `mfa-serial` parameter.

# General multifactor considerations

The goal of this section is to consider the multifactor options that are available on AWS without involving additional external infrastructure (e.g. by utilizing [federation](https://aws.amazon.com/identity/federation/)). Under these constraints there are two options available:

1. Restricing access by IPv4 addresses and/or networks, expressed in [CIDR](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing) notation
   1. Easy to implement in AWS and requires no additional effort except for the control of static IPv4 IPs or ranges
   2. Hard to restrict to the right people in an organization e.g. your operators network vs. your backoffice network or guest wifi network
   3. Some impact on remote workers - you'll need a VPN solution that routes all traffic for [AWS IP addresses](https://docs.aws.amazon.com/general/latest/gr/aws-ip-ranges.html) (or just all traffic)
   4. Difficult to spoof unless you are a very privileged attacker
2. Restricing access by proving the possession of a [virtual or hardware MFA device](https://aws.amazon.com/iam/details/mfa/), specifically a device that can emit 6-digits [TOTP](https://tools.ietf.org/html/rfc6238) one-time passwords (OTP)
   1. More difficult to implement in AWS since you'll want to introduce smartcards as sources for OTPs (unless users want to type in OTPs every _n_ minutes) - `awsu` supports [Yubikey](https://www.yubico.com/products/yubikey-hardware/) USB smartcards here
   2. Easy to restrict in usage to the right people since MFA devices are first-class entities in IAM
   3. No impact on remote workers
   4. Difficult to spoof due to the aggressive rate-limits of MFA authentication attempts - since AWS uses a 6-digit configuration it's nevertheless possible and failures to authenticate with MFA should be monitoried through AWS CloudTrail

When possible we recommend to require both seconds factors.

## Providing second factors

IP based second factors are provided implicitly by virtue of being the address from which a request originates.

MFA based second factors are provided through one of the following mechanisms:

1. Using long-term credentials (the access key pair associated with an IAM user) to get a [session token via STS](https://docs.aws.amazon.com/STS/latest/APIReference/API_GetSessionToken.html)
   1. this yields new short-term credentials (a new access key pair plus the session token itself) that now contains
      1. the information about the MFA context but
      2. no new principal (the IAM user remains the principal)
   2. a session token can be requested to be valid for a period between 15 minutes and 12 hours - there is no way of directly configuring an upper ceiling  
2. Using long- or short-term credentials to [assume a role via STS](https://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRole.html)
   1. this yields new short-term credentials with a new principal (the role)
      1. when using long-term credentials the MFA serial and OTP token must be passed to the STS API when assuming the role
      2. when using short-credentials aquired from the session token mechanism described above the MFA serial and OTP token can be omitted
      3. in both cases the role now includes a MFA context
   2. a role can be requested to be valid for a period between 15 minutes and 12 hours - it's upper ceiling (between 1 hour and 12 hours, with 1 hour being the default) is configured directly on the IAM role

## Requiring second factors

Second factor requirements are expressed in the form of [global conditions keys](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_condition-keys.html) in the the optional `Condition` block that is found in most of the available types of IAM policies such as

1. [Trust policies](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_terms-and-concepts.html) of roles
2. Policies attached to e.g. roles or users
3. [Permission boundaries](https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies_boundaries.html) attached to e.g. roles or users
4. Selected service policies e.g. [S3 bucket policies](https://docs.aws.amazon.com/AmazonS3/latest/dev/example-bucket-policies.html#example-bucket-policies-use-case-3)

[Service control policies](https://docs.aws.amazon.com/organizations/latest/userguide/orgs_reference_scp-syntax.html) that can be applied globally on [AWS Organization](https://docs.aws.amazon.com/organizations/latest/userguide/orgs_introduction.html) members are not supported at the moment.

Regardless of its location the `Condition` for an IP-based second factor would look like this:

```json
"IpAddress": {
  "aws:SourceIp": "203.0.113.0/24"
}
```

The `Condition` content for an MFA-based second factor can have two forms.

The "age" form:

```json
"NumericLessThan": {
  "aws:MultiFactorAuthAge": "3600"
}
```

Or the "boolean" form:

```json
"Bool": {
  "aws:MultiFactorAuthPresent": "true"
}
```

The "age" form is preferred as the means of requiring a MFA. Since a session token (see the section before) cannot directly be configured for a maximum duration the "age" form is the only way to set an acceptable window for the proof of possession of a MFA device.

## Location implications

The type of policy where a second factor requirement is enforced at has some implications.

When applied to a non-trust policy, the requirement is enforced on _every_ interaction. That also means that if a permission is denied (e.g. due to `aws:MultiFactorAuthAge`) the user has not always a way of knowing exactly why it's denied (in some AWS APIs the `UnauthorizedOperation` response contains an `EncodedMessage` that can be decoded using the [appropriate STS API](https://docs.aws.amazon.com/STS/latest/APIReference/API_DecodeAuthorizationMessage.html) - this may or may not clarify a permission issue).

When applied to a trust policy, the requirement is only enforced when _assuming the role_. Since a trust policy usually only regulates the acceptable principals (e.g. an AWS account) and second factor conditions the user should be able to deduce why the permission to assume the role has been denied (e.g. no access to the admin role in this account, not in the office / not dialed into the VPN, no MFA device active). Both approaches can also be mixed either through policies or through permission boundaries (under the usability constraints described above).

Since an assumed role also has an explicit time-to-live that is clearly visible to users only using trust policies for MFA conditions makes MFA based second factors easier to handle in practice. This approach can be summarized as following:

1. A role requires the proof of posession once at the beginning of it's lifetime through an `aws:MultiFactorAuthAge` condition that should set to the duration attribute (the maximum lifetime) of a role (as discussed above: since an older session token could have been used to assume the role, the "boolean" form is not sufficient)
2. The longer a role can last, the greater the window of opportunity gets for an attacker to gain control of the principal (e.g. by gaining access to the computer where the short-term credentials reside)
3. Therefore long lasting session should have less permissions than shorter lasting roles

This summary leads to the following basic setup recommendations:

1. Plan a permission schema that differentiates between the levels of access different functional roles should have - [typical examples](https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies_job-functions.html) include the AWS Managed Policies for Job Functions as well as e.g. roles for removing MFA devices from IAM users in self-service
2. Map those functional roles to IAM roles
   1. Role durations and matching MFA conditions should reflect the level of access the associated policies will have
   2. Consider that API operations carried with roles that expire mid-deployment will get aborted mid-deployment and that 1 hour might not be enough for certain operations to complete (such as e.g. RDS deployments) - `awsu` handles this problem by re-assuming roles that are due to expire in a configurable amount of time
   3. Consider adopting an [AWS Organization setup with multiple member accounts](https://aws.amazon.com/answers/account-management/aws-multi-account-security-strategy/) and locate these roles in different accounts
3. Map the right to assume these IAM roles to IAM groups
   1. Also give users in these groups (by using [policy variables](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_variables.html)) the right to [create](https://docs.aws.amazon.com/IAM/latest/APIReference/API_CreateVirtualMFADevice.html) and [enable](https://docs.aws.amazon.com/IAM/latest/APIReference/API_EnableMFADevice.html) a virtual MFA device for themselves
   2. Do not directly allow them to [deactive](https://docs.aws.amazon.com/IAM/latest/APIReference/API_DeactivateMFADevice.html) and [delete](https://docs.aws.amazon.com/IAM/latest/APIReference/API_DeleteVirtualMFADevice.html) their own MFA device (this would allow an attacker to simply replace an MFA device with one they control) - monitor deletions via CloudTrail and consider outsourcing this to dedicated privileged users or - at minimum - require second factors for the MFA deletion alone (directly or via a role dedicated to this purpose)
4. Do not give IAM users any other direct permissions
5. Add IAM users to IAM groups appropriate for their function roles

## Protecting second factor requirements

With second factor conditions in various policies (likely in multiple accounts) in place one must consider protecting them from modification by regular and privileged users.

Privileged users in this context are users that _can_ create and update IAM resources beyond the self-service permissions describe in the section above. Such privileged users are an unfortunate reality in AWS where service-linked permissions are not always fine-granular enough or are not supported at all. The unfortunate part is that it's currently [impossible to restrict](https://docs.aws.amazon.com/IAM/latest/UserGuide/list_identityandaccessmanagement.html) role creation to service principals which implies that privileged users can create cross-account roles to arbitrary accounts.

The following approaches are recommend in order to protect your second factor requirements:

1. Implement all IAM  that regulate access under a dedicated [path](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_identifiers.html), e.g. `/master/`
2. Create an IAM policy ("master") that restricts create, update and delete operations on IAM resources under this path
   1. attach this IAM policy as a permission boundary to your IAM roles
   2. enforce the attachment of this permission boundary by expressing it as a `iam:PermissionsBoundary` condition to all [applicate IAM actions](https://docs.aws.amazon.com/IAM/latest/UserGuide/list_identityandaccessmanagement.html)
3. Carefully monitor the creation of roles closely and consider implementing indirections (e.g. using the [AWS Service Catalog](https://aws.amazon.com/servicecatalog/)) as an alternative to directly creating IAM roles (e.g. for usage in [instance profiles](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_use_switch-role-ec2_instance-profiles.html)) and IAM users (e.g. for creating [SMTP credentials](https://docs.aws.amazon.com/ses/latest/DeveloperGuide/smtp-credentials.html) for SES)

Please don't hesitate to [contact us](https://kreuzwerker.de/) when you need consulting around the design of your organizational setup.
