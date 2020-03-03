package config

import (
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
	ini "gopkg.in/ini.v1"
)

const (
	errFailedToOpen = "failed to open %q"
)

// Profiles is map of profile names to profiles
type Profiles map[string]*Profile

// Source resolves a short-term credential source profile
func (p Profiles) Source(profile string) *Profile {
	return p[profile]
}

// Load will load profiles from multiple files, merging definitions on the way
func Load(files ...string) (Profiles, error) {

	var (
		loaded   []string
		profiles = make(map[string]*Profile)
	)

	for _, file := range files {

		if file == "" {
			continue
		}

		buf, err := ioutil.ReadFile(file)

		if err != nil {
			return nil, errors.Wrapf(err, errFailedToOpen, file)
		}

		loaded = append(loaded, file)

		f, err := ini.Load(buf)

		if err != nil {
			return nil, err
		}

		for _, section := range f.Sections() {

			name := section.Name()
			name = strings.TrimPrefix(name, "profile ")

			if name == "preview" {
				continue
			}

		init:

			profile, ok := profiles[name]

			if !ok {
				profiles[name] = new(Profile)
				goto init
			}

			sectionProfile := new(Profile)
			if err := section.MapTo(sectionProfile); err != nil {
				return nil, err
			}
			sectionProfile.Name = name
			profile.Merge(sectionProfile)
		}

	}

	return profiles, nil

}
