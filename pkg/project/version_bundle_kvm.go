package project

import (
	"github.com/giantswarm/versionbundle"
)

var versionBundleKVM = versionbundle.Bundle{
	Changelogs: []versionbundle.Changelog{
		{
			Component:   "cluster-operator",
			Description: "TODO",
			Kind:        versionbundle.KindChanged,
			URLs: []string{
				"",
			},
		},
	},
	Components: []versionbundle.Component{
		{
			Name:    "kube-state-metrics",
			Version: "1.9.2",
		},
		{
			Name:    "nginx-ingress-controller",
			Version: "0.28.0",
		},
		{
			Name:    "node-exporter",
			Version: "0.18.1",
		},
		{
			Name:    "coredns",
			Version: "1.6.5",
		},
		{
			Name:    "metrics-server",
			Version: "0.3.3",
		},
		{
			Name:    "net-exporter",
			Version: "1.6.0",
		},
	},
	Name:     "cluster-operator",
	Provider: "kvm",
	Version:  BundleVersion(),
}
