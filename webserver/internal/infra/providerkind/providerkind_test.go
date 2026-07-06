package providerkind

import "testing"

func TestFromKind(t *testing.T) {
	cases := []struct {
		kind string
		want string
	}{
		{"DockerCluster", Docker},
		{"DockerMachine", Docker},
		{"VSphereCluster", VSphere},
		{"VSphereMachine", VSphere},
		{"AWSCluster", Unknown},
		{"", Unknown},
	}

	for _, c := range cases {
		if got := FromKind(c.kind); got != c.want {
			t.Errorf("FromKind(%q) = %q, want %q", c.kind, got, c.want)
		}
	}
}
