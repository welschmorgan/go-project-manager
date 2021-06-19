package vcs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/fs"
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
	if fi, err := os.Stat(filepath.Join(path, ".git")); err != nil {
		return false, err
	} else if !fi.IsDir() {
		return false, fmt.Errorf("%s: not a directory", path)
	}
	return true, nil
}

func (g *Git) ListBranches(options VersionControlOptions) ([]string, error) {
	fs.Pushd(g.path)
	defer fs.Popd()
	var opts BranchOptions
	if ret, err := getOptions(options, BranchOptions{}); err != nil {
		return nil, err
	} else {
		opts = ret.(BranchOptions)
	}
	args := []string{}
	args = append(args, "branch")
	if opts.All {
		args = append(args, "--all")
	}
	if opts.Verbose {
		args = append(args, "--verbose")
	}
	if len(opts.SetUpstreamTo) > 0 {
		args = append(args, "--set-upstream-to", opts.SetUpstreamTo)
	}
	code, out, errTxt, err := runCommand("git", args...)
	dumpCommandErrors(code, errTxt)
	if err != nil {
		return nil, err
	}
	for i, _ := range out {
		if strings.HasPrefix(out[i], "*") {
			out[i] = strings.TrimSpace(strings.Replace(out[i], "*", "", 1))
		}
	}
	return out, nil
}

func (g *Git) Open(p string) error {
	g.path = p
	if remotes, err := g.ListRemotes(nil); err != nil {
		return err
	} else {
		if len(remotes) == 0 {
			if config.Get().DryRun {
				remotes = map[string]string{
					"fake-for-dry-run": "http://fake.com",
				}
			} else {
				fmt.Fprintf(os.Stderr, "%s[\033[1;31m-\033[0m] no remotes configured for '%s'\n", strings.Repeat("\t", config.Get().Indent), filepath.Base(g.path))
			}
		}
		g.url = ""
		for _, r := range remotes {
			g.url = r
			break
		}
	}
	return nil
}

