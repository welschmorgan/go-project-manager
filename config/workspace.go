package config

type Person struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

type Workspace struct {
	Name          string   `json:"name"`
	Path          string   `json:"path"`
	Projects      []string `json:"projects"`
	SourceControl string   `json:"sourceControl"`
	Author        Person   `json:"author"`
	Manager       Person   `json:"manager"`
	Developpers   []Person `json:"developpers"`
}
