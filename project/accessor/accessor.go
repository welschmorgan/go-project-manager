package accessor

import "github.com/welschmorgan/go-release-manager/config"

type ProjectAccessor interface {
	// Retrieve the name of this accessor
	AccessorName() string

	// Retrieve the path of the project
	Path() string

	// Retrieve the name of the project file
	DescriptionFile() string

	// Initialize a new project
	Initialize(p string, proj *config.Project) error

	// Detect if this accessor is suitable to retrieve data for the given project
	Detect(p string) (bool, error)

	// Open a path, and retrieve all possible data
	Open(p string) error

	// Retrieve the current project version
	Version() (string, error)

	// Retrieve the current project name
	Name() (string, error)
}
