package configmap

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger/microloggertest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/cluster-operator/pkg/label"
)

const (
	basicMatchJSON = `
	{
		"controller": {
			"replicas": 3,
			"service": {
				"enabled": false
			}
		},
		"global": {
			"controller": {
				"tempReplicas": 2,
				"useProxyProtocol": true
			},
			"migration": {
				"enabled": true
			}
		},
		"image": {
			"registry": "quay.io"
		}
	}
	`
	differentWorkerCountJSON = `
	{
		"controller": {
			"replicas": 7,
			"service": {
				"enabled": false
			}
		},
		"global": {
			"controller": {
				"tempReplicas": 4,
				"useProxyProtocol": true
			},
			"migration": {
				"enabled": true
			}
		},
		"image": {
			"registry": "quay.io"
		}
	}
	`
	differentSettingsJSON = `
	{
		"controller": {
			"replicas": 1,
			"service": {
				"enabled": true
			}
		},
		"global": {
			"controller": {
				"tempReplicas": 1,
				"useProxyProtocol": false
			},
			"migration": {
				"enabled": false
			}
		},
		"image": {
			"registry": "quay.io"
		}
	}
	`
	alreadyMigratedJSON = `
	{
		"controller": {
			"replicas": 3,
			"service": {
				"enabled": false
			}
		},
		"global": {
			"controller": {
				"tempReplicas": 2,
				"useProxyProtocol": false
			},
			"migration": {
				"enabled": false
			}
		},
		"image": {
			"registry": "quay.io"
		}
	}
	`
)

