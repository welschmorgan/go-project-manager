package ui

import (
	"github.com/welschmorgan/go-project-manager/config"
)

func AskProject(label string, defaults *config.Project, validators ...ObjValidator) (*config.Project, error) {
	if proj, err := AskObject(label, defaults, validators...); err != nil {
		return nil, err
	} else {
		ret := proj.(config.Project)
		return &ret, nil
	}
}
