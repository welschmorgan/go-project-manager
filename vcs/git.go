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
	if remotes, err := g.Remotes(nil); err != nil {
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
func (g *Git) Clone(url, path string, options VersionControlOptions) error {
	fs.Pushd(g.path)
	defer fs.Popd()
	var opts CloneOptions
	if ret, err := getOptions(options, CloneOptions{
		Branch:   "",
		Insecure: true,
	}); err != nil {
		return err
	} else {
		opts = ret.(CloneOptions)
	}
	args := []string{}
	if opts.Insecure {
		args = append(args, "--config", "http.sslVerify=false")
	}
	args = append(args, "clone", url, path)
	if len(strings.TrimSpace(opts.Branch)) > 0 {
		args = append(args, "--branch", opts.Branch)
	}
	_, err, _ := runCommand("git", args...)
	return err
}
func (g *Git) Checkout(branch string, options VersionControlOptions) error {
	fs.Pushd(g.path)
	defer fs.Popd()
	var opts CheckoutOptions
	if ret, err := getOptions(options, CheckoutOptions{
		CreateBranch: false,
	}); err != nil {
		return err
	} else {
		opts = ret.(CheckoutOptions)
	}
	args := []string{
		"checkout",
	}
	if opts.CreateBranch {
		args = append(args, "-b")
	}
	args = append(args, branch)
	_, err, _ := runCommand("git", args...)
	return err
}
func (g *Git) Pull(options VersionControlOptions) error {
	fs.Pushd(g.path)
	defer fs.Popd()
	var opts PullOptions
	if ret, err := getOptions(options, PullOptions{
		Force: false,
		All:   false,
	}); err != nil {
		return err
	} else {
		opts = ret.(PullOptions)
	}
	args := []string{
		"pull",
	}
	if opts.Force {
		args = append(args, "--force")
	}
	if opts.All {
		args = append(args, "--all")
	}
	_, err, _ := runCommand("git", args...)
	return err
}
func (g *Git) Push(options VersionControlOptions) error {
	fs.Pushd(g.path)
	defer fs.Popd()
	var opts PullOptions
	if ret, err := getOptions(options, PullOptions{
		Force: false,
		All:   false,
	}); err != nil {
		return err
	} else {
		opts = ret.(PullOptions)
	}
	args := []string{
		"push",
	}
	if opts.Force {
		args = append(args, "--force")
	}
	if opts.All {
		args = append(args, "--force")
	}
	_, err, _ := runCommand("git", args...)
	return err
}
func (g *Git) Tag(name, commit, message string, options VersionControlOptions) error {
	fs.Pushd(g.path)
	defer fs.Popd()
	_, err, _ := runCommand("git", "--config", "http.sslVerify=false", "tag", "-a", name, "-m", message, commit)
	return err
}
func (g *Git) Merge(source, dest string, options VersionControlOptions) error {
	fs.Pushd(g.path)
	defer fs.Popd()
	var opts MergeOptions
	if ret, err := getOptions(options, MergeOptions{
		FastForwardOnly: false,
		NoFastForward:   true,
	}); err != nil {
		return err
	} else {
		opts = ret.(MergeOptions)
	}
	args := []string{
		"merge",
	}
	if opts.FastForwardOnly {
		args = append(args, "--ff-only")
	}
	if opts.NoFastForward {
		args = append(args, "--no-ff")
	}
	args = append(args, source)
	if err := g.Checkout(dest, nil); err != nil {
		return err
	}
	_, err, _ := runCommand("git", args...)
	return err
}
func (g *Git) Authors(options VersionControlOptions) ([]*models.Person, error) {
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

func (g *Git) Remotes(options VersionControlOptions) (map[string]string, error) {
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
