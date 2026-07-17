package cmd

import (
	"context"
	"net/http"
	"os"

	gh "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
	mcpaggregator "github.com/knabben/observatio/webserver/internal/infra/mcp"
	"github.com/knabben/observatio/webserver/internal/web/handlers"
	"github.com/spf13/cobra"

	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/knabben/observatio/webserver/internal/web"
)

var (
	address           string
	development       bool
	toolSourcesConfig string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the webserver",
	Long: `Serve the webserver on a default 8080 port and consuming K8s from .kube/config. 
Check args to change default values.`,
	RunE: RunE,
}

func init() {
	log.SetLogger(zap.New(
		zap.UseDevMode(development),
		zap.WriteTo(os.Stdout),
		zap.JSONEncoder(),
	))

	serveCmd.PersistentFlags().StringVar(&address, "address", ":8080", "Webserver default listening HTTP port. Default 8080.")
	serveCmd.PersistentFlags().BoolVar(&development, "dev", false, "Development mode, no static hosting. Default false")
	serveCmd.PersistentFlags().StringVar(&toolSourcesConfig, "tool-sources-config", "",
		"Path to a YAML file registering external MCP tool sources for the AI assistant (see "+
			"specs/009-mcp-server-client-aggregator/contracts/tool-sources-config.md). Unset means no "+
			"external sources - only the built-in kubectl capability is offered. Also settable via "+
			"the TOOL_SOURCES_CONFIG env var.")

	rootCmd.AddCommand(serveCmd)
}

// resolveToolSourcesConfig lets TOOL_SOURCES_CONFIG override the --tool-sources-config flag at
// deploy time, matching the ANTHROPIC_MODEL env-override convention already used in
// internal/infra/llm/observation.go.
func resolveToolSourcesConfig() string {
	if v := os.Getenv("TOOL_SOURCES_CONFIG"); v != "" {
		return v
	}
	return toolSourcesConfig
}

func RunE(cmd *cobra.Command, args []string) error {
	client, config, err := clusterapi.NewClient()
	if err != nil {
		return err
	}

	// The tool source aggregator is built once, here, and shared across every WebSocket chat
	// connection - not rebuilt per connection - since registering a source can involve a real MCP
	// handshake (specs/009-mcp-server-client-aggregator, research.md R2).
	ctx := context.Background()
	localSource, err := mcpaggregator.NewLocalToolSource(ctx)
	if err != nil {
		return err
	}
	externalSources, err := mcpaggregator.BuildExternalSources(ctx, resolveToolSourcesConfig())
	if err != nil {
		return err
	}

	sources := []mcpaggregator.ToolSource{localSource}
	for _, src := range externalSources {
		src.StartHealthChecking(ctx)
		sources = append(sources, src)
	}
	aggregator := mcpaggregator.NewAggregator(sources...)

	router := mux.NewRouter()
	router.Use(web.WithKubernetes(client, config))
	router.Use(web.WithLogger())
	handlers.DefaultHandlers(router, development, aggregator)

	allowDomain := gh.AllowedOrigins([]string{"*"})
	allowMethods := gh.AllowedMethods([]string{"GET", "OPTIONS", "POST"})
	allowContentType := gh.AllowedHeaders([]string{"Content-Type", "Authorization"})

	log.FromContext(context.Background()).Info("Starting server, listening on " + address)
	return http.ListenAndServe(address, gh.CORS(allowDomain, allowMethods, allowContentType)(router))
}
