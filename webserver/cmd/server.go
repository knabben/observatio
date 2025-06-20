package cmd

import (
	"context"
	"net/http"
	"os"

	gh "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
	"github.com/knabben/observatio/webserver/internal/web/handlers"
	"github.com/spf13/cobra"

	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/knabben/observatio/webserver/internal/web"
)

var (
	address     string
	development bool
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

	rootCmd.AddCommand(serveCmd)
}

func RunE(cmd *cobra.Command, args []string) error {
	client, config, err := clusterapi.NewClient()
	if err != nil {
		return err
	}

	router := mux.NewRouter()
	router.Use(web.WithKubernetes(client, config))
	router.Use(web.WithLogger())
	handlers.DefaultHandlers(router, development)

	allowDomain := gh.AllowedOrigins([]string{"*"})
	allowMethods := gh.AllowedMethods([]string{"GET", "OPTIONS", "POST"})
	allowContentType := gh.AllowedHeaders([]string{"Content-Type", "Authorization"})

	log.FromContext(context.Background()).Info("Starting server, listening on " + address)
	return http.ListenAndServe(address, gh.CORS(allowDomain, allowMethods, allowContentType)(router))
}
