package init

import (
	"strings"

	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/ui"
)

func askManager(wksp *config.Workspace) error {
	var username string
	var err error
	if wksp.Manager != nil {
		username = wksp.Manager.Name
	}
	if len(strings.TrimSpace(username)) == 0 && wksp.Author != nil {
		username = wksp.Author.Name
	}
	if wksp.Manager, err = ui.AskPerson("Manager", &config.Person{Name: username}); err != nil {
		return err
	}
	return nil
}
