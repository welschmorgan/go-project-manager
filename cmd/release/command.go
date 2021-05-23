package release

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/project"
	"github.com/welschmorgan/go-release-manager/ui"
	"github.com/welschmorgan/go-release-manager/vcs"
)

// The release context
type Context struct {
	startingBranch string // The branch the repository was on before starting release
	oldBranch      string // The branch before checking out the current one
	releaseBranch  string // The release branch
	devBranch      string // The development branch
	prodBranch     string // The production branch
	version        string // The version the project is in
	nextVersion    string // The next version the project will be in after release
	hasRemotes     bool   // Wether the repository has remotes or not
}

var errAbortRelease = errors.New("release aborted")
var context Context
var undoActions []func() error

var Command = &cobra.Command{
	Use:   "release [OPTIONS...]",
	Short: "Release all projects included in this workspace",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if _, err = os.Stat(config.Get().WorkspacePath); err != nil && os.IsNotExist(err) {
			panic("Workspace has not been initialized, run `grlm init`")
		}
		for _, prj := range config.Get().Workspace.Projects {
			if err = release(prj); err != nil {
				return err
			}
		}
		return nil
	},
}

func release(p *config.Project) (err error) {
	if err = os.Chdir(p.Path); err != nil {
		return err
	}
	var vc vcs.VersionControlSoftware
	if vc, err = vcs.Open(p.Path); err != nil {
		return err
	}
	// cleanup release on ctrl-c
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		if err := abortRelease(p, vc); err != nil {
			panic(err.Error())
		}
		os.Exit(0)
	}()

	// acquire current version
	version := "0.1"
	switch config.Get().AcquireVersionFrom {
	case "tags":
		if tags, err := vc.ListTags(nil); err != nil {
			return err
		} else if len(tags) == 0 {
			return errors.New("cannot acquire version from tags, no tags yet")
		} else {
			version = tags[len(tags)-1]
		}
	case "package":
		accessor, err := project.Open(p.Path)
		if err != nil {
			return err
		}
		if version, err = accessor.CurrentVersion(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("cannot acquire version from '%s', don't know what to do", config.Get().AcquireVersionFrom)
	}
	curBranch, err := vc.CurrentBranch()
	if err != nil {
		return err
	}
	context = Context{
		startingBranch: curBranch,
		releaseBranch:  strings.ReplaceAll(config.Get().BranchNames["release"], "$VERSION", version),
		devBranch:      config.Get().BranchNames["development"],
		prodBranch:     config.Get().BranchNames["production"],
		version:        version,
		nextVersion:    "",
		hasRemotes:     false,
	}

	remotes := map[string]string{}
	if remotes, err = vc.ListRemotes(nil); err != nil {
		return err
	} else if len(remotes) > 0 {
		context.hasRemotes = true
	} else {
		fmt.Printf("[\033[1;33m!\033[0m] Project has not remotes, won't push or pull\n")
	}
	// stash modifications
	if err = stashModifications(p, vc); err != nil {
		return err
	}

	// checkout development and production branches
	if err = updateRepository(p, vc); err != nil {
		return err
	}

	// start release
	if err = releaseStart(p, vc); err != nil {
		return err
	}

	// wait for user to manually edit release

	if err = waitUserToConfirm(p, vc); err != nil {
		return err
	}

	// finish release
	if err = releaseFinish(p, vc); err != nil {
		return err
	}
	if err = bumpVersion(p, vc); err != nil {
		return err
	}
	return nil
}

