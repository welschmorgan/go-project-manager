package release

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/fs"
	"github.com/welschmorgan/go-release-manager/ui"
	"github.com/welschmorgan/go-release-manager/version"
)

var errAbortRelease = errors.New("release aborted")

var Command = &cobra.Command{
	Use:   "release <major|minor|build|revision|preRelease|buildMetaTag> [OPTIONS...]",
	Short: "Release all projects included in this workspace",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires exactly 1 argument: release_type")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if typ, err := version.ParsePart(args[0]); err != nil {
			return err
		} else {
			config.Get().ReleaseType = typ
			return doRelease()
		}
	},
}

func doRelease() (err error) {
	if _, err = os.Stat(config.Get().WorkspacePath); err != nil && os.IsNotExist(err) {
		panic("Workspace has not been initialized, run `grlm init`")
	}

	releases := []*Release{}

	// cleanup release on ctrl-c
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGQUIT)
	rollback := func() {
		fmt.Fprintf(os.Stderr, "Rolling back all releases...\n")
		fs.DumpDirStack(os.Stderr)
		errs := []error{}
		for _, r := range releases {
			// fmt.Fprintf(os.Stderr, "  Rolling back %s v%s:\n", r.Project.Name, r.Context.version)
			// for _, u := range r.UndoActions {
			// 	fmt.Fprintf(os.Stderr, "    - Undo '%s' (path: %s, action: %s)\n", u.Title, u.Path, u.Name)
			// }
			if err := r.Undo(); err != nil {
				errs = append(errs, err)
			}
		}
		errStr := ""
		for _, e := range errs {
			if len(errStr) > 0 {
				errStr += "\n"
			}
			errStr += e.Error()
		}
		if len(errStr) > 0 {
			panic(errStr)
		}
	}
	go func() {
		<-sigs
		rollback()
		os.Exit(0)
	}()
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Printf("Recovered from: %v\n", r)
	// 		rollback()
	// 	}
	// }()

	for _, prj := range config.Get().Workspace.Projects {
		if r, err := NewRelease(prj); err != nil {
			return err
		} else {
			releases = append(releases, r)
			if err = r.Do(); err != nil {
				return err
			}
		}
	}

	println("Check if everything is OK, if it isn't, answering 'no' now will rollback what has been done.")
	if ok, err := ui.AskYN("Is everything ok"); err != nil || !ok {
		rollback()
	}
	return nil
}
