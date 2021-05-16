package ui

import (
	"github.com/welschmorgan/go-project-manager/models"
)

func AskProject(label string, defaults *models.Project, validators ...ObjValidator) (*models.Project, error) {
	if proj, err := AskObject(label, defaults, validators...); err != nil {
		return nil, err
	} else {
		ret := proj.(models.Project)
		return &ret, nil
	}
}
