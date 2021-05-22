package ui

import (
	"github.com/welschmorgan/go-project-manager/config"
)

func AskPerson(label string, defaults *config.Person, validators ...ObjValidator) (*config.Person, error) {
	if pers, err := AskObject(label, defaults, validators...); err != nil {
		return nil, err
	} else {
		ret := pers.(config.Person)
		return &ret, nil
	}
}
