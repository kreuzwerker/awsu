package command

import (
	"os"

	"github.com/kreuzwerker/awsu/config"
	"github.com/kreuzwerker/awsu/log"
	"github.com/kreuzwerker/awsu/session"
)

func newSession(nocache bool, profile string, profiles config.Profiles) (*session.Session, error) {

	if nocache {
		return session.New(profile, profiles)
	}

	cwd, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	sess, err := session.Load(cwd, profile)

	if err != nil {

		log.Log("no previous session, creating")
		return restoreSession(cwd, profile, profiles)

	} else if !sess.IsValid() {

		log.Log("invalid previous session, recreating")
		return restoreSession(cwd, profile, profiles)

	}

	return sess, nil

}

func restoreSession(cwd string, profile string, profiles config.Profiles) (*session.Session, error) {

	sess, err := session.New(profile, profiles)

	if err != nil {
		return nil, err
	}

	if err = sess.Save(cwd); err != nil {
		return nil, err
	}

	return sess, err

}
