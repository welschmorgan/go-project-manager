package project

import (
	"github.com/welschmorgan/go-release-manager/config"
)

type ProjectAccessor interface {
	// Retrieve the basic project infos: url, name, path, sourceControl
	Infos() *config.Project

	// Detect if this accessor is suitable to retrieve data for the given project
	Detect() (bool, error)

	// Open a path, and retrieve all possible data
	Open(p string) error

	// Retrieve the project version
	CurrentVersion() (string, error)
}
