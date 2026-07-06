package models

// ProviderStatus represents whether an infrastructure provider is installed
// in the connected environment, and its version.
type ProviderStatus struct {
	// Installed is true when this provider's clusterctl inventory entry was found.
	Installed bool `json:"installed"`

	// Version is the installed provider version; empty when Installed is false.
	Version string `json:"version"`
}

// InfrastructureCapability represents which infrastructure providers are
// installed in the connected environment, derived from the clusterctl
// provider inventory.
type InfrastructureCapability struct {
	// Docker is the detection result for the Docker infrastructure provider.
	Docker ProviderStatus `json:"docker"`

	// VSphere is the detection result for the vSphere infrastructure provider.
	VSphere ProviderStatus `json:"vsphere"`
}
