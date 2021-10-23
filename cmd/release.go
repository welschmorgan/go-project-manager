package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/welschmorgan/go-release-manager/release"
)

var ReleaseCmd = &cobra.Command{
	Use:   "release <major|minor|build|revision|preRelease|buildMetaTag> [OPTIONS...]",
	Short: "Release all projects included in this workspace",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires exactly 1 argument: release_type")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return release.DoRelease(args[0])
	},
}
