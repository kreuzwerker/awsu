package session

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/kreuzwerker/awsu/config"
	"github.com/kreuzwerker/awsu/log"
	"github.com/kreuzwerker/awsu/yubikey"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/yawn/envmap"
)

var (
	errFailedToAssumeRole        = "failed to assume role %q: %s"
	errFailedToGetCallerIdentity = "failed to get caller identity during MFA autodetection: %s"
	errFailedToGetSessionToken   = "failed to get session token: %s"
	errNoMFAFound                = "no MFA found in profile or source profile definition"
	errNoSuchProfile             = "no such profile %q configured"
	errNoSuchSourceProfile       = "no such source profile %q configured"
	expires                      = 1 * time.Hour
)

type Session struct {
	credentials.Value
	Expires time.Time
	Profile string
}

// New creates a new session with the given profile, using the passed in profiles
func New(profile string, profiles config.Profiles) (*Session, error) {

	var (
		lt, st *config.Profile
		sess   = session.Must(session.NewSession())
	)

	if lt = profiles[profile]; lt == nil {
		return nil, fmt.Errorf(errNoSuchProfile, profile)
	}

	// split into long and short term
	if !lt.IsLongTerm() {
		st = lt
		lt = profiles[st.SourceProfile]
	}

	if lt == nil {
		return nil, fmt.Errorf(errNoSuchSourceProfile, profile)
	}

	session := &Session{
		Expires: time.Now().Add(expires),
		Profile: profile,
		Value:   lt.Value(),
	}

	{

		client := sts.New(sess, aws.NewConfig().WithCredentials(session.Credentials()))

		var mfa string

		if mfa = lt.MFASerial; mfa != "" {
			log.Log("get session token on profile %q, using MFA from long-term credentials", lt.Name)
		} else if mfa = st.MFASerial; mfa != "" {
			log.Log("get session token on profile %q, using MFA from short-term credentials", lt.Name)
		} else {

			// TODO: autodetection

			// arn:aws:iam::113030722353:user/joern.barthel@kreuzwerker.de vs MFA of
			// arn:aws:iam::113030722353:mfa/joern.barthel@kreuzwerker.de

			return nil, fmt.Errorf(errNoMFAFound)

		}

		token, err := yubikey.Generate(mfa)

		if err != nil {
			return nil, err
		}

		res, err := client.GetSessionToken(&sts.GetSessionTokenInput{
			DurationSeconds: aws.Int64(int64(expires.Seconds())),
			SerialNumber:    &mfa,
			TokenCode:       &token,
		})

		if err != nil {
			return nil, fmt.Errorf(errFailedToGetSessionToken, err)
		}

		session.setCredentials(res.Credentials)

	}

	if st != nil {

		client := sts.New(sess, aws.NewConfig().WithCredentials(session.Credentials()))

		log.Log("assuming role %q using profile %s", st.RoleARN, st.Name)

		req := &sts.AssumeRoleInput{
			DurationSeconds: aws.Int64(int64(expires.Seconds())),
			RoleArn:         &st.RoleARN,
			RoleSessionName: aws.String("test"), // TODO: session names based on user id
		}

		if st.ExternalID != "" {
			req.ExternalId = &st.ExternalID
		}

		res, err := client.AssumeRole(req)

		if err != nil {
			return nil, fmt.Errorf(errFailedToAssumeRole, st.RoleARN, err)
		}

		session.setCredentials(res.Credentials)

	}

	return session, nil

}

// Credentials returns a new set of static credentials
func (s *Session) Credentials() *credentials.Credentials {
	return credentials.NewStaticCredentialsFromCreds(s.Value)
}

// Exec sets an appropriate runtime environment and execs the passed in app
func (s *Session) Exec(app string, args []string) error {

	env := envmap.Import()

	env["AWSU_EXPIRES"] = s.Expires.Format(time.RFC3339)
	env["AWSU_PROFILE"] = s.Profile
	env["AWS_ACCESS_KEY_ID"] = s.Value.AccessKeyID
	env["AWS_SECRET_ACCESS_KEY"] = s.Value.SecretAccessKey
	env["AWS_SESSION_TOKEN"] = s.Value.SessionToken

	cmd, err := exec.LookPath(app)

	if err != nil {
		return err
	}

	log.Log("running %q with args %q", cmd, args)

	return syscall.Exec(cmd, args, env.ToEnv())

}

// IsValid determines if the session is still not expired
func (s *Session) IsValid() bool {

	grace := 15 * time.Minute

	return time.Now().Add(grace).Before(s.Expires)
}

// Load loads a cached session
func Load(dir, profile string) (*Session, error) {

	path, err := sessionPath(dir, profile)

	if err != nil {
		return nil, err
	}

	log.Log("loading session from %q", path)

	raw, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var session *Session

	if err := json.Unmarshal(raw, &session); err != nil {
		return nil, err
	}

	if !session.IsValid() {
		return nil, fmt.Errorf("existing session is invalid")
	}

	return session, nil

}

// Save caches a session
func (s *Session) Save(dir string) error {

	path, err := sessionPath(dir, s.Profile)

	if err != nil {
		return err
	}

	log.Log("saving session to %q", path)

	out, err := json.Marshal(s)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, out, 0600)

}

// sessionPath determines a caching for the given session
func sessionPath(dir, profile string) (string, error) {

	home, err := homedir.Dir()

	if err != nil {
		return "", err
	}

	dir = filepath.Join(home, ".awsu", "sessions")

	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}

	return filepath.Join(dir, fmt.Sprintf("%s.json", profile)), nil

}

// setCredentials imports a new set of STS credentials
func (s *Session) setCredentials(c *sts.Credentials) {
	s.AccessKeyID = *c.AccessKeyId
	s.SecretAccessKey = *c.SecretAccessKey
	s.SessionToken = *c.SessionToken
}

// String returns a string representation of this session, suitable for eval()
func (s *Session) String() string {

	parts := []string{
		fmt.Sprintf("export AWS_ACCESS_KEY_ID=%s", s.AccessKeyID),
		fmt.Sprintf("export AWS_SECRET_ACCESS_KEY=%s", s.SecretAccessKey),
	}

	if s.SessionToken != "" {
		parts = append(parts, fmt.Sprintf("export AWS_SESSION_TOKEN=%s", s.SessionToken))
	}

	return strings.Join(parts, "\n")

}
