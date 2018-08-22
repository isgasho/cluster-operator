package configmap

import (
	"context"
	"encoding/json"

	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/cluster-operator/pkg/label"
)

type configMapGenerator func(ctx context.Context, configMapValues ConfigMapValues, projectName string) (*corev1.ConfigMap, error)

type BasicConfigMap struct {
	Image Image `json:"image"`
}

type IngressController struct {
	Controller IngressControllerController `json:"controller"`
	Global     IngressControllerGlobal     `json:"global"`
	Image      Image                       `json:"image"`
}

type IngressControllerController struct {
	Replicas int                                `json:"replicas"`
	Service  IngressControllerControllerService `json:"service"`
}

type IngressControllerControllerService struct {
	Enabled bool `json:"enabled"`
}

type IngressControllerGlobal struct {
	Controller IngressControllerGlobalController `json:"controller"`
	Migration  IngressControllerGlobalMigration  `json:"migration"`
}

type IngressControllerGlobalController struct {
	Replicas         int  `json:"replicas"`
	UseProxyProtocol bool `json:"useProxyProtocol"`
}

type IngressControllerGlobalMigration struct {
	Enabled bool `json:"enabled"`
}

type Image struct {
	Registry string `json:"registry"`
}

type CertExporter struct {
	Namespace string `json:"namespace"`
}

type NetExporter struct {
	Namespace string `json:"namespace"`
}

func (s *Service) GetDesiredState(ctx context.Context, configMapValues ConfigMapValues) ([]*corev1.ConfigMap, error) {
	desiredConfigMaps := make([]*corev1.ConfigMap, 0)

	generators := []configMapGenerator{
		s.newCertExporterConfigMap,
		s.newIngressControllerConfigMap,
		s.newKubeStateMetricsConfigMap,
		s.newNetExporterConfigMap,
		s.newNodeExporterConfigMap,
	}

	for _, g := range generators {
		configMap, err := g(ctx, configMapValues, s.projectName)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		desiredConfigMaps = append(desiredConfigMaps, configMap)
	}

	return desiredConfigMaps, nil
}

func (s *Service) newCertExporterConfigMap(ctx context.Context, configMapValues ConfigMapValues, projectName string) (*corev1.ConfigMap, error) {
	configMapName := "cert-exporter-values"
	appName := "cert-exporter"
	labels := newConfigMapLabels(configMapValues, appName, projectName)

	values := CertExporter{
		Namespace: metav1.NamespaceSystem,
	}
	json, err := json.Marshal(values)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	data := map[string]string{
		"values.json": string(json),
	}

	newConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: metav1.NamespaceSystem,
			Labels:    labels,
		},
		Data: data,
	}

	return newConfigMap, nil
}

func (s *Service) newIngressControllerConfigMap(ctx context.Context, configMapValues ConfigMapValues, projectName string) (*corev1.ConfigMap, error) {
	configMapName := "nginx-ingress-controller-values"
	appName := "nginx-ingress-controller"
	labels := newConfigMapLabels(configMapValues, appName, projectName)

	// controllerServiceEnabled needs to be set separately for the chart
	// migration logic but is the reverse of migration enabled.
	controllerServiceEnabled := !configMapValues.IngressControllerMigrationEnabled

	values := IngressController{
		Controller: IngressControllerController{
			Replicas: configMapValues.WorkerCount,
			Service: IngressControllerControllerService{
				Enabled: controllerServiceEnabled,
			},
		},
		Global: IngressControllerGlobal{
			Controller: IngressControllerGlobalController{
				Replicas:         configMapValues.WorkerCount,
				UseProxyProtocol: configMapValues.IngressControllerUseProxyProtocol,
			},
			Migration: IngressControllerGlobalMigration{
				Enabled: configMapValues.IngressControllerMigrationEnabled,
			},
		},
		Image: Image{
			Registry: s.registryDomain,
		},
	}
	json, err := json.Marshal(values)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	data := map[string]string{
		"values.json": string(json),
	}

	newConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: metav1.NamespaceSystem,
			Labels:    labels,
		},
		Data: data,
	}

	return newConfigMap, nil
}

func (s *Service) newKubeStateMetricsConfigMap(ctx context.Context, configMapValues ConfigMapValues, projectName string) (*corev1.ConfigMap, error) {
	return s.newBasicConfigMap(ctx, configMapValues, projectName, "kube-state-metrics")
}

func (s *Service) newNetExporterConfigMap(ctx context.Context, configMapValues ConfigMapValues, projectName string) (*corev1.ConfigMap, error) {
	configMapName := "net-exporter-values"
	appName := "net-exporter"
	labels := newConfigMapLabels(configMapValues, appName, projectName)

	values := NetExporter{
		Namespace: metav1.NamespaceSystem,
	}
	json, err := json.Marshal(values)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	data := map[string]string{
		"values.json": string(json),
	}

	newConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: metav1.NamespaceSystem,
			Labels:    labels,
		},
		Data: data,
	}

	return newConfigMap, nil
}

func (s *Service) newNodeExporterConfigMap(ctx context.Context, configMapValues ConfigMapValues, projectName string) (*corev1.ConfigMap, error) {
	return s.newBasicConfigMap(ctx, configMapValues, projectName, "node-exporter")
}

func (s *Service) newBasicConfigMap(ctx context.Context, configMapValues ConfigMapValues, projectName string, appName string) (*corev1.ConfigMap, error) {
	configMapName := appName + "-values"
	labels := newConfigMapLabels(configMapValues, appName, projectName)

	values := BasicConfigMap{
		Image: Image{
			Registry: s.registryDomain,
		},
	}
	json, err := json.Marshal(values)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	data := map[string]string{
		"values.json": string(json),
	}

	newConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: metav1.NamespaceSystem,
			Labels:    labels,
		},
		Data: data,
	}

	return newConfigMap, nil
}

func newConfigMapLabels(configMapValues ConfigMapValues, appName, projectName string) map[string]string {
	return map[string]string{
		label.App:          appName,
		label.Cluster:      configMapValues.ClusterID,
		label.ManagedBy:    projectName,
		label.Organization: configMapValues.Organization,
		label.ServiceType:  label.ServiceTypeManaged,
	}
}