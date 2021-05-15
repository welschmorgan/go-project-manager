package ui

import (
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/welschmorgan/go-project-manager/config"
)

type Validator func(string) error
type ObjValidator func(k, v string) error

func Ask(label, defaultValue string, validators ...Validator) (string, error) {
	// templates := &promptui.PromptTemplates{
	// 	Prompt:  "{{ . }} ",
	// 	Valid:   "{{ . | green }} ",
	// 	Invalid: "{{ . | red }} ",
	// 	Success: "{{ . | bold }} ",
	// }

	prompt := promptui.Prompt{
		Label: label,
		// Templates: templates,
		AllowEdit: true,
		Default:   defaultValue,
	}
	prompt.Validate = func(s string) error {
		var err error
		for _, v := range validators {
			if v != nil {
				if err = v(s); err != nil {
					return err
				}
			}
		}
		return nil
	}
	return prompt.Run()
}

func AskPerson(label string, defaults *config.Person, validators ...ObjValidator) (*config.Person, error) {
	defaultName, defaultEmail, defaultPhone := "", "", ""
	if defaults != nil {
		defaultName = defaults.Name
		defaultEmail = defaults.Email
		defaultPhone = defaults.Phone
	}
	validator := func(k string) []Validator {
		ret := []Validator{}
		for _, validator := range validators {
			if validator != nil {
				ret = append(ret, func(v string) error {
					return validator(k, v)
				})
			}
		}
		return ret
	}
	var ret *config.Person = nil
	if name, err := Ask(label+".name", defaultName, validator("name")...); err != nil {
		return nil, err
	} else if email, err := Ask(label+".email", defaultEmail, validator("email")...); err != nil {
		return nil, err
	} else if phone, err := Ask(label+".phone", defaultPhone, validator("phone")...); err != nil {
		return nil, err
	} else {
		ret = &config.Person{
			Name:  name,
			Email: email,
			Phone: phone,
		}
	}
	return ret, nil
}

func AskProject(label string, defaults *config.Project, validators ...ObjValidator) (*config.Project, error) {
	defaultName, defaultPath := "", ""
	if defaults != nil {
		defaultName = defaults.Name
		defaultPath = defaults.Path
	}
	validator := func(k string) []Validator {
		ret := []Validator{}
		for _, validator := range validators {
			if validator != nil {
				ret = append(ret, func(v string) error {
					return validator(k, v)
				})
			}
		}
		return ret
	}
	var ret *config.Project = nil
	if name, err := Ask(label+".name", defaultName, validator("name")...); err != nil {
		return nil, err
	} else {
		var path string
		var err error
		if len(strings.TrimSpace(name)) > 0 {
			if path, err = Ask(label+".path", defaultPath, validator("path")...); err != nil {
				return nil, err
			}
		}
		ret = &config.Project{
			Name: name,
			Path: path,
		}
	}
	return ret, nil
}

func Select(label string, items []string, validator func(string) error) (string, error) {

	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	_, result, err := prompt.Run()
	return result, err
}
