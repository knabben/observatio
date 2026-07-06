package kubernetes

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func Test_parseRawObjectQuery(t *testing.T) {
	cases := []struct {
		name      string
		query     url.Values
		wantGVR   schema.GroupVersionResource
		wantNS    string
		wantName  string
		wantError bool
	}{
		{
			name: "valid cluster GVR",
			query: url.Values{
				"group": {"cluster.x-k8s.io"}, "version": {"v1beta1"}, "resource": {"clusters"},
				"namespace": {"default"}, "name": {"capi-workload"},
			},
			wantGVR:  schema.GroupVersionResource{Group: "cluster.x-k8s.io", Version: "v1beta1", Resource: "clusters"},
			wantNS:   "default",
			wantName: "capi-workload",
		},
		{
			name:      "missing name",
			query:     url.Values{"version": {"v1beta1"}, "resource": {"clusters"}, "namespace": {"default"}},
			wantError: true,
		},
		{
			name:      "missing version and namespace",
			query:     url.Values{"resource": {"clusters"}, "name": {"c1"}},
			wantError: true,
		},
		{
			name:      "empty query",
			query:     url.Values{},
			wantError: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			gvr, ns, name, err := parseRawObjectQuery(tt.query)
			if tt.wantError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantGVR, gvr)
			assert.Equal(t, tt.wantNS, ns)
			assert.Equal(t, tt.wantName, name)
		})
	}
}