func Test_ConfigMap_GetDesiredState(t *testing.T) {
	testCases := []struct {
		name                            string
		configMapConfig                 ConfigMapConfig
		configMapValues                 ConfigMapValues
		ingressControllerReleasePresent bool
		expectedConfigMaps              []*corev1.ConfigMap
	}{
		{
			name: "case 0: basic match",
			configMapConfig: ConfigMapConfig{
				ClusterID:      "5xchu",
				GuestAPIDomain: "5xchu.aws.giantswarm.io",
				Namespaces:     []string{},
			},
			configMapValues: ConfigMapValues{
				ClusterID:                         "5xchu",
				IngressControllerMigrationEnabled: true,
				IngressControllerUseProxyProtocol: true,
				Organization:                      "giantswarm",
				WorkerCount:                       3,
			},
			ingressControllerReleasePresent: false,
			expectedConfigMaps: []*corev1.ConfigMap{
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "cert-exporter-values",
						Namespace: metav1.NamespaceSystem,
						Labels: map[string]string{
							label.App:          "cert-exporter",
							label.Cluster:      "5xchu",
							label.ManagedBy:    "cluster-operator",
							label.Organization: "giantswarm",
							label.ServiceType:  "managed",
						},
					},
					Data: map[string]string{
						"values.json": "{\"namespace\":\"kube-system\"}",
					},
				},
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "nginx-ingress-controller-values",
						Namespace: metav1.NamespaceSystem,
						Labels: map[string]string{
							label.App:          "nginx-ingress-controller",
							label.Cluster:      "5xchu",
							label.ManagedBy:    "cluster-operator",
							label.Organization: "giantswarm",
							label.ServiceType:  "managed",
						},
					},
					Data: map[string]string{
						"values.json": basicMatchJSON,
					},
				},
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "kube-state-metrics-values",
						Namespace: metav1.NamespaceSystem,
						Labels: map[string]string{
							label.App:          "kube-state-metrics",
							label.Cluster:      "5xchu",
							label.ManagedBy:    "cluster-operator",
							label.Organization: "giantswarm",
							label.ServiceType:  "managed",
						},
					},
					Data: map[string]string{
						"values.json": "{\"image\":{\"registry\":\"quay.io\"}}",
					},
				},
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "net-exporter-values",
						Namespace: metav1.NamespaceSystem,
						Labels: map[string]string{
							label.App:          "net-exporter",
							label.Cluster:      "5xchu",
							label.ManagedBy:    "cluster-operator",
							label.Organization: "giantswarm",
							label.ServiceType:  "managed",
						},
					},
					Data: map[string]string{
						"values.json": "{\"namespace\":\"kube-system\"}",
					},
				},
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "node-exporter-values",
						Namespace: metav1.NamespaceSystem,
						Labels: map[string]string{
							label.App:          "node-exporter",
							label.Cluster:      "5xchu",
							label.ManagedBy:    "cluster-operator",
							label.Organization: "giantswarm",
							label.ServiceType:  "managed",
						},
					},
					Data: map[string]string{
						"values.json": "{\"image\":{\"registry\":\"quay.io\"}}",
					},
				},
			},
		},
		{
			name: "case 1: different worker count",
			configMapConfig: ConfigMapConfig{
				ClusterID:      "5xchu",
				GuestAPIDomain: "5xchu.aws.giantswarm.io",
				Namespaces:     []string{},
			},
			configMapValues: ConfigMapValues{
				ClusterID:                         "5xchu",
				Organization:                      "giantswarm",
				IngressControllerMigrationEnabled: true,
				IngressControllerUseProxyProtocol: true,
				WorkerCount:                       7,
			},
			ingressControllerReleasePresent: false,
			expectedConfigMaps: []*corev1.ConfigMap{
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "cert-exporter-values",
						Namespace: metav1.NamespaceSystem,
						Labels: map[string]string{
							label.App:          "cert-exporter",
							label.Cluster:      "5xchu",
							label.ManagedBy:    "cluster-operator",
							label.Organization: "giantswarm",
							label.ServiceType:  "managed",
						},
					},
					Data: map[string]string{
						"values.json": "{\"namespace\":\"kube-system\"}",
					},
				},
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "nginx-ingress-controller-values",
						Namespace: metav1.NamespaceSystem,
						Labels: map[string]string{
							label.App:          "nginx-ingress-controller",
							label.Cluster:      "5xchu",
							label.ManagedBy:    "cluster-operator",
							label.Organization: "giantswarm",
							label.ServiceType:  "managed",
						},
					},
					Data: map[string]string{
						"values.json": differentWorkerCountJSON,
					},
				},
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "kube-state-metrics-values",
						Namespace: metav1.NamespaceSystem,
						Labels: map[string]string{
							label.App:          "kube-state-metrics",
							label.Cluster:      "5xchu",
							label.ManagedBy:    "cluster-operator",
							label.Organization: "giantswarm",
							label.ServiceType:  "managed",
						},
					},
					Data: map[string]string{
						"values.json": "{\"image\":{\"registry\":\"quay.io\"}}",
					},
				},
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "net-exporter-values",
						Namespace: metav1.NamespaceSystem,
						Labels: map[string]string{
							label.App:          "net-exporter",
							label.Cluster:      "5xchu",
							label.ManagedBy:    "cluster-operator",
							label.Organization: "giantswarm",
							label.ServiceType:  "managed",
						},
					},
					Data: map[string]string{
						"values.json": "{\"namespace\":\"kube-system\"}",
					},
				},
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "node-exporter-values",
						Namespace: metav1.NamespaceSystem,
						Labels: map[string]string{
							label.App:          "node-exporter",
							label.Cluster:      "5xchu",
							label.ManagedBy:    "cluster-operator",
							label.Organization: "giantswarm",
							label.ServiceType:  "managed",
						},
					},
					Data: map[string]string{
						"values.json": "{\"image\":{\"registry\":\"quay.io\"}}",
					},
				},
			},
		},
		{
			name: "case 2: different ingress controller settings",
			configMapConfig: ConfigMapConfig{
				ClusterID:      "5xchu",
				GuestAPIDomain: "5xchu.aws.giantswarm.io",
				Namespaces:     []string{},
			},
			configMapValues: ConfigMapValues{
				ClusterID:                         "5xchu",
				IngressControllerMigrationEnabled: false,
				IngressControllerUseProxyProtocol: false,
				Organization:                      "giantswarm",
				WorkerCount:                       1,
			},
			ingressControllerReleasePresent: false,
			expectedConfigMaps: []*corev1.ConfigMap{
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "cert-exporter-values",
						Namespace: metav1.NamespaceSystem,
						Labels: map[string]string{
							label.App:          "cert-exporter",
							label.Cluster:      "5xchu",
							label.ManagedBy:    "cluster-operator",
							label.Organization: "giantswarm",
							label.ServiceType:  "managed",
						},
					},
					Data: map[string]string{
						"values.json": "{\"namespace\":\"kube-system\"}",
					},
				},
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "nginx-ingress-controller-values",
						Namespace: metav1.NamespaceSystem,
						Labels: map[string]string{
							label.App:          "nginx-ingress-controller",
							label.Cluster:      "5xchu",
							label.ManagedBy:    "cluster-operator",
							label.Organization: "giantswarm",
							label.ServiceType:  "managed",
						},
					},
					Data: map[string]string{
						"values.json": differentSettingsJSON,
					},
				},
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "kube-state-metrics-values",
						Namespace: metav1.NamespaceSystem,
						Labels: map[string]string{
							label.App:          "kube-state-metrics",
							label.Cluster:      "5xchu",
							label.ManagedBy:    "cluster-operator",
							label.Organization: "giantswarm",
							label.ServiceType:  "managed",
						},
					},
					Data: map[string]string{
						"values.json": "{\"image\":{\"registry\":\"quay.io\"}}",
					},
				},
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "net-exporter-values",
						Namespace: metav1.NamespaceSystem,
						Labels: map[string]string{
							label.App:          "net-exporter",
							label.Cluster:      "5xchu",
							label.ManagedBy:    "cluster-operator",
							label.Organization: "giantswarm",
							label.ServiceType:  "managed",
						},
					},
					Data: map[string]string{
						"values.json": "{\"namespace\":\"kube-system\"}",
					},
				},
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "node-exporter-values",
						Namespace: metav1.NamespaceSystem,
						Labels: map[string]string{
							label.App:          "node-exporter",
							label.Cluster:      "5xchu",
							label.ManagedBy:    "cluster-operator",
							label.Organization: "giantswarm",
							label.ServiceType:  "managed",
						},
					},
					Data: map[string]string{
						"values.json": "{\"image\":{\"registry\":\"quay.io\"}}",
					},
				},
			},
		},
		{
			name: "case 3: ingress controller already migrated",
			configMapConfig: ConfigMapConfig{
				ClusterID:      "5xchu",
				GuestAPIDomain: "5xchu.aws.giantswarm.io",
				Namespaces:     []string{},
			},
			configMapValues: ConfigMapValues{
				ClusterID:                         "5xchu",
				IngressControllerMigrationEnabled: true,
				Organization:                      "giantswarm",
				WorkerCount:                       3,
			},
			ingressControllerReleasePresent: true,
			expectedConfigMaps: []*corev1.ConfigMap{
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "cert-exporter-values",
						Namespace: metav1.NamespaceSystem,
						Labels: map[string]string{
							label.App:          "cert-exporter",
							label.Cluster:      "5xchu",
							label.ManagedBy:    "cluster-operator",
							label.Organization: "giantswarm",
							label.ServiceType:  "managed",
						},
					},
					Data: map[string]string{
						"values.json": "{\"namespace\":\"kube-system\"}",
					},
				},
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "nginx-ingress-controller-values",
						Namespace: metav1.NamespaceSystem,
						Labels: map[string]string{
							label.App:          "nginx-ingress-controller",
							label.Cluster:      "5xchu",
							label.ManagedBy:    "cluster-operator",
							label.Organization: "giantswarm",
							label.ServiceType:  "managed",
						},
					},
					Data: map[string]string{
						"values.json": alreadyMigratedJSON,
					},
				},
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "kube-state-metrics-values",
						Namespace: metav1.NamespaceSystem,
						Labels: map[string]string{
							label.App:          "kube-state-metrics",
							label.Cluster:      "5xchu",
							label.ManagedBy:    "cluster-operator",
							label.Organization: "giantswarm",
							label.ServiceType:  "managed",
						},
					},
					Data: map[string]string{
						"values.json": "{\"image\":{\"registry\":\"quay.io\"}}",
					},
				},
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "net-exporter-values",
						Namespace: metav1.NamespaceSystem,
						Labels: map[string]string{
							label.App:          "net-exporter",
							label.Cluster:      "5xchu",
							label.ManagedBy:    "cluster-operator",
							label.Organization: "giantswarm",
							label.ServiceType:  "managed",
						},
					},
					Data: map[string]string{
						"values.json": "{\"namespace\":\"kube-system\"}",
					},
				},
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "node-exporter-values",
						Namespace: metav1.NamespaceSystem,
						Labels: map[string]string{
							label.App:          "node-exporter",
							label.Cluster:      "5xchu",
							label.ManagedBy:    "cluster-operator",
							label.Organization: "giantswarm",
							label.ServiceType:  "managed",
						},
					},
					Data: map[string]string{
						"values.json": "{\"image\":{\"registry\":\"quay.io\"}}",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			helmClient := &helmMock{}
			if !tc.ingressControllerReleasePresent {
				helmClient.defaultError = microerror.Newf("No such release: nginx-ingress-controller")
			}

			c := Config{
				Tenant: &tenantMock{
					fakeTenantHelmClient: helmClient,
				},
				Logger:         microloggertest.New(),
				ProjectName:    "cluster-operator",
				RegistryDomain: "quay.io",
			}
			newService, err := New(c)
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}

			configMaps, err := newService.GetDesiredState(context.TODO(), tc.configMapConfig, tc.configMapValues)
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}

			if len(configMaps) != len(tc.expectedConfigMaps) {
				t.Fatal("expected", len(tc.expectedConfigMaps), "got", len(configMaps))
			}

			for _, expectedConfigMap := range tc.expectedConfigMaps {
				configMap, err := getConfigMapByNameAndNamespace(configMaps, expectedConfigMap.Name, expectedConfigMap.Namespace)
				if IsNotFound(err) {
					t.Fatalf("expected config map '%s' not found", expectedConfigMap.Name)
				} else if err != nil {
					t.Fatalf("expected nil, got %#v", err)
				}

				if !reflect.DeepEqual(configMap.ObjectMeta.Labels, expectedConfigMap.ObjectMeta.Labels) {
					t.Fatalf("expected config map labels %#v, got %#v", expectedConfigMap.ObjectMeta.Labels, configMap.ObjectMeta.Labels)
				}

				for expectedKey, expectedValues := range expectedConfigMap.Data {
					values, ok := configMap.Data[expectedKey]
					if !ok {
						t.Fatalf("expected key '%s' not found", expectedKey)
					}

					equalValues, err := compareJSON(expectedValues, values)
					if err != nil {
						t.Fatal("expected", nil, "got", err)
					}
					if !equalValues {
						t.Fatal("expected", expectedValues, "got", values)
					}
				}
			}
		})
	}
}

