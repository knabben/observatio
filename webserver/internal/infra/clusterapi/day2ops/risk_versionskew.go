package day2ops

import (
	"fmt"
	"strings"
)

// CRDVersionInfo captures the minimal apiextensions data needed for the version-skew check.
type CRDVersionInfo struct {
	Name           string
	ServedVersions []string
	StoredVersions []string
}

// ComputeVersionSkewRisk flags a CRD with objects stored under a version the API server no longer
// serves — a directly observable upgrade hazard (research.md R6). This is a more reliably
// determinable signal than comparing a provider's release version against its CRD version, which
// use unrelated versioning schemes and can't be compared correctly in general.
func ComputeVersionSkewRisk(objectRef ObjectRef, info CRDVersionInfo) *RiskWarning {
	served := make(map[string]bool, len(info.ServedVersions))
	for _, v := range info.ServedVersions {
		served[v] = true
	}
	var stale []string
	for _, v := range info.StoredVersions {
		if !served[v] {
			stale = append(stale, v)
		}
	}
	if len(stale) == 0 {
		return nil
	}
	return &RiskWarning{
		ObjectRef: objectRef,
		Kind:      RiskVersionSkew,
		Detail: fmt.Sprintf(
			"%s has objects stored under version(s) %s, which are no longer served",
			info.Name, strings.Join(stale, ", "),
		),
		CheckStatus: RiskCheckEvaluated,
	}
}
