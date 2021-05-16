package init

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/welschmorgan/go-project-manager/models"
	"github.com/welschmorgan/go-project-manager/vcs"
)

type ProjectMenu struct {
	workspace *models.Workspace
	projects  []*models.Project
	names     []string
	indices   map[string]int
}

func NewProjectMenu(wksp *models.Workspace) (*ProjectMenu, error) {
	menu := &ProjectMenu{
		workspace: wksp,
		projects:  wksp.Projects,
		names:     make([]string, 0),
		indices:   map[string]int{},
	}
	menu.Update()
	if err := menu.Discover(); err != nil {
		return nil, err
	}
	menu.Update()
	return menu, nil
}

func (m *ProjectMenu) Get(name string) *models.Project {
	if id, ok := m.indices[name]; ok {
		return m.projects[id]
	} else {
		return nil
	}
}

func (m *ProjectMenu) ActionEdit(id int, name, path, url, vcs string) error {
	if id < 0 || id >= len(m.projects) {
		return errors.New("invalid project")
	}
	project := m.projects[id]
	if len(strings.TrimSpace(project.Name)) == 0 {
		project.Name = name
	}
	if len(strings.TrimSpace(project.Path)) == 0 {
		project.Path = path
	}
	if len(strings.TrimSpace(project.SourceControl)) == 0 {
		project.SourceControl = vcs
	}
	if len(strings.TrimSpace(project.SourceControl)) == 0 {
		project.Url = url
	}
	return nil
}

func (m *ProjectMenu) Create(name, path, url, vcs string) error {
	m.projects = append(m.projects, models.NewProject(name, path, url, vcs))
	return nil
}

func (m *ProjectMenu) ActionRemove(name string) {
	if id, ok := m.indices[name]; ok {
		m.projects = append(m.projects[:id], m.projects[id+1:]...)
	}
}

func (m *ProjectMenu) ActionClear() {
	m.projects = []*models.Project{}
}

func (m *ProjectMenu) Discover() error {
	var cwd string
	var err error
	var entries []fs.DirEntry
	if cwd, err = os.Getwd(); err != nil {
		return err
	}
	if entries, err = os.ReadDir(cwd); err != nil {
		return err
	}
	for _, dir := range entries {
		if dir.IsDir() && !strings.HasPrefix(strings.TrimSpace(dir.Name()), ".") {
			if sourceControl, err := vcs.Open(filepath.Join(cwd, dir.Name())); err != nil {
				log.Printf("failed to open folder '%s'", err.Error())
			} else {
				if id, ok := m.indices[dir.Name()]; ok {
					if err = m.ActionEdit(id, dir.Name(), sourceControl.Path(), sourceControl.Url(), sourceControl.Name()); err != nil {
						return err
					}
				} else {
					if err = m.Create(dir.Name(), sourceControl.Path(), sourceControl.Url(), sourceControl.Name()); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (m *ProjectMenu) RenderProjects() {
	s := fmt.Sprintf("Found %d projects: ", len(m.projects))
	for id, proj := range m.projects {
		if id > 0 {
			s += ", "
		}
		s += proj.Name
	}
	println(s)
}

func (m *ProjectMenu) Update() {
	if m.projects == nil {
		m.projects = []*models.Project{}
	}
	m.names = []string{}
	m.indices = map[string]int{}
	for _, p := range m.projects {
		m.names = append(m.names, p.Name)
		m.indices[p.Name] = len(m.names) - 1
	}
}
