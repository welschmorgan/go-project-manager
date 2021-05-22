package init

import (
	"os/user"
	"strings"

	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/ui"
)

func askAuthor(wksp *config.Workspace) error {
	var currentUser *user.User
	var err error
	var defaultAuthor *config.Person = wksp.Author
	if defaultAuthor != nil && len(strings.TrimSpace(defaultAuthor.Name)) == 0 {
		if currentUser, err = user.Current(); err != nil {
			return err
		}
		defaultAuthor.Name = currentUser.Name
		if len(strings.TrimSpace(defaultAuthor.Name)) == 0 {
			defaultAuthor.Name = currentUser.Username
		}
	}
	if wksp.Author, err = ui.AskPerson("Author", defaultAuthor); err != nil {
		return err
	}
	return nil
}
