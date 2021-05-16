package init

import (
	"fmt"
	"strings"

	"github.com/welschmorgan/go-project-manager/models"
	"github.com/welschmorgan/go-project-manager/ui"
	"github.com/welschmorgan/go-project-manager/vcs"
)

func askDeveloppers(wksp *models.Workspace) error {
	var developpers []*models.Person = wksp.Developpers
	if developpers == nil {
		developpers = []*models.Person{}
	}
	developperNames := []string{}
	developperIds := map[string]int{}
	getDevelopperId := func(name string) int {
		for id, d := range developpers {
			if d != nil && d.Name == name {
				return id
			}
		}
		return -1
	}
	for _, project := range wksp.Projects {
		s := vcs.Get(project.SourceControl)
		if err := s.Open(project.Path); err != nil {
			return err
		}
		if projectDeveloppers, err := s.Authors(nil); err != nil {
			return err
		} else {
			for _, tmpDev := range projectDeveloppers {
				id := getDevelopperId(tmpDev.Name)
				if id != -1 {
					dev := developpers[id]
					if len(strings.TrimSpace(dev.Name)) == 0 {
						dev.Name = tmpDev.Name
					}
					if len(strings.TrimSpace(dev.Phone)) == 0 {
						dev.Phone = tmpDev.Phone
					}
					if len(strings.TrimSpace(dev.Email)) == 0 {
						dev.Email = tmpDev.Email
					}
				} else {
					developpers = append(developpers, tmpDev)
				}
			}
		}
	}
	printDeveloppers := func() {
		developperNames = []string{}
		developperIds = map[string]int{}
		s := fmt.Sprintf("Found %d developpers: ", len(developpers))
		for id, a := range developpers {
			if id > 0 {
				s += ", "
			}
			s += a.Name
			developperNames = append(developperNames, a.Name)
			developperIds[a.Name] = id
		}
		println(s)
	}
	done := false
	var action string
	var developperName string
	var err error
	for !done {
		printDeveloppers()
		if action, err = ui.Select("Action", []string{"ActionQuit", "ActionAdd", "ActionRemove", "ActionEdit", "ActionClear"}, nil); err != nil {
			return err
		}
		if action == "ActionEdit" || action == "ActionRemove" {
			if developperName, err = ui.Select("Developper", developperNames, nil); err != nil {
				return err
			}
		}
		defaultDevelopper := models.Person{}
		if action == "ActionEdit" {
			developper := developpers[developperIds[developperName]]
			defaultDevelopper.Name = developper.Name
			defaultDevelopper.Email = developper.Email
			defaultDevelopper.Phone = developper.Phone
		}
		if action == "ActionEdit" || action == "ActionAdd" {
			if auth, err := ui.AskPerson("Developper", &defaultDevelopper, nil); err != nil {
				return err
			} else {
				oldId := developperIds[auth.Name]
				developpers = append(developpers, auth)
				if action == "ActionEdit" {
					developpers = append(developpers[:oldId], developpers[oldId+1:]...)
				}
			}
		}
		if action == "ActionRemove" {
			id := developperIds[developperName]
			developpers = append(developpers[:id], developpers[id+1:]...)
		}
		if action == "ActionClear" {
			developpers = []*models.Person{}
		}
		if action == "ActionQuit" {
			done = true
		}
	}
	wksp.Developpers = developpers
	return nil
}
