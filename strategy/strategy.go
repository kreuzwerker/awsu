package strategy

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	human "github.com/dustin/go-humanize"
	"github.com/kreuzwerker/awsu/config"
	"github.com/kreuzwerker/awsu/log"
	"github.com/kreuzwerker/awsu/strategy/credentials"
)

const (
	errProfileAquisitionFailed = "failed to aquire credentials for profile %q: %s"
	errProfileCacheLoadExpired = "ignoring expired cached profile %q"
	errProfileCacheLoadFailed  = "failed to load cached profile %q: %s"
	errProfileCacheSaveFailed  = "failed to save cached profile %q: %s"
	errProfileNotFound         = "no such profile %q configured"
	logSessionExpires          = "session will expire (after applying grace) %s"
	logStrategyWithProfile     = "using strategy %q (cache: %t) for profile %q"
)

// Strategy identifies a way of aquiring short-term, cacheable credentials
type Strategy interface {
	// Credentials aquires actual credentials
	Credentials(*session.Session) (*credentials.Credentials, error)

	// Name returns the name of this strategy
	Name() string

	// IsCacheable indicates the output of this strategy can be cached
	IsCacheable() bool

	// Profiles returns the name of the profile used (if applicable)
	Profile() *config.Profile
}

// Apply applies all configured strategy, depending on the given Config
func Apply(cfg *config.Config) (*credentials.Credentials, error) {

	var (
		sess   *session.Session
		source *config.Profile
		target = cfg.Profiles[cfg.Profile]
	)

	if target == nil {
		return nil, fmt.Errorf(errProfileNotFound, cfg.Profile)
	}

	source = cfg.Profiles[target.SourceProfile]

	strategies := []Strategy{
		&LongTerm{
			Profiles: []*config.Profile{source, target},
		},
		&SessionToken{
			Duration:  cfg.Duration,
			Generator: cfg.Generator,
			Grace:     cfg.Grace,
			MFASerial: cfg.MFASerial,
			Profiles:  []*config.Profile{source, target},
		},
		&AssumeRole{
			Duration: cfg.Duration,
			Grace:    cfg.Grace,
			Profiles: []*config.Profile{source, target},
		},
	}

	var last *credentials.Credentials

	for _, a := range strategies {

		var current *credentials.Credentials

		profile := a.Profile()

		if profile == nil {
			continue
		}

		cache := !cfg.NoCache && a.IsCacheable()

		log.Debug(logStrategyWithProfile, a.Name(), cache, profile.Name)

		// try to load
		if cache {

			creds, err := credentials.Load(profile.Name)

			if err != nil {
				log.Debug(errProfileCacheLoadFailed, profile.Name, err)
			} else if !creds.IsValid() {
				log.Debug(errProfileCacheLoadExpired, profile.Name)
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

		if sess == nil {
			sess = current.NewSession()
		} else {
			sess = current.UpdateSession(sess)
		}

	}

	if last.Expires.Second() > 0 {
		log.Debug(logSessionExpires, human.Time(last.Expires))
	}

	return last, nil

}
