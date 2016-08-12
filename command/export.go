package command

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/cobra"
	"github.com/yawn/envmap"
)

var exportCmd = &cobra.Command{

	Use:     "export",
	Short:   "Exports shell variables",
	PreRunE: stsClient.PreRun,
	RunE: func(cmd *cobra.Command, args []string) error {		// var (
			// 	exp = fmt.Sprintf(`^(?:_)*(?:%s)$`, strings.Join([]string{aki, ask, ast}, "|"))
			// 	f   = regexp.MustCompile(exp).MatchString
			// )
			//
			// e1 := envmap.Import()
			// e2 := e1.Push("_", f).Push("_", f).Push("_", f)
			//
			// fmt.Println(e2.Subset(f))
			//
			// if true {
			// 	return nil
			// }

		hash := config.section.KeysHash()

		params := &sts.AssumeRoleInput{
			DurationSeconds: aws.Int64(3600),
			RoleSessionName: this.Session(),
		}

		var (
			externalID = config.section.HasKey("external_id")
			role       = config.section.HasKey("role_arn")
		)

		if role {

			params.RoleArn = aws.String(hash["role_arn"])

			if externalID {
				params.ExternalId = aws.String(hash["external_id"])
			}

			res, err := stsClient.AssumeRole(params)

			if err != nil {
				return fmt.Errorf("error during assume role call: %s", err)
			}

			var e2  envmap.Envmap = make(map[string]string, 3)

			e2[aki] = *res.Credentials.AccessKeyId
			e2[ask] = *res.Credentials.SecretAccessKey
			e2[ast] = *res.Credentials.SessionToken

			for k, v := range e2 {
				fmt.Printf("unset %s\n", k)
				fmt.Printf("export %s=%s\n", k, v)
			}

		} else {
			return fmt.Errorf("no way to assume role for profile %q with keys %s",
				rootFlags.profile,
				config.section.KeyStrings())
		}

		return nil

	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
