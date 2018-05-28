// +build k8srequired

package teardown

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/giantswarm/e2e-harness/pkg/framework"
	awsclient "github.com/giantswarm/e2eclients/aws"
	"github.com/giantswarm/helmclient"
	"github.com/giantswarm/microerror"
	"k8s.io/helm/pkg/helm"
)

func Resources(c *awsclient.Client, f *framework.Host, helmClient *helmclient.Client) error {
	items := []string{
		"cluster-operator",
		"cluster-operator-resource",
		"cert-operator",
		"cert-resource-lab",
		"aws-operator",
		"aws-resource-lab",
		"node-operator",
	}

	for _, item := range items {
		err := helmClient.DeleteRelease(item, helm.DeletePurge(true))
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func HostPeerVPC(c *awsclient.Client) error {
	log.Printf("Deleting Host Peer VPC stack")

	_, err := c.CloudFormation.DeleteStack(&cloudformation.DeleteStackInput{
		StackName: aws.String("host-peer-" + os.Getenv("CLUSTER_NAME")),
	})
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}