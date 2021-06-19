package config

type Project struct {
	Type          string `json:"type,omitempty"`
	Name          string `json:"name,omitempty"`
	Path          string `json:"path,omitempty"`
	Url           string `json:"url,omitempty"`
	SourceControl string `json:"sourceControl,omitempty"`
}

func NewProject(typ, name, path, url, sourceControl string) *Project {
	return &Project{
		Name:          name,
		Path:          path,
		Url:           url,
		SourceControl: sourceControl,
		Type:          typ,
	}
}
