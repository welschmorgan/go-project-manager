package vcs

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/welschmorgan/go-project-manager/models"
)

type Git struct {
	VersionControlSoftware
	path string
}

func (g *Git) Name() string { return "Git" }
func (g *Git) Path() string { return g.path }
func (g *Git) Detect(path string) (bool, error) {
	if _, err := os.Stat(filepath.Join(path, ".git")); err != nil {
		return false, err
	}
	return true, nil
}
func (g *Git) Open(p string) error {
	g.path = p
	return nil
}
func (g *Git) Clone(url, path string) error           { return errNotYetImpl }
func (g *Git) Checkout(branch string) error           { return errNotYetImpl }
func (g *Git) Pull() error                            { return errNotYetImpl }
func (g *Git) Push() error                            { return errNotYetImpl }
func (g *Git) Tag(name, commit, message string) error { return errNotYetImpl }
func (g *Git) Merge(source, dest string) error        { return errNotYetImpl }
func (g *Git) Authors() ([]*models.Person, error) {
	Pushd(g.path)
	defer Popd()
	var stdout bytes.Buffer
	cmd := exec.Command("git", "log", "--format=%cn <%ce>")
	cmd.Stderr = os.Stderr
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	ret := []*models.Person{}
	lines := map[string]bool{}
	for line, err := stdout.ReadString('\n'); err == nil; line, err = stdout.ReadString('\n') {
		line = strings.TrimSpace(line)
		if ok := lines[line]; !ok {
			lines[line] = true
		}
	}
	for line, _ := range lines {
		rule := regexp.MustCompile("(.*)<(.*?)>")
		matches := rule.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			ret = append(ret, models.NewPerson(strings.TrimSpace(match[1]), strings.TrimSpace(match[2]), ""))
		}
	}
	return ret, nil
}
