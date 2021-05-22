package init

import (
	"os"
	"path"
	"strings"

	"github.com/welschmorgan/go-project-manager/models"
	"github.com/welschmorgan/go-project-manager/ui"
)

func askName(wksp *models.Workspace) error {
	var dir string
	var err error
	if len(strings.TrimSpace(wksp.Name)) == 0 {
		if dir, err = os.Getwd(); err != nil {
			return err
		}
		wksp.Name = path.Base(dir)
	}
	if wksp.Name, err = ui.Ask("Name", wksp.Name, strMustNotContainOnlySpaces); err != nil {
		return err
	}
	return nil
}
