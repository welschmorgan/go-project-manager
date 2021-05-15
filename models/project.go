package models

type Project struct {
	Name          string `json:"name,omitempty"`
	Path          string `json:"path,omitempty"`
	SourceControl string `json:"sourceControl,omitempty"`
}

func NewProject(name, path, sourceControl string) *Project {
	return &Project{
		Name:          name,
		Path:          path,
		SourceControl: sourceControl,
	}
}
