package chartoperator

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/apiextensions/pkg/clientset/versioned/fake"
	"github.com/giantswarm/apprclient/apprclienttest"
	"github.com/giantswarm/helmclient"
	"github.com/giantswarm/helmclient/helmclienttest"
	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/spf13/afero"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientgofake "k8s.io/client-go/kubernetes/fake"

	"github.com/giantswarm/cluster-operator/pkg/cluster"
)

func Test_Chart_GetCurrentState(t *testing.T) {
	testCases := []struct {
		name           string
		obj            interface{}
		releaseContent *helmclient.ReleaseContent
		releaseHistory *helmclient.ReleaseHistory
		helmError      error
		expectedState  ResourceState
		expectedError  bool
	}{
		{
			name: "case 0: basic match",
			obj: v1alpha1.ClusterGuestConfig{
				DNSZone: "5xchu.aws.giantswarm.io",
				ID:      "5xchu",
				Owner:   "giantswarm",
			},
			releaseContent: &helmclient.ReleaseContent{
				Name:   "chart-operator",
				Status: "DEPLOYED",
				Values: map[string]interface{}{
					"key": "value",
				},
			},
			releaseHistory: &helmclient.ReleaseHistory{
				Name:    "chart-operator",
				Version: "0.1.2",
			},
			expectedState: ResourceState{
				ChartName:      "chart-operator-chart",
				ReleaseName:    "chart-operator",
				ReleaseStatus:  "DEPLOYED",
				ReleaseVersion: "0.1.2",
			},
		},
		{
			name: "case 1: different release status",
			obj: v1alpha1.ClusterGuestConfig{
				DNSZone: "5xchu.aws.giantswarm.io",
				ID:      "5xchu",
				Owner:   "giantswarm",
			},
			releaseContent: &helmclient.ReleaseContent{
				Name:   "chart-operator",
				Status: "FAILED",
				Values: map[string]interface{}{
					"key": "value",
				},
			},
			releaseHistory: &helmclient.ReleaseHistory{
				Name:    "chart-operator",
				Version: "0.1.2",
			},
			expectedState: ResourceState{
				ChartName:      "chart-operator-chart",
				ReleaseName:    "chart-operator",
				ReleaseStatus:  "FAILED",
				ReleaseVersion: "0.1.2",
			},
		},
		{
			name: "case 2: state is empty when release does not exist",
			obj: v1alpha1.ClusterGuestConfig{
				DNSZone: "5xchu.aws.giantswarm.io",
				ID:      "5xchu",
				Owner:   "giantswarm",
			},
			releaseContent: &helmclient.ReleaseContent{},
			releaseHistory: &helmclient.ReleaseHistory{},
			helmError:      fmt.Errorf("No such release: chart-operator"),
			expectedState:  ResourceState{},
			expectedError:  false,
		},
		{
			name: "case 3: unexpected error",
			obj: v1alpha1.ClusterGuestConfig{
				DNSZone: "5xchu.aws.giantswarm.io",
				ID:      "5xchu",
				Owner:   "giantswarm",
			},
			releaseContent: &helmclient.ReleaseContent{},
			releaseHistory: &helmclient.ReleaseHistory{},
			helmError:      fmt.Errorf("Unexpected error"),
			expectedError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var helmClient helmclient.Interface
			{
				c := helmclienttest.Config{
					DefaultReleaseContent: tc.releaseContent,
					DefaultReleaseHistory: tc.releaseHistory,
					DefaultError:          tc.helmError,
				}
				helmClient = helmclienttest.New(c)
			}

			c := Config{
				ApprClient: apprclienttest.New(apprclienttest.Config{}),
				BaseClusterConfig: cluster.Config{
					ClusterID: "test-cluster",
				},
				ClusterIPRange: "172.31.0.0/16",
				Fs:             afero.NewMemMapFs(),
				G8sClient:      fake.NewSimpleClientset(),
				K8sClient:      clientgofake.NewSimpleClientset(),
				Logger:         microloggertest.New(),
				ProjectName:    "cluster-operator",
				RegistryDomain: "quay.io",
				Tenant: &tenantMock{
					fakeTenantHelmClient: helmClient,
				},
				ToClusterGuestConfigFunc: func(v interface{}) (v1alpha1.ClusterGuestConfig, error) {
					return v.(v1alpha1.ClusterGuestConfig), nil
				},
				ToClusterObjectMetaFunc: func(v interface{}) (metav1.ObjectMeta, error) {
					return metav1.ObjectMeta{}, nil
				},
			}

			r, err := New(c)
			if err != nil {
				t.Fatalf("error == %#v, want nil", err)
			}

			result, err := r.GetCurrentState(context.TODO(), tc.obj)
			switch {
			case err != nil && !tc.expectedError:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.expectedError:
				t.Fatalf("error == nil, want non-nil")
			}

			if result != nil && !reflect.DeepEqual(tc.expectedState, ResourceState{}) {
				chartState, err := toResourceState(result)
				if err != nil {
					t.Fatalf("error == %#v, want nil", err)
				}

				if !reflect.DeepEqual(chartState, tc.expectedState) {
					t.Fatalf("ChartState == %q, want %q", chartState, tc.expectedState)
				}
			}
		})
	}
}