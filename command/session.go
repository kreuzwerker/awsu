package command

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kreuzwerker/awsu/aquirer"
	"github.com/kreuzwerker/awsu/config"
	"github.com/kreuzwerker/awsu/log"
)

const (
	errProfileAquisitionFailed = "failed to aquire credentials for profile %q: %s"
	errProfileCacheLoadExpired = "ignoring expired cached profile %q"
	errProfileCacheLoadFailed  = "failed to load cached profile %q: %s"
	errProfileCacheSaveFailed  = "failed to save cached profile %q: %s"
	errProfileNotFound         = "no such profile %q configured"
	msgAquirerWithProfile      = "using aquirer %q (cache: %t) for profile %q"
)

func newSession(cfg *config.Config) (*aquirer.Credentials, error) {

	var (
		sess   = session.Must(session.NewSession())
		source *config.Profile
		target = cfg.Profiles[cfg.Profile]
	)

	if target == nil {
		return nil, fmt.Errorf(errProfileNotFound, cfg.Profile)
	}

	source = cfg.Profiles[target.SourceProfile]

	aquirers := []aquirer.Aquirer{
		&aquirer.LongTerm{
			Profiles: []*config.Profile{source, target},
		},
		&aquirer.SessionToken{
			Duration:  cfg.SessionTTL,
			Grace:     cfg.SessionTTL / 2,
			MFASerial: cfg.MFASerial,
			Profiles:  []*config.Profile{source, target},
		},
		&aquirer.AssumeRole{
			Duration: cfg.CacheTTL,
			Grace:    cfg.CacheTTL / 2,
			Profiles: []*config.Profile{source, target},
		},
	}

	var last *aquirer.Credentials

	for _, a := range aquirers {

		var current *aquirer.Credentials

		profile := a.Profile()

		if profile == nil {
			continue
		}

		cache := !cfg.NoCache && a.IsCacheable()

		log.Log(msgAquirerWithProfile, a.Name(), cache, profile.Name)

		// try to load
		if cache {

			creds, err := aquirer.LoadCredentials(profile.Name)

			if err != nil {
				log.Log(errProfileCacheLoadFailed, profile.Name, err)
			} else if !creds.IsValid() {
				log.Log(errProfileCacheLoadExpired, profile.Name)
			} else {
				current = creds
			}

		}

		// try to aquire
		if current == nil {

			creds, err := a.Credentials(sess)

			if err != nil {
				return nil, fmt.Errorf(errProfileAquisitionFailed, profile.Name, err)
			}

			current = creds

			// try to save
			if cache {

				if err := current.Save(); err != nil {
					return nil, fmt.Errorf(errProfileCacheSaveFailed, profile.Name, err)
				}

			}

		}

		last = current
		sess = current.UpdateSession(sess)

	}

	return last, nil

}
