package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestWithKubernetes(t *testing.T) {
	tests := []struct {
		name           string
		client         any
		config         any
		expectedClient any
		expectedConfig any
	}{
		{
			name:           "valid client and config",
			client:         "client-value",
			config:         "config-value",
			expectedClient: "client-value",
			expectedConfig: "config-value",
		},
		{
			name:           "nil client and valid config",
			client:         nil,
			config:         "config-value",
			expectedClient: nil,
			expectedConfig: "config-value",
		},
		{
			name:           "valid client and nil config",
			client:         "client-value",
			config:         nil,
			expectedClient: "client-value",
			expectedConfig: nil,
		},
		{
			name:           "nil client and nil config",
			client:         nil,
			config:         nil,
			expectedClient: nil,
			expectedConfig: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := WithKubernetes(tt.client, tt.config)
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				client := ctx.Value("client")
				if client != tt.expectedClient {
					t.Errorf("expected client=%v, got %v", tt.expectedClient, client)
				}
				config := ctx.Value("config")
				if config != tt.expectedConfig {
					t.Errorf("expected config=%v, got %v", tt.expectedConfig, config)
				}
			})

			router := mux.NewRouter()
			router.Use(middleware)
			router.Handle("/validate", handler)
			req := httptest.NewRequest(http.MethodGet, "/validate", nil)
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)
			if res.Code != http.StatusOK {
				t.Errorf("expected status code 200, got %d", res.Code)
			}
		})
	}
}
