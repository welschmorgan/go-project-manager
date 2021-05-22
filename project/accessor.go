package project

import (
	"fmt"
	"reflect"
	"strings"
)

type ProjectAccessor interface {
	// Retrieve the name of this accessor
	Name() string

	// Retrieve the path of the project
	Path() string

	// Detect if this accessor is suitable to retrieve data for the given project
	Detect(p string) (bool, error)

	// Open a path, and retrieve all possible data
	Open(p string) error

	// Retrieve the project version
	CurrentVersion() (string, error)
}

var All []ProjectAccessor = []ProjectAccessor{
	&NodeProjectAccessor{},
}

func instanciate(a ProjectAccessor) ProjectAccessor {
	inst := reflect.New(reflect.TypeOf(a).Elem())
	return inst.Interface().(ProjectAccessor)
}

func Get(n string) ProjectAccessor {
	loName := strings.ToLower(n)
	for _, a := range All {
		if strings.ToLower(a.Name()) == loName {
			return instanciate(a)
		}
	}
	return nil
}

func Detect(p string) (ProjectAccessor, error) {
	for _, a := range All {
		if ok, err := a.Detect(p); ok {
			if err != nil {
				return nil, err
			}
			return instanciate(a), nil
		}
	}
	return nil, fmt.Errorf("no accessor found for '%s'", p)
}

func Open(p string) (ProjectAccessor, error) {
	if a, err := Detect(p); err != nil {
		return nil, err
	} else {
		ret := instanciate(a)
		if err = ret.Open(p); err != nil {
			return nil, err
		}
		return ret, nil
	}
}
