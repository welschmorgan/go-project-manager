package ui

import (
	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/vcs"
)

func AskProject(label string, defaults *config.Project, validators ...ObjValidator) (*config.Project, error) {
	defaultName := ""
	defaultPath := ""
	defaultUrl := ""
	if defaults != nil {
		defaultName = defaults.Name
		defaultPath = defaults.Path.Raw()
		defaultUrl = defaults.Url
	}
	if proj, err := AskObject(label, defaults, map[string]ItemFieldType{
		"Name":          NewItemFieldType(ItemFieldText, defaultName),
		"Path":          NewItemFieldType(ItemFieldText, defaultPath),
		"Url":           NewItemFieldType(ItemFieldText, defaultUrl),
		"SourceControl": NewItemFieldType(ItemFieldList, vcs.AllNames),
	}, validators...); err != nil {
		return nil, err
	} else {
		ret := proj.(config.Project)
		return &ret, nil
	}
}
