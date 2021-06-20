package ui

import (
	"github.com/welschmorgan/go-release-manager/config"
)

func AskPerson(label string, defaults *config.Person, validators ...ObjValidator) (*config.Person, error) {
	defaultName := ""
	defaultEmail := ""
	defaultPhone := ""
	if defaults != nil {
		defaultName = defaults.Name
		defaultEmail = defaults.Email
		defaultPhone = defaults.Phone
	}
	if pers, err := AskObject(label, defaults, map[string]ItemFieldType{
		"Name":  NewItemFieldType(ItemFieldText, defaultName),
		"Email": NewItemFieldType(ItemFieldText, defaultEmail),
		"Phone": NewItemFieldType(ItemFieldText, defaultPhone),
	}, validators...); err != nil {
		return nil, err
	} else {
		ret := pers.(config.Person)
		return &ret, nil
	}
}
