package command

import (
	"os"

	"github.com/kreuzwerker/awsu/config"
	"github.com/kreuzwerker/awsu/session"
)

func newSession(workspace string) (*session.Session, error) {

	cwd, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	cfg, err := config.Load(cwd)

	if err != nil {
		return nil, err
	}

	if workspace == "" {
		workspace = cfg.DetectWorkspace()
	}

	profile, err := cfg.Get(workspace)

	if err != nil {
		return nil, err
	}

	sess, err := session.Load(cwd, profile)

	if err != nil || !sess.IsValid() {

		sess, err = session.New(profile)

		if err != nil {
			return nil, err
		}

		if err := sess.Save(cwd); err != nil {
			return nil, err
		}

	}

	return sess, nil

}
