package service

import (
	"github.com/giantswarm/versionbundle"

	awsv1 "github.com/giantswarm/cluster-operator/service/controller/aws/v1"
	awsv2 "github.com/giantswarm/cluster-operator/service/controller/aws/v2"
	awsv3 "github.com/giantswarm/cluster-operator/service/controller/aws/v3"
	azurev1 "github.com/giantswarm/cluster-operator/service/controller/azure/v1"
	azurev2 "github.com/giantswarm/cluster-operator/service/controller/azure/v2"
	azurev3 "github.com/giantswarm/cluster-operator/service/controller/azure/v3"
	kvmv1 "github.com/giantswarm/cluster-operator/service/controller/kvm/v1"
	kvmv2 "github.com/giantswarm/cluster-operator/service/controller/kvm/v2"
	kvmv3 "github.com/giantswarm/cluster-operator/service/controller/kvm/v3"
)

func NewVersionBundles() []versionbundle.Bundle {
	var versionBundles []versionbundle.Bundle

	versionBundles = append(versionBundles, awsv1.VersionBundle())
	versionBundles = append(versionBundles, azurev1.VersionBundle())
	versionBundles = append(versionBundles, kvmv1.VersionBundle())

	versionBundles = append(versionBundles, awsv2.VersionBundle())
	versionBundles = append(versionBundles, azurev2.VersionBundle())
	versionBundles = append(versionBundles, kvmv2.VersionBundle())

	versionBundles = append(versionBundles, awsv3.VersionBundle())
	versionBundles = append(versionBundles, azurev3.VersionBundle())
	versionBundles = append(versionBundles, kvmv3.VersionBundle())

	return versionBundles
}