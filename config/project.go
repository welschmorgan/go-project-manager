package config

import "github.com/welschmorgan/go-release-manager/fs"

type Project struct {
	Type          string  `yaml:"type,omitempty"`
	Name          string  `yaml:"name,omitempty"`
	Path          fs.Path `yaml:"path,omitempty"`
	Url           string  `yaml:"url,omitempty"`
	SourceControl string  `yaml:"source_control,omitempty"`
}

func NewProject(typ, name, path, url, sourceControl string) *Project {
	return &Project{
		Name:          name,
		Path:          fs.Path(path),
		Url:           url,
		SourceControl: sourceControl,
		Type:          typ,
	}
}
