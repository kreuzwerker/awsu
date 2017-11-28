package session

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"encoding/json"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kreuzwerker/awsu/log"
	"github.com/kreuzwerker/awsu/metadata"
	"github.com/kreuzwerker/awsu/yubikey"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/yawn/envmap"
)

type Session struct {
	credentials.Value
	Expires time.Time
	metadata.Metadata
	Profile string
}

func (s *Session) Exec(app string, args []string) error {

	env := envmap.Import()

	env["AWSU_EXPIRES"] = s.Expires.Format(time.RFC3339)
	env["AWSU_PROFILE"] = s.Profile

	env["AWS_ACCESS_KEY_ID"] = s.Value.AccessKeyID
	env["AWS_SECRET_ACCESS_KEY"] = s.Value.SecretAccessKey

	if s.Value.SessionToken != "" {
		env["AWS_SESSION_TOKEN"] = s.Value.SessionToken
	}

	cmd, err := exec.LookPath(app)

	if err != nil {
		return err
	}

	log.Log("running %q with args %q", cmd, args)

	return syscall.Exec(cmd, args, env.ToEnv())

}

func (s *Session) IsValid() bool {

	grace := 1 * time.Minute

	return time.Now().Add(grace).Before(s.Expires)
}

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

func New(profile string) (*Session, error) {

	log.Log("creating new session with profile %q", profile)

	var (
		sess *session.Session
		tp   = new(yubikey.TokenProvider)
	)

	// this is currently the only way to ensure this ttl
	stscreds.DefaultDuration = 1 * time.Hour

	sess = session.Must(session.NewSessionWithOptions(session.Options{
		AssumeRoleTokenProvider: tp.Provide,
		Profile:                 profile,
		SharedConfigState:       session.SharedConfigEnable,
	}))

	tp.Session = sess

	value, err := sess.Config.Credentials.Get()

	if err != nil {
		return nil, err
	}

	session := &Session{
		Expires:  time.Now().Add(stscreds.DefaultDuration),
		Metadata: metadata.New(sess.Config.Credentials),
		Profile:  profile,
		Value:    value,
	}

	return session, nil

}

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

func sessionPath(dir, profile string) (string, error) {

	hash := sha1.Sum([]byte(dir))

	home, err := homedir.Dir()

	if err != nil {
		return "", err
	}

	dir = filepath.Join(home, ".awsu", "sessions", fmt.Sprintf("%x", hash[:]))

	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}

	return filepath.Join(dir, fmt.Sprintf("%s.json", profile)), nil

}
