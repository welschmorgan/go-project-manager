package init

import (
	"os"
	"strings"

	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/fs"
	"github.com/welschmorgan/go-release-manager/ui"
)

func askPath(wksp *config.Workspace) error {
	var dir string
	var err error
	var path string = wksp.Path.Expand()
	if len(strings.TrimSpace(path)) == 0 {
		if dir, err = os.Getwd(); err != nil {
			return err
		}
		path = dir
	}
	if path, err = ui.Ask("Path", path, ui.StrMustBeNonEmpty, ui.StrMustNotContainOnlySpaces, ui.PathMustBeDir); err != nil {
		return err
	}
	wksp.Path = fs.Path(path)
	return nil
}
