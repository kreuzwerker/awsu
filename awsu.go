package awsu

const (
	AccessKeyID     = "AWS_ACCESS_KEY_ID"
	SecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	SessionActive   = "AWSU_SESSION_ACTIVE"
	SessionToken    = "AWS_SESSION_TOKEN"
)

var AllKeys = []string{
	AccessKeyID,
	SecretAccessKey,
	SessionActive,
	SessionToken,
}
