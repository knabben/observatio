package fetchers

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/day2ops"
)

var secretGVR = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "secrets"}

// CAPICertSecretSuffixes are the standard CAPI-managed cert Secret name suffixes for a cluster
// (research.md R4 in specs/006-day2-ops-dashboard/): "<cluster>-ca", "<cluster>-etcd",
// "<cluster>-proxy".
var CAPICertSecretSuffixes = []string{"ca", "etcd", "proxy"}

// FetchClusterCertExpiries reads the `<cluster>-ca`/`-etcd`/`-proxy` Secrets for a cluster and
// parses each `tls.crt`'s NotAfter. A Secret that doesn't exist is skipped, not an error — not
// every provider populates all three. A non-NotFound error (e.g. an RBAC Forbidden) is returned so
// the caller can distinguish "checked, nothing found" from "could not check" (FR-018).
func FetchClusterCertExpiries(ctx context.Context, dyn dynamic.Interface, namespace, clusterName string) ([]day2ops.CertExpiry, error) {
	var expiries []day2ops.CertExpiry
	for _, suffix := range CAPICertSecretSuffixes {
		secretName := fmt.Sprintf("%s-%s", clusterName, suffix)
		obj, err := dyn.Resource(secretGVR).Namespace(namespace).Get(ctx, secretName, metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				continue
			}
			return nil, err
		}
		notAfter, err := parseTLSCertNotAfter(obj)
		if err != nil {
			continue
		}
		expiries = append(expiries, day2ops.CertExpiry{SecretName: secretName, NotAfter: notAfter})
	}
	return expiries, nil
}

// parseTLSCertNotAfter extracts the `NotAfter` timestamp from a Secret's `data["tls.crt"]` field.
func parseTLSCertNotAfter(obj *unstructured.Unstructured) (time.Time, error) {
	dataMap, found, err := unstructured.NestedStringMap(obj.Object, "data")
	if err != nil || !found {
		return time.Time{}, fmt.Errorf("secret has no data field")
	}
	encoded, ok := dataMap["tls.crt"]
	if !ok {
		return time.Time{}, fmt.Errorf("secret has no tls.crt key")
	}
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return time.Time{}, err
	}
	block, _ := pem.Decode(raw)
	if block == nil {
		return time.Time{}, fmt.Errorf("tls.crt is not valid PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return time.Time{}, err
	}
	return cert.NotAfter, nil
}
