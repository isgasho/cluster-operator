package configmap

import (
	"context"
	"fmt"

	"github.com/giantswarm/errors/tenant"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/controller/context/resourcecanceledcontext"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/cluster-operator/pkg/label"
	"github.com/giantswarm/cluster-operator/pkg/project"
)

func (s *Service) GetCurrentState(ctx context.Context, clusterConfig ClusterConfig) ([]*corev1.ConfigMap, error) {
	var currentConfigMaps []*corev1.ConfigMap

	tenantK8sClient, err := s.newTenantK8sClient(ctx, clusterConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// Namespaces used by all providers. Uses a map for deduping.
	namespaces := map[string]bool{
		metav1.NamespaceSystem: true,
	}

	// Add any provider specific namespaces.
	for _, namespace := range clusterConfig.Namespaces {
		namespaces[namespace] = true
	}

	listOptions := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s, %s=%s", label.ServiceType, label.ServiceTypeManaged, label.ManagedBy, project.Name()),
	}

	for namespace := range namespaces {
		configMapList, err := tenantK8sClient.CoreV1().ConfigMaps(namespace).List(listOptions)
		if tenant.IsAPINotAvailable(err) {
			s.logger.LogCtx(ctx, "level", "debug", "message", "tenant cluster is not available yet")
			s.logger.LogCtx(ctx, "level", "debug", "message", "canceling resource")
			resourcecanceledcontext.SetCanceled(ctx)
			return nil, nil

		} else if err != nil {
			return nil, microerror.Mask(err)
		}

		for _, item := range configMapList.Items {
			c := item.DeepCopy()
			currentConfigMaps = append(currentConfigMaps, c)
		}
	}

	return currentConfigMaps, nil
}