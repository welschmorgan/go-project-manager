package accessor

import (
	"reflect"
	"strings"
)

type ProjectAccessor interface {
	// Retrieve the name of this accessor
	AccessorName() string

	// Retrieve the path of the project
	Path() string

	// Retrieve the name of the project file
	DescriptionFile() string

	// Initialize a new project
	Scaffold(ctx *FinalizationContext) error

	// Retrieve project scaffolder
	Scaffolder() Scaffolder

	// Detect if this accessor is suitable to retrieve data for the given project
	Detect(p string) (bool, error)

	// Open a path, and retrieve all possible data
	Open(p string) error

	// Retrieve the current project version
	Version() (string, error)

	// Retrieve the current project name
	Name() (string, error)
}

var accessors []ProjectAccessor = []ProjectAccessor{}

func instanciate(a ProjectAccessor) ProjectAccessor {
	inst := reflect.New(reflect.TypeOf(a).Elem())
	return inst.Interface().(ProjectAccessor)
}

func Register(a ProjectAccessor) {
	accessors = append(accessors, a)
}

func Get(n string) ProjectAccessor {
	loName := strings.ToLower(n)
	for _, a := range accessors {
		if strings.ToLower(a.AccessorName()) == loName {
			return instanciate(a)
		}
	}
	return nil
}

func GetAll() []ProjectAccessor {
	ret := []ProjectAccessor{}
	for _, a := range accessors {
		ret = append(ret, instanciate(a))
	}
	return ret
}

func GetAllNames() []string {
	ret := []string{}
	for _, a := range accessors {
		ret = append(ret, a.AccessorName())
	}
	return ret
}
