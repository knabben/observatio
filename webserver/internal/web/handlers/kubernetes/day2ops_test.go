package kubernetes

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func Test_parseDay2OpsDetailQuery(t *testing.T) {
	cases := []struct {
		name      string
		query     url.Values
		wantGVR   schema.GroupVersionResource
		wantNS    string
		wantName  string
		wantError bool
	}{
		{
			name: "valid machine GVR",
			query: url.Values{
				"group": {"cluster.x-k8s.io"}, "version": {"v1beta1"}, "resource": {"machines"},
				"namespace": {"default"}, "name": {"worker-0"},
			},
			wantGVR:  schema.GroupVersionResource{Group: "cluster.x-k8s.io", Version: "v1beta1", Resource: "machines"},
			wantNS:   "default",
			wantName: "worker-0",
		},
		{
			name:      "missing name",
			query:     url.Values{"version": {"v1beta1"}, "resource": {"machines"}, "namespace": {"default"}},
			wantError: true,
		},
		{
			name:      "missing version and namespace",
			query:     url.Values{"resource": {"machines"}, "name": {"worker-0"}},
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
			gvr, ns, name, err := parseDay2OpsDetailQuery(tt.query)
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
