package models

type Workspace struct {
	Name        string     `json:"name"`
	Path        string     `json:"path"`
	Projects    []*Project `json:"projects"`
	Author      *Person    `json:"author"`
	Manager     *Person    `json:"manager"`
	Developpers []*Person  `json:"developpers"`
}

func NewWorkspace(name, path string, projects []*Project, sourceControl string, author *Person, manager *Person, developpers []*Person) *Workspace {
	return &Workspace{
		Name:        name,
		Path:        path,
		Projects:    projects,
		Author:      author,
		Manager:     manager,
		Developpers: developpers,
	}
}
