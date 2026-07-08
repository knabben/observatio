package day2ops

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ComputeCertExpiryRisk_WithinWindow(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	notAfter := now.Add(10 * 24 * time.Hour) // expires in 10 days, inside the 30-day window

	risk := ComputeCertExpiryRisk(ObjectRef{Name: "prod-1"}, "prod-1-ca", notAfter, now, DefaultCertExpiryWarningWindow)

	require.NotNil(t, risk)
	assert.Equal(t, RiskCertExpiry, risk.Kind)
	assert.Equal(t, RiskCheckEvaluated, risk.CheckStatus)
	assert.Contains(t, risk.Detail, "prod-1-ca")
}

func Test_ComputeCertExpiryRisk_OutsideWindow(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	notAfter := now.Add(365 * 24 * time.Hour) // expires in a year - no warning

	risk := ComputeCertExpiryRisk(ObjectRef{Name: "prod-1"}, "prod-1-ca", notAfter, now, DefaultCertExpiryWarningWindow)
	assert.Nil(t, risk)
}

func Test_ComputeCertExpiryRisk_ExactBoundary(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	notAfter := now.Add(DefaultCertExpiryWarningWindow) // exactly at the boundary

	risk := ComputeCertExpiryRisk(ObjectRef{Name: "prod-1"}, "prod-1-ca", notAfter, now, DefaultCertExpiryWarningWindow)
	require.NotNil(t, risk, "boundary case must warn, not silently miss the last eligible day")
}

func Test_ComputeClusterCertRisks_NotEvaluableOnFetchError(t *testing.T) {
	risks := ComputeClusterCertRisks(ObjectRef{Name: "prod-1"}, nil, assertError{}, time.Now(), DefaultCertExpiryWarningWindow)
	require.Len(t, risks, 1)
	assert.Equal(t, RiskCheckNotEvaluable, risks[0].CheckStatus)
}

func Test_ComputeClusterCertRisks_EvaluatesEachExpiry(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	expiries := []CertExpiry{
		{SecretName: "prod-1-ca", NotAfter: now.Add(5 * 24 * time.Hour)},
		{SecretName: "prod-1-etcd", NotAfter: now.Add(400 * 24 * time.Hour)},
	}

	risks := ComputeClusterCertRisks(ObjectRef{Name: "prod-1"}, expiries, nil, now, DefaultCertExpiryWarningWindow)
	require.Len(t, risks, 1)
	assert.Contains(t, risks[0].Detail, "prod-1-ca")
}

type assertError struct{}

func (assertError) Error() string { return "forbidden" }
