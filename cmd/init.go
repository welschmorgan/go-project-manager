package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/welschmorgan/go-project-manager/config"
	"github.com/welschmorgan/go-project-manager/ui"
	"github.com/welschmorgan/go-project-manager/vcs"
)

func strMustBeNonEmpty(s string) error {
	if len(s) == 0 {
		return errors.New("value must be non-empty")
	}
	return nil
}

func strMustNotContainOnlySpaces(s string) error {
	if len(strings.TrimSpace(s)) == 0 {
		return errors.New("value must not contain only spaces")
	}
	return nil
}

func pathMustExist(p string) error {
	if _, err := os.Stat(p); err != nil && os.IsNotExist(err) {
		return errors.New("path does not exist")
	}
	return nil
}

func pathMustBeDir(p string) error {
	if fi, err := os.Stat(p); err != nil {
		return err
	} else if !fi.IsDir() {
		return errors.New("path is not a directory")
	}
	return nil
}

func askName(wksp *config.Workspace) error {
	if dir, err := os.Getwd(); err != nil {
		return err
	} else if res, err := ui.Ask("Name", path.Base(dir), strMustNotContainOnlySpaces); err != nil {
		return err
	} else {
		wksp.Name = res
	}
	return nil
}

func askPath(wksp *config.Workspace) error {
	if dir, err := os.Getwd(); err != nil {
		return err
	} else if res, err := ui.Ask("Path", dir, strMustBeNonEmpty, strMustNotContainOnlySpaces, pathMustBeDir); err != nil {
		return err
	} else {
		wksp.Path = res
	}
	return nil
}

func askProjects(wksp *config.Workspace) error {
	var projects []config.Project = make([]config.Project, 0)
	var entries []fs.DirEntry
	var cwd string
	var err error
	if cwd, err = os.Getwd(); err != nil {
		return err
	} else if entries, err = os.ReadDir(cwd); err != nil {
		return err
	} else {
		for _, dir := range entries {
			if dir.IsDir() && !strings.HasPrefix(strings.TrimSpace(dir.Name()), ".") {
				projects = append(projects, config.Project{
					Name: dir.Name(),
					Path: filepath.Join(cwd, dir.Name()),
				})
			}
		}
	}
	fmt.Printf("Found %d projects:\n", len(projects))
	for _, proj := range projects {
		fmt.Printf("- %s\n", proj)
	}
	fmt.Printf("Add empty project to stop\n")
	done := false
	for !done {
		if res, err := ui.AskProject("Project", nil, nil); err != nil {
			return err
		} else if len(strings.TrimSpace(res.Name)) == 0 {
			done = true
		} else {
			projects = append(projects, *res)
		}
	}
	wksp.Projects = projects
	return nil
}

func askSourceControlType(wksp *config.Workspace) error {
	names := []string{}
	for _, s := range vcs.VersionControlSoftwares {
		names = append(names, s.Name())
	}
	if res, err := ui.Select("Source Control Type", names, nil); err != nil {
		return err
	} else {
		wksp.SourceControl = res
	}
	return nil
}

func askAuthor(wksp *config.Workspace) error {
	if currentUser, err := user.Current(); err != nil {
		return err
	} else {
		username := currentUser.Name
		if len(strings.TrimSpace(username)) == 0 {
			username = currentUser.Username
		}
		if author, err := ui.AskPerson("Author", &config.Person{Name: username}); err != nil {
			return err
		} else {
			wksp.Author = *author
		}
	}
	return nil
}

func askManager(wksp *config.Workspace) error {
	username := wksp.Author.Name
	if manager, err := ui.AskPerson("Manager", &config.Person{Name: username}); err != nil {
		return err
	} else {
		wksp.Manager = *manager
	}
	return nil
}

func askDeveloppers(wksp *config.Workspace) error {
	return nil
}

var (
	initCmd = &cobra.Command{
		Use:   "init [sub]",
		Short: "Initialize the current folder as a workspace",
		Long: `Initialize the current folder and turns it into a workspace.
This will write 'grlm.worspace.yaml' and will interactively ask a few questions.
`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			wksp := config.Workspace{}
			if err = askName(&wksp); err != nil {
				return err
			}
			if err = askPath(&wksp); err != nil {
				return err
			}
			if err = askProjects(&wksp); err != nil {
				return err
			}
			if err = askSourceControlType(&wksp); err != nil {
				return err
			}
			if err = askAuthor(&wksp); err != nil {
				return err
			}
			if err = askManager(&wksp); err != nil {
				return err
			}
			if err = askDeveloppers(&wksp); err != nil {
				return err
			}
			fmt.Printf("Write: %+v\n", wksp)
			return nil
		},
	}
)
