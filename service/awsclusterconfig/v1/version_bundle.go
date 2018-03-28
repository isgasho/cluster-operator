package v1

import (
	"time"

	"github.com/giantswarm/versionbundle"
)

func VersionBundle() versionbundle.Bundle {
	return versionbundle.Bundle{
		Changelogs: []versionbundle.Changelog{
			{
				Component:   "Cluster Operator",
				Description: "Initial version for AWS",
				Kind:        "added",
			},
		},
		Components: []versionbundle.Component{
			{
				Name:    "aws-operator",
				Version: "1.0.0",
			},
		},
		Dependencies: []versionbundle.Dependency{},
		Deprecated:   false,
		Name:         "cluster-operator",
		Time:         time.Date(2018, time.March, 27, 12, 00, 0, 0, time.UTC),
		Version:      "0.1.0",
		WIP:          true,
	}
}
