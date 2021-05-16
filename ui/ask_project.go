package ui

import (
	"strings"

	"github.com/welschmorgan/go-project-manager/models"
)

func AskProject(label string, defaults *models.Project, validators ...ObjValidator) (*models.Project, error) {
	defaultName, defaultPath, defaultUrl, defaultSourceControl := "", "", "", ""
	if defaults != nil {
		defaultName = defaults.Name
		defaultPath = defaults.Path
		defaultUrl = defaults.Url
		defaultSourceControl = defaults.SourceControl
	}
	validator := NewMultiObjValidator(validators...)
	var err error
	ret := &models.Project{}
	if ret.Name, err = Ask(label+".name", defaultName, validator("name")...); err != nil {
		return nil, err
	} else {
		if len(strings.TrimSpace(ret.Name)) > 0 {
			if ret.Path, err = Ask(label+".path", defaultPath, validator("path")...); err != nil {
				return nil, err
			}
			if ret.Url, err = Ask(label+".url", defaultUrl, validator("url")...); err != nil {
				return nil, err
			}
			if ret.SourceControl, err = Ask(label+".sourceControl", defaultSourceControl, validator("sourceControl")...); err != nil {
				return nil, err
			}
		}

	}
	return ret, nil
}
