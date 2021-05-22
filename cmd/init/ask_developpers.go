package init

import "github.com/welschmorgan/go-project-manager/config"

func askDeveloppers(wksp *config.Workspace) error {
	var err error
	var menu *DevelopperMenu
	if menu, err = NewDevelopperMenu(wksp); err != nil {
		panic(err)
	} else {
		if err := menu.Render(); err != nil {
			panic(err)
		}
	}
	// wksp.Developpers = developpers
	return nil
}
