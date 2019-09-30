package v20

import (
	"github.com/giantswarm/versionbundle"
)

func VersionBundle() versionbundle.Bundle {
	return versionbundle.Bundle{
		Changelogs: []versionbundle.Changelog{
			{
				Component:   "cluster-operator",
				Description: "Add internal Kubernetes API domain into API certificate alternative names.",
				Kind:        versionbundle.KindChanged,
			},
			{
				Component:   "chart-operator",
				Description: "Install chart-operator from default catalog.",
				Kind:        versionbundle.KindChanged,
			},
		},
		Components: []versionbundle.Component{
			{
				Name:    "coredns",
				Version: "1.6.2",
			},
			{
				Name:    "kube-state-metrics",
				Version: "1.7.2",
			},
			{
				Name:    "nginx-ingress-controller",
				Version: "0.25.1",
			},
			{
				Name:    "node-exporter",
				Version: "0.18.0",
			},
			{
				Name:    "metrics-server",
				Version: "0.3.1",
			},
		},
		Name:    "cluster-operator",
		Version: "0.20.0",
	}
}