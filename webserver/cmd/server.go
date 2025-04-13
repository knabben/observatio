package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"os/user"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the webserver",
	Long: `Serve the webserver on a default 8080 port and consuming K8s from .kube/config. 
Check args to change default values.`,
	RunE: RunE,
}

func init() {
	rootCmd.AddCommand(serveCmd)
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	serveCmd.PersistentFlags().Int("port", 8080, "Webserver default listening HTTP port. Default 8080.")
	serveCmd.PersistentFlags().String("kubeconfig", fmt.Sprintf("%s/.kube/config", currentUser.HomeDir), "Kubernetes configuration file path. Default $HOME/.kube/config")
}

func RunE(cmd *cobra.Command, args []string) error {
	fmt.Println("server")
	return nil
}
