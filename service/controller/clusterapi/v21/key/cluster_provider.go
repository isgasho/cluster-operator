package key

import (
	"encoding/json"

	g8sv1alpha1 "github.com/giantswarm/apiextensions/pkg/apis/cluster/v1alpha1"
	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/runtime"
	cmav1alpha1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

func clusterProviderSpec(cluster cmav1alpha1.Cluster) g8sv1alpha1.AWSClusterSpec {
	spec, err := g8sClusterSpecFromCMAClusterSpec(cluster.Spec.ProviderSpec)
	if err != nil {
		panic(err)
	}
	return spec
}

func clusterProviderStatus(cluster cmav1alpha1.Cluster) g8sv1alpha1.AWSClusterStatus {
	return g8sClusterStatusFromCMAClusterStatus(cluster.Status.ProviderStatus)
}

func g8sClusterSpecFromCMAClusterSpec(cmaSpec cmav1alpha1.ProviderSpec) (g8sv1alpha1.AWSClusterSpec, error) {
	if cmaSpec.Value == nil {
		return g8sv1alpha1.AWSClusterSpec{}, microerror.Maskf(notFoundError, "provider spec extension for AWS not found")
	}

	var g8sSpec g8sv1alpha1.AWSClusterSpec
	{
		if len(cmaSpec.Value.Raw) == 0 {
			return g8sSpec, nil
		}

		err := json.Unmarshal(cmaSpec.Value.Raw, &g8sSpec)
		if err != nil {
			return g8sv1alpha1.AWSClusterSpec{}, microerror.Mask(err)
		}
	}

	return g8sSpec, nil
}

func g8sClusterCommonStatusFromCMAClusterStatus(cmaStatus *runtime.RawExtension) g8sv1alpha1.CommonClusterStatus {
	// Whatever provider status we unmarshal here, the wrapper should extract the
	// common cluster status types that way.
	type wrapper struct {
		Cluster g8sv1alpha1.CommonClusterStatus `json:"cluster" yaml:"cluster"`
	}

	var w wrapper
	{
		if cmaStatus == nil {
			return w.Cluster
		}

		if len(cmaStatus.Raw) == 0 {
			return w.Cluster
		}

		err := json.Unmarshal(cmaStatus.Raw, &w)
		if err != nil {
			panic(err)
		}
	}

	return w.Cluster
}

func g8sClusterStatusFromCMAClusterStatus(cmaStatus *runtime.RawExtension) g8sv1alpha1.AWSClusterStatus {
	var g8sStatus g8sv1alpha1.AWSClusterStatus
	{
		if cmaStatus == nil {
			return g8sStatus
		}

		if len(cmaStatus.Raw) == 0 {
			return g8sStatus
		}

		err := json.Unmarshal(cmaStatus.Raw, &g8sStatus)
		if err != nil {
			panic(err)
		}
	}

	return g8sStatus
}

func withG8sClusterStatusToCMAClusterStatus(cluster cmav1alpha1.Cluster, status g8sv1alpha1.AWSClusterStatus) cmav1alpha1.Cluster {
	var err error

	if cluster.Status.ProviderStatus == nil {
		cluster.Status.ProviderStatus = &runtime.RawExtension{}
	}

	cluster.Status.ProviderStatus.Raw, err = json.Marshal(status)
	if err != nil {
		panic(err)
	}

	return cluster
}