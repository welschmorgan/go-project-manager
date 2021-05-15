package workspace

type Workspace struct {
	Name     string   `json:"name"`
	Path     string   `json:"path"`
	Projects []string `json:"projects"`
}
