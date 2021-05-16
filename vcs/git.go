package vcs

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/welschmorgan/go-project-manager/fs"
	"github.com/welschmorgan/go-project-manager/models"
)

type Git struct {
	VersionControlSoftware
	path string
	url  string
}

func (g *Git) Name() string { return "Git" }
func (g *Git) Path() string { return g.path }
func (g *Git) Url() string  { return g.url }
func (g *Git) Detect(path string) (bool, error) {
	if _, err := os.Stat(filepath.Join(path, ".git")); err != nil {
		return false, err
	}
	return true, nil
}
func (g *Git) Open(p string) error {
	g.path = p
	if remotes, err := g.Remotes(); err != nil {
		return err
	} else if len(remotes) == 0 {
		return fmt.Errorf("no remotes configured for '%s'", filepath.Base(g.path))
	} else {
		g.url = ""
		for _, r := range remotes {
			g.url = r
			break
		}
	}
	return nil
}
func (g *Git) Clone(url, path string) error           { return errNotYetImpl }
func (g *Git) Checkout(branch string) error           { return errNotYetImpl }
func (g *Git) Pull() error                            { return errNotYetImpl }
func (g *Git) Push() error                            { return errNotYetImpl }
func (g *Git) Tag(name, commit, message string) error { return errNotYetImpl }
func (g *Git) Merge(source, dest string) error        { return errNotYetImpl }
func (g *Git) Authors() ([]*models.Person, error) {
	fs.Pushd(g.path)
	defer fs.Popd()
	var lines []string
	var err error
	if lines, err, _ = runCommand("git", "log", "--format=%cn <%ce>"); err != nil {
		return nil, err
	}
	ret := []*models.Person{}
	for _, line := range lines {
		rule := regexp.MustCompile("(.*)<(.*?)>")
		matches := rule.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			ret = append(ret, models.NewPerson(strings.TrimSpace(match[1]), strings.TrimSpace(match[2]), ""))
		}
	}
	return ret, nil
}

func (g *Git) Remotes() (map[string]string, error) {
	fs.Pushd(g.path)
	defer fs.Popd()
	var lines []string
	var err error
	if lines, err, _ = runCommand("git", "remote", "-v"); err != nil {
		return nil, err
	}
	ret := map[string]string{}
	for _, line := range lines {
		rule := regexp.MustCompile(`(\w+)\s+(.*)\s+\((\w+)\)`)
		matches := rule.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			ret[strings.TrimSpace(match[1])] = strings.TrimSpace(match[2])
		}
	}
	return ret, nil
}
