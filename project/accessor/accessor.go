package accessor

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/welschmorgan/go-release-manager/version"
)

type VersionFileManipulator struct {
	// Detect if a specific file contains matching version string
	Detect func(file, content string, currentVersion *version.Version) bool

	// Update a specific file, replacing currentVersion with nextVersion
	// It returns the modified content
	Update func(file, content string, currentVersion, nextVersion *version.Version) (string, error)
}

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
	ReadVersion() (version.Version, error)

	// Define the current project version
	WriteVersion(v *version.Version) error

	// Retrieve the possible list of files that need version updates
	VersionManipulators() map[string]VersionFileManipulator

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

func Detect(p string) (found ProjectAccessor, err error) {
	errs := []string{}
	for _, a := range accessors {
		if _, err = a.Detect(p); err != nil {
			errs = append(errs, err.Error())
		} else {
			found = a
			break
		}
	}
	err = nil
	if found == nil {
		extra := ""
		if len(errs) > 0 {
			extra += fmt.Sprintf(":\n\t- %s", strings.Join(errs, "\n\t- "))
		}
		err = fmt.Errorf("failed to find suitable project accessor for '%s'%s", p, extra)
	}
	return found, err
}

func Open(p string) (found ProjectAccessor, err error) {
	if found, err = Detect(p); err != nil {
		return nil, err
	}
	err = found.Open(p)
	return found, err
}
