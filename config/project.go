package config

type Project struct {
	Type          string `yaml:"type,omitempty"`
	Name          string `yaml:"name,omitempty"`
	Path          string `yaml:"path,omitempty"`
	Url           string `yaml:"url,omitempty"`
	SourceControl string `yaml:"source_control,omitempty"`
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
