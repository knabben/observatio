package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/knabben/observatio/webserver/internal/infra/models"
	"github.com/knabben/observatio/webserver/internal/infra/providerkind"
)

func Test_resolveInfraProvider(t *testing.T) {
	both := models.InfrastructureCapability{
		Docker:  models.ProviderStatus{Installed: true, Version: "v1.10.10"},
		VSphere: models.ProviderStatus{Installed: true, Version: "v1.12.0"},
	}
	dockerOnly := models.InfrastructureCapability{
		Docker: models.ProviderStatus{Installed: true, Version: "v1.10.10"},
	}
	vsphereOnly := models.InfrastructureCapability{
		VSphere: models.ProviderStatus{Installed: true, Version: "v1.12.0"},
	}
	neither := models.InfrastructureCapability{}

	cases := []struct {
		name         string
		requested    string
		capability   models.InfrastructureCapability
		wantProvider string
		wantOk       bool
	}{
		{"auto-selects docker when both installed", "", both, providerkind.Docker, true},
		{"auto-selects docker when only docker installed", "", dockerOnly, providerkind.Docker, true},
		{"auto-selects vsphere when only vsphere installed", "", vsphereOnly, providerkind.VSphere, true},
		{"auto-select with neither installed returns empty but ok", "", neither, "", true},
		{"explicit docker request", "docker", vsphereOnly, providerkind.Docker, true},
		{"explicit vsphere request", "vsphere", dockerOnly, providerkind.VSphere, true},
		{"unrecognized provider is rejected", "aws", both, "", false},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			provider, ok := resolveInfraProvider(tt.requested, tt.capability)
			assert.Equal(t, tt.wantOk, ok)
			assert.Equal(t, tt.wantProvider, provider)
		})
	}
}
