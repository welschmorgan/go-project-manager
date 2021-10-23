package cmd

import (
	"github.com/spf13/cobra"
	"github.com/welschmorgan/go-release-manager/api"
	"github.com/welschmorgan/go-release-manager/config"
)

var APICmd = &cobra.Command{
	Use:   "api [OPTIONS...]",
	Short: "Run the API server",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := api.NewAPIServer(config.Get().API.ListenAddr)
		s.Serve()
		return nil
	},
}
