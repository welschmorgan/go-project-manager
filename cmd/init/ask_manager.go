package init

import (
	"strings"

	"github.com/welschmorgan/go-project-manager/models"
	"github.com/welschmorgan/go-project-manager/ui"
)

func askManager(wksp *models.Workspace) error {
	var username string
	var err error
	if wksp.Manager != nil {
		username = wksp.Manager.Name
	}
	if len(strings.TrimSpace(username)) == 0 && wksp.Author != nil {
		username = wksp.Author.Name
	}
	if wksp.Manager, err = ui.AskPerson("Manager", &models.Person{Name: username}); err != nil {
		return err
	}
	return nil
}
