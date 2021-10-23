package config

import "github.com/welschmorgan/go-release-manager/fs"

type Project struct {
	Type          string  `json:"type,omitempty" yaml:"type,omitempty"`
	Name          string  `json:"name,omitempty" yaml:"name,omitempty"`
	Path          fs.Path `json:"path,omitempty" yaml:"path,omitempty"`
	Url           string  `json:"url,omitempty" yaml:"url,omitempty"`
	SourceControl string  `json:"source_control,omitempty" yaml:"source_control,omitempty"`
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