func abortRelease(p *config.Project, v vcs.VersionControlSoftware) error {
	errs := []error{}
	for _, action := range undoActions {
		if err := action(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		errStr := "there were errors while undoing release:\n"
		for _, e := range errs {
			errStr += fmt.Sprintf(" - %s\n", e.Error())
		}
		return errors.New(errStr)
	}
	return nil
}

func waitUserToConfirm(p *config.Project, v vcs.VersionControlSoftware) error {
	ok, _ := ui.AskYN("Finish release now")
	if !ok {
		return errAbortRelease
	}
	return nil
}

func updateRepository(p *config.Project, v vcs.VersionControlSoftware) error {
	var err error
	if err = checkoutAndPullBranch(p, v, context.prodBranch); err != nil {
		return err
	}
	if err = checkoutAndPullBranch(p, v, context.devBranch); err != nil {
		return err
	}
	if err = pullTags(p, v); err != nil {
		return err
	}
	return nil
}

func undoStashSave(vc vcs.VersionControlSoftware) func() error {
	return func() error {
		_, err := vc.Stash(vcs.StashOptions{
			Pop: true,
		})
		return err
	}
}

func undoCreateBranch(vc vcs.VersionControlSoftware, oldBranch, newBranch string) func() error {
	return func() error {
		if err := vc.DeleteBranch(newBranch, nil); err != nil {
			return err
		}
		return nil
	}
}

func undoCheckout(vc vcs.VersionControlSoftware, oldBranch, newBranch string) func() error {
	return func() error {
		return vc.Checkout(oldBranch, nil)
	}
}

func undoMerge(vc vcs.VersionControlSoftware, source, target string) func() error {
	return func() error {
		if err := vc.Checkout(target, nil); err != nil {
			return err
		}
		if err := vc.Reset(vcs.ResetOptions{
			Commit: "HEAD~1",
			Hard:   true,
		}); err != nil {
			return err
		}
		return nil
	}
}

func undoTag(vc vcs.VersionControlSoftware, name string) func() error {
	return func() error {
		return vc.Tag(name, vcs.TagOptions{
			Delete: true,
		})
	}
}

func undoBumpVersion(vc vcs.VersionControlSoftware, oldVersion, newVersion string) func() error {
	return func() error {
		return nil
	}
}

func stashModifications(p *config.Project, v vcs.VersionControlSoftware) error {
	oldDryRun := config.Get().DryRun
	config.Get().DryRun = false
	status, err := v.Status(vcs.StatusOptions{Short: true})
	config.Get().DryRun = oldDryRun
	if err != nil {
		return err
	} else if len(status) != 0 {
		message := fmt.Sprintf("Before release %s, on branch %s", context.version, context.startingBranch)
		fmt.Printf("The current repository is dirty:\n%v\n\t-> stashed under '%s'\n", status, message)
		if _, err := v.Stash(vcs.StashOptions{
			Save:             true,
			IncludeUntracked: true,
			Message:          message,
		}); err != nil {
			return err
		}
		undoActions = append(undoActions, undoStashSave(v))
	}
	return nil
}

func checkoutBranch(p *config.Project, v vcs.VersionControlSoftware, branch string) error {
	var err error
	if context.oldBranch, err = v.CurrentBranch(); err != nil {
		return err
	} else if err = v.Checkout(branch, vcs.CheckoutOptions{CreateBranch: false}); err != nil {
		return err
	} else {
		undoActions = append(undoActions, undoCheckout(v, context.oldBranch, branch))
	}
	return nil
}

func pullBranch(p *config.Project, v vcs.VersionControlSoftware) error {
	if err := v.Pull(vcs.PullOptions{All: false, ListTags: false, Force: false}); err != nil {
		return err
	}
	return nil
}

func checkoutAndPullBranch(p *config.Project, v vcs.VersionControlSoftware, branch string) error {
	if err := checkoutBranch(p, v, branch); err != nil {
		return err
	}
	if context.hasRemotes {
		if err := pullBranch(p, v); err != nil {
			return err
		}
	}
	return nil
}

func pullTags(p *config.Project, v vcs.VersionControlSoftware) error {
	if context.hasRemotes {
		if err := v.Pull(vcs.PullOptions{All: false, ListTags: true, Force: true}); err != nil {
			return err
		}
	}
	return nil
}

func releaseStart(p *config.Project, v vcs.VersionControlSoftware) error {
	var err error
	if context.oldBranch, err = v.CurrentBranch(); err != nil {
		return err
	} else if err = v.Checkout(context.releaseBranch, vcs.CheckoutOptions{
		StartingPoint: context.devBranch,
		CreateBranch:  true,
	}); err != nil {
		return err
	}
	undoActions = append(undoActions, undoCheckout(v, context.oldBranch, context.releaseBranch), undoCreateBranch(v, context.oldBranch, context.releaseBranch))
	return nil
}

func releaseFinish(p *config.Project, v vcs.VersionControlSoftware) error {
	// merge release branch into prod branch
	if err := v.Merge(context.releaseBranch, context.prodBranch, vcs.MergeOptions{NoFastForward: true}); err != nil {
		return err
	}
	undoActions = append(undoActions, undoMerge(v, context.releaseBranch, context.prodBranch))
	// tag prod branch
	if err := v.Tag(context.version, vcs.TagOptions{Annotated: true, Message: fmt.Sprintf("Release %s: %s", context.version, "TODO")}); err != nil {
		return err
	}
	undoActions = append(undoActions, undoTag(v, context.version))
	// retro merge tag into dev branch
	if err := v.Merge(context.version, context.devBranch, vcs.MergeOptions{NoFastForward: true}); err != nil {
		return err
	}
	undoActions = append(undoActions, undoMerge(v, context.version, context.devBranch))
	return nil
}

func bumpVersion(p *config.Project, v vcs.VersionControlSoftware) error {
	undoActions = append(undoActions, undoBumpVersion(v, context.version, context.nextVersion))
	return nil
}