func Test_ConfigMap_setIngressControllerTempReplicas(t *testing.T) {
	testCases := []struct {
		name                 string
		workerCount          int
		expectedTempReplicas int
		errorMatcher         func(error) bool
	}{
		{
			name:                 "case 0: basic match",
			workerCount:          3,
			expectedTempReplicas: 2,
		},
		{
			name:                 "case 1: single node",
			workerCount:          1,
			expectedTempReplicas: 1,
		},
		{
			name:                 "case 2: large cluster",
			workerCount:          20,
			expectedTempReplicas: 10,
		},
		{
			name:                 "case 3: larger cluster",
			workerCount:          50,
			expectedTempReplicas: 25,
		},
		{
			name:         "case 4: 0 workers",
			workerCount:  0,
			errorMatcher: IsInvalidExecution,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempReplicas, err := setIngressControllerTempReplicas(tc.workerCount)

			switch {
			case err == nil && tc.errorMatcher == nil:
				// correct; carry on
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case !tc.errorMatcher(err):
				t.Fatalf("error == %#v, want matching", err)
			}

			if tempReplicas != tc.expectedTempReplicas {
				t.Fatal("expected", tc.expectedTempReplicas, "got", tempReplicas)
			}
		})
	}
}

func compareJSON(expectedJSON, valuesJSON string) (bool, error) {
	var err error

	expectedValues := make(map[string]interface{})
	err = json.Unmarshal([]byte(expectedJSON), &expectedValues)
	if err != nil {
		return false, microerror.Mask(err)
	}

	values := make(map[string]interface{})
	err = json.Unmarshal([]byte(valuesJSON), &values)
	if err != nil {
		return false, microerror.Mask(err)
	}

	return reflect.DeepEqual(expectedValues, values), nil
}