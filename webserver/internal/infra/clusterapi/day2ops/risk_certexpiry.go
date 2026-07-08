package day2ops

import (
	"fmt"
	"time"
)

// DefaultCertExpiryWarningWindow is the default lead time before certificate expiry that
// triggers a warning (FR-008; SC-003's default of 30 days).
const DefaultCertExpiryWarningWindow = 30 * 24 * time.Hour

// CertExpiry is the parsed NotAfter for one CAPI-managed cert Secret. Only the expiry timestamp
// is ever kept — raw certificate/key bytes are read solely to extract this and are discarded
// (plan.md Constraints).
type CertExpiry struct {
	SecretName string
	NotAfter   time.Time
}

// ComputeCertExpiryRisk flags a cert nearing expiry within the warning window. Returns nil when
// the cert isn't close enough to expiry to warrant a warning.
func ComputeCertExpiryRisk(objectRef ObjectRef, secretName string, notAfter time.Time, now time.Time, window time.Duration) *RiskWarning {
	if notAfter.Sub(now) > window {
		return nil
	}
	return &RiskWarning{
		ObjectRef:   objectRef,
		Kind:        RiskCertExpiry,
		Detail:      fmt.Sprintf("%s expires %s", secretName, notAfter.Format("2006-01-02")),
		CheckStatus: RiskCheckEvaluated,
	}
}

// ComputeClusterCertRisks maps a cluster's fetched cert expiries into risk warnings. When
// fetchErr is non-nil (the check could not be performed, e.g. an RBAC Forbidden), it returns a
// single not-evaluable warning instead of silently omitting the risk category (FR-018).
func ComputeClusterCertRisks(objectRef ObjectRef, expiries []CertExpiry, fetchErr error, now time.Time, window time.Duration) []RiskWarning {
	if fetchErr != nil {
		return []RiskWarning{{ObjectRef: objectRef, Kind: RiskCertExpiry, CheckStatus: RiskCheckNotEvaluable}}
	}
	var risks []RiskWarning
	for _, e := range expiries {
		if risk := ComputeCertExpiryRisk(objectRef, e.SecretName, e.NotAfter, now, window); risk != nil {
			risks = append(risks, *risk)
		}
	}
	return risks
}
