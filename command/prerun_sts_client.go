package command

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/cobra"
)

var stsClient = new(prerunSTSClient)

type prerunSTSClient struct {
	*sts.STS
}

func (p *prerunSTSClient) PreRun(cmd *cobra.Command, args []string) error {

	sess, err := session.NewSession()

	if err != nil {
		return fmt.Errorf("error while opening new session: %s", err)
	}

	p.STS = sts.New(sess)

	return nil

}
