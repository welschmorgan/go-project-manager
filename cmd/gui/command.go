package version

import (
	"github.com/spf13/cobra"
	"github.com/webview/webview"
	"github.com/welschmorgan/go-release-manager/config"
)

var projects []*config.Project

var Command = &cobra.Command{
	Use:   "gui",
	Short: "Interface to show workspace",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		debug := true
		w := webview.New(debug)
		defer w.Destroy()
		w.SetTitle("Minimal webview example")
		w.SetSize(800, 600, webview.HintNone)
		w.Navigate("https://en.m.wikipedia.org/wiki/Main_Page")
		w.Run()
		return nil
	},
}
