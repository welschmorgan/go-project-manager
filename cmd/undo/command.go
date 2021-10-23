package undo

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/log"
	"github.com/welschmorgan/go-release-manager/release"
	"github.com/welschmorgan/go-release-manager/ui"
	"github.com/welschmorgan/go-release-manager/vcs"
)

var Command = &cobra.Command{
	Use:   "undo [OPTIONS...]",
	Short: "List / Undo release",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		var releaseUndoActions = []release.UndoAction{}

		undoActions, err := release.ListUndos()
		if err != nil {
			return err
		}
		undoReleases := []string{}
		for k, _ := range undoActions {
			undoReleases = append(undoReleases, k)
		}

		done := false
		for !done {
			release := ""
			action := ""
			if release, err = ui.Select("Undo release", undoReleases); err != nil {
				return err
			}

			if action, err = ui.Select("Action", []string{"View", "Run"}); err != nil {
				return err
			}
			releaseUndoActions = undoActions[release]
			switch action {
			case "View":
				log.Main().SetReportCaller(false)
				names := []string{"All"}
				for i, u := range releaseUndoActions {
					names = append(names, fmt.Sprintf("[%d] %s", i, u.Title))
				}
				if step, err := ui.Select("View Step", names); err != nil {
					return err
				} else {
					for i, u := range releaseUndoActions {
						if fmt.Sprintf("[%d] %s", i, u.Title) == step || step == "All" {
							log.Main().Infof("[%d] undo '%s' -> %s", i, u.Name, u.Title)
							log.Main().Debugf("params:\n\tpath: %s\n\tvcs: %s\n\ttitle: %s\n\tparams: %s", u.Path, u.SourceControl, u.Title, u.Params)
							log.Main().SetReportCaller(false)
						}
					}
				}
			case "Run":
				for i, u := range releaseUndoActions {
					log.Warnf("[%d] running undo '%s' -> %s", i, u.Name)
					if u.VC, err = vcs.Open(u.Path); err != nil {
						return err
					}
					if err = u.Run(); err != nil {
						return err
					}
					dir := config.Get().Workspace.Path.Join(".grlm", "undos").Expand()
					if err = os.Remove(filepath.Join(dir, release)); err != nil && !os.IsNotExist(err) {
						return err
					}
					delete(undoActions, release)
				}
			}
			if len(undoActions) == 0 {
				done = true
			}
		}
		return nil
	},
}
