package ui

import (
	"github.com/welschmorgan/go-project-manager/models"
)

func AskPerson(label string, defaults *models.Person, validators ...ObjValidator) (*models.Person, error) {
	if pers, err := AskObject(label, defaults, validators...); err != nil {
		return nil, err
	} else {
		ret := pers.(models.Person)
		return &ret, nil
	}
}
