package fetchers

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func selfSignedCertPEM(t *testing.T, notAfter time.Time) string {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "test-ca"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     notAfter,
		IsCA:         true,
	}
	der, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	require.NoError(t, err)

	var pemBuf []byte
	pemBuf = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	return base64.StdEncoding.EncodeToString(pemBuf)
}

func Test_parseTLSCertNotAfter(t *testing.T) {
	expiry := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	obj := &unstructured.Unstructured{Object: map[string]interface{}{
		"data": map[string]interface{}{"tls.crt": selfSignedCertPEM(t, expiry)},
	}}

	notAfter, err := parseTLSCertNotAfter(obj)
	require.NoError(t, err)
	assert.WithinDuration(t, expiry, notAfter, time.Second)
}

func Test_parseTLSCertNotAfter_MissingKey(t *testing.T) {
	obj := &unstructured.Unstructured{Object: map[string]interface{}{
		"data": map[string]interface{}{"tls.key": "some-key"},
	}}

	_, err := parseTLSCertNotAfter(obj)
	assert.Error(t, err)
}

func Test_parseTLSCertNotAfter_InvalidBase64(t *testing.T) {
	obj := &unstructured.Unstructured{Object: map[string]interface{}{
		"data": map[string]interface{}{"tls.crt": "not-valid-base64!!!"},
	}}

	_, err := parseTLSCertNotAfter(obj)
	assert.Error(t, err)
}
