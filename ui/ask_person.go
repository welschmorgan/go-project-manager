package ui

import "github.com/welschmorgan/go-project-manager/models"

func AskPerson(label string, defaults *models.Person, validators ...ObjValidator) (*models.Person, error) {
	defaultName, defaultEmail, defaultPhone := "", "", ""
	if defaults != nil {
		defaultName = defaults.Name
		defaultEmail = defaults.Email
		defaultPhone = defaults.Phone
	}
	validator := NewMultiObjValidator(validators...)
	var ret *models.Person = nil
	if name, err := Ask(label+".name", defaultName, validator("name")...); err != nil {
		return nil, err
	} else if email, err := Ask(label+".email", defaultEmail, validator("email")...); err != nil {
		return nil, err
	} else if phone, err := Ask(label+".phone", defaultPhone, validator("phone")...); err != nil {
		return nil, err
	} else {
		ret = &models.Person{
			Name:  name,
			Email: email,
			Phone: phone,
		}
	}
	return ret, nil
}