func (g *Git) Status(options VersionControlOptions) ([]string, error) {
	fs.Pushd(g.path)
	defer fs.Popd()
	var opts StatusOptions
	if ret, err := getOptions(options, StatusOptions{
		Short: true,
	}); err != nil {
		return nil, err
	} else {
		opts = ret.(StatusOptions)
	}
	args := []string{}
	args = append(args, "status")
	if opts.Short {
		args = append(args, "--short")
	}
	code, out, errTxt, err := runCommand("git", args...)
	dumpCommandErrors(code, errTxt)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (g *Git) Stash(options VersionControlOptions) ([]string, error) {
	fs.Pushd(g.path)
	defer fs.Popd()
	var opts StashOptions
	if ret, err := getOptions(options, StashOptions{
		Save:             true,
		IncludeUntracked: true,
	}); err != nil {
		return nil, err
	} else {
		opts = ret.(StashOptions)
	}
	args := []string{"stash"}
	if opts.Save {
		args = append(args, "save")
	} else if opts.List {
		args = append(args, "list")
	} else if opts.Apply {
		args = append(args, "apply")
	} else if opts.Pop {
		args = append(args, "pop")
	}
	if opts.IncludeUntracked {
		args = append(args, "-u")
	}
	if len(strings.TrimSpace(opts.Message)) > 0 {
		args = append(args, opts.Message)
	}
	code, out, errTxt, err := runCommand("git", args...)
	dumpCommandErrors(code, errTxt)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (g *Git) DeleteBranch(name string, options VersionControlOptions) error {
	fs.Pushd(g.path)
	defer fs.Popd()
	var opts DeleteBranchOptions
	if ret, err := getOptions(options, DeleteBranchOptions{
		Local:  true,
		Remote: false,
	}); err != nil {
		return err
	} else {
		opts = ret.(DeleteBranchOptions)
	}
	if opts.Local {
		code, _, errTxt, err := runCommand("git", []string{"branch", "-D", name}...)
		dumpCommandErrors(code, errTxt)
		if err != nil {
			return err
		}
	}
	if opts.Remote {
		remoteName := opts.RemoteName
		if len(remoteName) == 0 {
			remoteName = "origin"
		}
		code, _, errTxt, err := runCommand("git", []string{"push", remoteName, ":" + name}...)
		dumpCommandErrors(code, errTxt)
		if err != nil {
			return err
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
	code, _, errTxt, err := runCommand("git", args...)
	dumpCommandErrors(code, errTxt)
	return err
}
func (g *Git) Checkout(branch string, options VersionControlOptions) error {
	fs.Pushd(g.path)
	defer fs.Popd()
	var opts CheckoutOptions
	if ret, err := getOptions(options, CheckoutOptions{
		CreateBranch:     false,
		UpdateIfExisting: false,
		StartingPoint:    "",
	}); err != nil {
		return err
	} else {
		opts = ret.(CheckoutOptions)
	}
	args := []string{
		"checkout",
	}
	if opts.CreateBranch {
		if opts.UpdateIfExisting {
			args = append(args, "-B")
		} else {
			args = append(args, "-b")
		}
	}
	args = append(args, branch)
	if len(strings.TrimSpace(opts.StartingPoint)) > 0 {
		args = append(args, strings.TrimSpace(opts.StartingPoint))
	}
	code, _, errTxt, err := runCommand("git", args...)
	dumpCommandErrors(code, errTxt)
	return err
}

func (g *Git) Reset(options VersionControlOptions) error {
	fs.Pushd(g.path)
	defer fs.Popd()
	var opts ResetOptions
	if ret, err := getOptions(options, ResetOptions{
		Hard:   false,
		Commit: "",
	}); err != nil {
		return err
	} else {
		opts = ret.(ResetOptions)
	}
	args := []string{
		"reset",
	}
	if opts.Hard {
		args = append(args, "--hard")
	}
	if len(strings.TrimSpace(opts.Commit)) > 0 {
		args = append(args, strings.TrimSpace(opts.Commit))
	}
	code, _, errTxt, err := runCommand("git", args...)
	dumpCommandErrors(code, errTxt)
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
	if opts.ListTags {
		args = append(args, "--tags")
	}
	code, _, errTxt, err := runCommand("git", args...)
	dumpCommandErrors(code, errTxt)
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
	code, _, errTxt, err := runCommand("git", args...)
	dumpCommandErrors(code, errTxt)
	return err
}

func (g *Git) Tag(name string, options VersionControlOptions) error {
	fs.Pushd(g.path)
	defer fs.Popd()
	var opts TagOptions
	if ret, err := getOptions(options, TagOptions{
		Delete:    false,
		Annotated: false,
		Message:   "",
		Commit:    "",
	}); err != nil {
		return err
	} else {
		opts = ret.(TagOptions)
	}
	args := []string{
		"tag",
	}
	if opts.Annotated {
		args = append(args, "-a")
	} else if opts.Delete {
		args = append(args, "-d")
	}
	args = append(args, name)
	if len(strings.TrimSpace(opts.Message)) > 0 {
		args = append(args, "-m", opts.Message)
	}
	if len(strings.TrimSpace(opts.Commit)) > 0 {
		args = append(args, opts.Commit)
	}
	code, _, errTxt, err := runCommand("git", args...)
	dumpCommandErrors(code, errTxt)
	return err
}

func (g *Git) CurrentBranch() (string, error) {
	fs.Pushd(g.path)
	defer fs.Popd()
	args := []string{
		"rev-parse", "--abbrev-ref", "HEAD",
	}
	code, out, errTxt, err := runCommand("git", args...)
	dumpCommandErrors(code, errTxt)
	if err != nil {
		return "", err
	}
	if len(out) == 0 {
		if config.Get().DryRun {
			return "my-branch", nil
		} else {
			return "", errors.New("no branch name found")
		}
	}
	return out[0], nil
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
	if opts.FastForwardOnly && opts.NoFastForward {
		return errors.New("--ff-only and --no-ff are mutually exclusive, please pick one")
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
	code, _, errTxt, err := runCommand("git", args...)
	dumpCommandErrors(code, errTxt)
	return err
}
func (g *Git) ListAuthors(options VersionControlOptions) ([]*config.Person, error) {
	fs.Pushd(g.path)
	defer fs.Popd()
	var lines []string
	var err error
	var errTxt []string
	var code int
	if code, lines, errTxt, err = runCommand("git", "log", "--format=%cn <%ce>"); err != nil {
		return nil, err
	}
	dumpCommandErrors(code, errTxt)
	ret := []*config.Person{}
	for _, line := range lines {
		rule := regexp.MustCompile("(.*)<(.*?)>")
		matches := rule.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			ret = append(ret, config.NewPerson(strings.TrimSpace(match[1]), strings.TrimSpace(match[2]), ""))
		}
	}
	return ret, nil
}

func (g *Git) ListRemotes(options VersionControlOptions) (map[string]string, error) {
	fs.Pushd(g.path)
	defer fs.Popd()
	var lines []string
	var err error
	var errTxt []string
	var code int
	if code, lines, errTxt, err = runCommand("git", "remote", "-v"); err != nil {
		return nil, err
	}
	dumpCommandErrors(code, errTxt)
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

func (g *Git) ListTags(options VersionControlOptions) ([]string, error) {
	fs.Pushd(g.path)
	defer fs.Popd()
	var err error
	var opts ListTagsOptions
	if ret, err := getOptions(options, ListTagsOptions{
		SortByTaggerDate:    true,
		SortByCommitterDate: false,
	}); err != nil {
		return nil, err
	} else {
		opts = ret.(ListTagsOptions)
	}
	args := []string{
		"tag", "-l",
	}
	if opts.SortByCommitterDate {
		args = append(args, "--sort=committerdate")
	}
	if opts.SortByTaggerDate {
		args = append(args, "--sort=taggerdate")
	}
	code, out, errTxt, err := runCommand("git", args...)
	dumpCommandErrors(code, errTxt)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (g *Git) Initialize(path string, options VersionControlOptions) error {
	fs.Pushd(g.path)
	defer fs.Popd()
	var err error
	var opts InitOptions
	if ret, err := getOptions(options, InitOptions{
		Bare: false,
	}); err != nil {
		return err
	} else {
		opts = ret.(InitOptions)
	}
	args := []string{
		"init",
	}
	if opts.Bare {
		args = append(args, "--bare")
	}
	args = append(args, path)
	code, _, errTxt, err := runCommand("git", args...)
	dumpCommandErrors(code, errTxt)
	return err
}
