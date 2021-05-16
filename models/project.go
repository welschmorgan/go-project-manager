package models

type Project struct {
	Name          string `json:"name,omitempty"`
	Path          string `json:"path,omitempty"`
	Url           string `json:"url,omitempty"`
	SourceControl string `json:"sourceControl,omitempty"`
}

func NewProject(name, path, url, sourceControl string) *Project {
	return &Project{
		Name:          name,
		Path:          path,
		Url:           url,
		SourceControl: sourceControl,
	}
}
