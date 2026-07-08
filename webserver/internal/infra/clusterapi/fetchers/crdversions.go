package fetchers

import (
	"context"

	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/day2ops"
)

// FetchCRDVersionInfo reads a CustomResourceDefinition's served versions and stored versions,
// backing the version-skew risk check (research.md R6 in specs/006-day2-ops-dashboard/).
func FetchCRDVersionInfo(ctx context.Context, clientset *apiextensionsclientset.Clientset, crdName string) (day2ops.CRDVersionInfo, error) {
	crd, err := clientset.ApiextensionsV1().CustomResourceDefinitions().Get(ctx, crdName, metav1.GetOptions{})
	if err != nil {
		return day2ops.CRDVersionInfo{}, err
	}

	var served []string
	for _, v := range crd.Spec.Versions {
		if v.Served {
			served = append(served, v.Name)
		}
	}
	return day2ops.CRDVersionInfo{
		Name:           crdName,
		ServedVersions: served,
		StoredVersions: crd.Status.StoredVersions,
	}, nil
}
