package init

import (
	"os"
	"strings"

	"github.com/welschmorgan/go-project-manager/config"
	"github.com/welschmorgan/go-project-manager/ui"
)

func askPath(wksp *config.Workspace) error {
	var dir string
	var err error
	if len(strings.TrimSpace(wksp.Path)) == 0 {
		if dir, err = os.Getwd(); err != nil {
			return err
		}
		wksp.Path = dir
	}
	if wksp.Path, err = ui.Ask("Path", wksp.Path, ui.StrMustBeNonEmpty, ui.StrMustNotContainOnlySpaces, ui.PathMustBeDir); err != nil {
		return err
	}
	return nil
}
