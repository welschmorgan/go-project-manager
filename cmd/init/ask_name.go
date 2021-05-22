package init

import (
	"os"
	"path"
	"strings"

	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/ui"
)

func askName(wksp *config.Workspace) error {
	var dir string
	var err error
	if len(strings.TrimSpace(wksp.Name)) == 0 {
		if dir, err = os.Getwd(); err != nil {
			return err
		}
		wksp.Name = path.Base(dir)
	}
	if wksp.Name, err = ui.Ask("Name", wksp.Name, ui.StrMustNotContainOnlySpaces); err != nil {
		return err
	}
	return nil
}
