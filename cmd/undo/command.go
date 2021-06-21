package release

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/log"
	"github.com/welschmorgan/go-release-manager/release"
	"github.com/welschmorgan/go-release-manager/ui"
	"github.com/welschmorgan/go-release-manager/vcs"
	"gopkg.in/yaml.v2"
)

var Command = &cobra.Command{
	Use:   "undo [OPTIONS...]",
	Short: "List / Undo release",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		var releaseUndoActions = []release.UndoAction{}
		dir := filepath.Join(config.Get().Workspace.Path(), ".grlm-undos")
		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}
		undoActions := map[string][]release.UndoAction{}
		undoFiles := []string{}
		path := ""
		for _, e := range entries {
			path = filepath.Join(dir, e.Name())
			if content, err := os.ReadFile(path); err != nil {
				log.Errorf("Failed to load undo %s, %s", path, err.Error())
			} else {
				if err = yaml.Unmarshal(content, &releaseUndoActions); err != nil {
					log.Errorf("Failed to load undo %s, %s", path, err.Error())
				}
				undoActions[e.Name()] = releaseUndoActions
				undoFiles = append(undoFiles, e.Name())
			}
		}
		sort.Strings(undoFiles)

		done := false
		for !done {
			release := ""
			action := ""
			if release, err = ui.Select("Undo release", undoFiles); err != nil {
				return err
			}

			if action, err = ui.Select("Action", []string{"View", "Run"}); err != nil {
				return err
			}
			releaseUndoActions = undoActions[release]
			switch action {
			case "View":
				for i, u := range releaseUndoActions {
					log.Infof("[%d] loaded undo '%s'", i, u.Name)
				}
			case "Run":
				for i, u := range releaseUndoActions {
					log.Warnf("[%d] running undo '%s'", i, u.Name)
					if u.VC, err = vcs.Open(u.Path); err != nil {
						return err
					}
					if err = u.Run(); err != nil {
						return err
					}
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
