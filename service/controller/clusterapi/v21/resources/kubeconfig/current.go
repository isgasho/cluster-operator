package kubeconfig

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/controller/context/resourcecanceledcontext"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/cluster-operator/service/controller/clusterapi/v21/key"
)

func (r *Resource) GetCurrentState(ctx context.Context, obj interface{}) ([]*corev1.Secret, error) {
	cr, err := key.ToCluster(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// The kube config secret is deleted implicitely by the provider operator when
	// it deletes the tenant cluster namespace in the control plane.
	if key.IsDeleted(&cr) {
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("not deleting secret %#q for tenant cluster %#q", key.KubeConfigSecretName(&cr), key.ClusterID(&cr)))
		r.logger.LogCtx(ctx, "level", "debug", "message", "canceling resource")
		resourcecanceledcontext.SetCanceled(ctx)
		return nil, nil
	}

	var secret *corev1.Secret
	{
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("finding secret %#q for tenant cluster %#q", key.KubeConfigSecretName(&cr), key.ClusterID(&cr)))

		secret, err = r.k8sClient.CoreV1().Secrets(key.ClusterID(&cr)).Get(key.KubeConfigSecretName(&cr), metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("did not find secret %#q for tenant cluster %#q", key.KubeConfigSecretName(&cr), key.ClusterID(&cr)))
			return nil, nil

		} else if err != nil {
			return nil, microerror.Mask(err)
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("found secret %#q for tenant cluster %#q", key.KubeConfigSecretName(&cr), key.ClusterID(&cr)))
	}

	return []*corev1.Secret{secret}, nil
}