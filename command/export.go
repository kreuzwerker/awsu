package command

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kreuzwerker/awsu"
	"github.com/spf13/cobra"
	"github.com/yawn/envmap"
)

var exportCmd = &cobra.Command{

	Use:   "export",
	Short: "Exports shell variables",
	RunE: func(cmd *cobra.Command, args []string) error {

		f := envmap.PrefixedKeys("_", strings.Join(awsu.AllKeys, "|")).MatchString
		e := envmap.Import()

		if e.IsSet(awsu.SessionActive) {

			// temporarily remove keys
			for k := range e.Subset(f) {
				os.Unsetenv(k)
			}

		}

		var (
			hash = section.KeysHash()
			env  envmap.Envmap
		)

		client, err := awsu.NewClient()

		if err != nil {
			return err
		}

		var (
			aki = hash["aws_access_key_id"]
			arn = hash["role_arn"]
			ask = hash["aws_secret_access_key"]
			eid = hash["external_id"]
		)

		if aki != "" && ask != "" {

			log.Println("found keypair, re-exporting from credentials")

			res := map[string]string{
				awsu.AccessKeyID:     hash["aws_access_key_id"],
				awsu.SecretAccessKey: hash["aws_secret_access_key"],
			}

			env = res

		} else if arn != "" {

			log.Println("found role, assuming it")

			res, err := client.AssumeRole(&awsu.Options{
				ARN:        arn,
				ExternalID: eid,
			})

			if err != nil {
				return err
			}

			env = res

		} else {
			return fmt.Errorf("no strategy for handling profile %q", rootFlags.profile)
		}

		env[awsu.SessionActive] = rootFlags.profile

		for _, e := range env.Subset(f).ToEnv() {
			fmt.Printf("export %s\n", e)
		}

		return nil

	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
