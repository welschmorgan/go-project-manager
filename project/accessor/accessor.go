package accessor

import "github.com/welschmorgan/go-release-manager/config"

type ProjectAccessor interface {
	// Retrieve the name of this accessor
	Name() string

	// Retrieve the path of the project
	Path() string

	// Initialize a new project
	Initialize(p string, proj *config.Project) error

	// Detect if this accessor is suitable to retrieve data for the given project
	Detect(p string) (bool, error)

	// Open a path, and retrieve all possible data
	Open(p string) error

	// Retrieve the current project version
	CurrentVersion() (string, error)

	// Retrieve the current project name
	CurrentName() (string, error)
}
