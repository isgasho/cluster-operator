package kubeconfig

import (
	"context"

	"github.com/giantswarm/certs"
	"github.com/giantswarm/kubeconfig"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/client/k8srestconfig"
	"github.com/giantswarm/operatorkit/controller/context/reconciliationcanceledcontext"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/cluster-operator/pkg/label"
	"github.com/giantswarm/cluster-operator/pkg/v14/key"
	awskey "github.com/giantswarm/cluster-operator/service/controller/aws/v14/key"
)

func (r *StateGetter) GetDesiredState(ctx context.Context, obj interface{}) ([]*corev1.Secret, error) {
	cr, err := awskey.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	if awskey.IsDeleted(cr) {
		return []*corev1.Secret{}, nil
	}

	clusterGuestConfig := awskey.ClusterGuestConfig(cr)
	apiDomain, err := key.APIDomain(clusterGuestConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	appOperator, err := r.certsSearcher.SearchAppOperator(clusterGuestConfig.ID)
	if certs.IsTimeout(err) {
		r.logger.LogCtx(ctx, "level", "debug", "message", "did not get an app-operator-api cert for the tenant cluster")

		// We can't continue without a app-operator-api cert. We will retry during the
		// next execution.
		reconciliationcanceledcontext.SetCanceled(ctx)
		r.logger.LogCtx(ctx, "level", "debug", "message", "canceling reconciliation")

		return []*corev1.Secret{}, nil
	} else if err != nil {
		return nil, microerror.Mask(err)
	}

	c := k8srestconfig.Config{
		Logger: r.logger,

		Address:   apiDomain,
		InCluster: false,
		TLS: k8srestconfig.ConfigTLS{
			CAData:  appOperator.APIServer.CA,
			CrtData: appOperator.APIServer.Crt,
			KeyData: appOperator.APIServer.Key,
		},
	}
	restConfig, err := k8srestconfig.New(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	yamlBytes, err := kubeconfig.NewKubeConfigForRESTConfig(ctx, restConfig, key.KubeConfigClusterName(clusterGuestConfig), "")
	if err != nil {
		return nil, microerror.Mask(err)
	}

	secretName := key.KubeConfigSecretName(clusterGuestConfig)

	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: r.resourceNamespace,
			Labels: map[string]string{
				label.Cluster:      clusterGuestConfig.ID,
				label.ManagedBy:    r.projectName,
				label.Organization: clusterGuestConfig.Owner,
				label.ServiceType:  label.ServiceTypeManaged,
			},
		},
		Data: map[string][]byte{
			"kubeConfig": yamlBytes,
		},
	}

	return []*corev1.Secret{&secret}, nil
}