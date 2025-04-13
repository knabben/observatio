package cmd

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/knabben/observatio/webserver/internal/infra/kubernetes"
	"github.com/knabben/observatio/webserver/internal/web"
)

var (
	address     string
	development bool
	timeout     time.Duration = 15 * time.Second
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the webserver",
	Long: `Serve the webserver on a default 8080 port and consuming K8s from .kube/config. 
Check args to change default values.`,
	RunE: RunE,
}

func init() {
	log.SetLogger(zap.New())

	serveCmd.PersistentFlags().StringVar(&address, "address", ":8080", "Webserver default listening HTTP port. Default 8080.")
	serveCmd.PersistentFlags().BoolVar(&development, "dev", false, "Development mode, no static hosting. Default false")

	rootCmd.AddCommand(serveCmd)
}

func RunE(cmd *cobra.Command, args []string) error {
	router := mux.NewRouter()
	client, config, err := kubernetes.NewClient()
	if err != nil {
		return err
	}
	router.Use(web.WithKubernetes(client, config))
	web.DefaultHandlers(router, development)

	ctx := context.Background()
	log.FromContext(ctx).Info("Listening on " + address)
	srv := &http.Server{Handler: router, Addr: address, WriteTimeout: timeout, ReadTimeout: timeout}
	return srv.ListenAndServe()
}
