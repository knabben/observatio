// Package providerkind derives a normalized infrastructure provider name from
// the standard CAPI infrastructureRef.Kind field.
package providerkind

const (
	// Docker identifies a Docker (CAPD) infrastructure provider.
	Docker = "docker"

	// VSphere identifies a vSphere (CAPV) infrastructure provider.
	VSphere = "vsphere"

	// Unknown identifies an infrastructure reference kind that isn't a
	// recognized/supported provider.
	Unknown = "unknown"
)

// FromKind maps a Cluster/Machine infrastructureRef.Kind to a normalized
// provider name: "docker", "vsphere", or "unknown".
func FromKind(kind string) string {
	switch kind {
	case "DockerCluster", "DockerMachine":
		return Docker
	case "VSphereCluster", "VSphereMachine":
		return VSphere
	default:
		return Unknown
	}
}
