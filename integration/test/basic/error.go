package basic

import (
	"github.com/giantswarm/microerror"
)

var emptyChartConfigListError = microerror.New("empty chart config list")

// IsEmptyChartConfigList asserts emptyChartConfigListError.
func IsEmptyChartConfigList(err error) bool {
	return microerror.Cause(err) == emptyChartConfigListError
}

var releaseStatusNotMatchingError = microerror.New("release status not matching")

// IsReleaseStatusNotMatching asserts releaseStatusNotMatchingError
func IsReleaseStatusNotMatching(err error) bool {
	return microerror.Cause(err) == releaseStatusNotMatchingError
}