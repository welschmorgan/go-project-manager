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
	"github.com/welschmorgan/go-release-manager/version"
)

var errAbortRelease = errors.New("release aborted")

var Command = &cobra.Command{
	Use:   "release [OPTIONS...]",
	Short: "Release all projects included in this workspace",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if _, err = os.Stat(config.Get().WorkspacePath); err != nil && os.IsNotExist(err) {
			panic("Workspace has not been initialized, run `grlm init`")
		}

		releases := []*Release{}

		// cleanup release on ctrl-c
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigs
			errs := []error{}
			for _, r := range releases {
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
			os.Exit(0)
		}()

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
		return nil
	},
}

type Release struct {
	Project     *config.Project
	Vc          vcs.VersionControlSoftware
	Context     Context
	UndoActions []*UndoAction
}

func NewRelease(p *config.Project) (r *Release, err error) {
	r = &Release{
		Project: p,
		Context: Context{
			startingBranch: "",
			releaseBranch:  config.Get().BranchNames["release"],
			devBranch:      config.Get().BranchNames["development"],
			prodBranch:     config.Get().BranchNames["production"],
			version:        "",
			nextVersion:    "",
			hasRemotes:     false,
		},
		UndoActions: []*UndoAction{},
	}

	if r.Vc, err = vcs.Open(r.Project.Path); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Release) PushUndoAction(name, path, vc string, params map[string]interface{}) error {
	if act, err := NewUndoAction(name, path, vc, params); err != nil {
		return err
	} else {
		r.UndoActions = append(r.UndoActions, act)
	}
	return nil
}

func (r *Release) PrepareContext() error {
	// acquire current version
	curVersion := "0.1"
	switch config.Get().AcquireVersionFrom {
	case "tags":
		if tags, err := r.Vc.ListTags(nil); err != nil {
			return err
		} else if len(tags) == 0 {
			return errors.New("cannot acquire version from tags, no tags yet")
		} else {
			curVersion = tags[len(tags)-1]
		}
	case "package":
		accessor, err := project.Open(r.Project.Path)
		if err != nil {
			return err
		}
		if curVersion, err = accessor.CurrentVersion(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("cannot acquire version from '%s', don't know what to do", config.Get().AcquireVersionFrom)
	}
	curBranch, err := r.Vc.CurrentBranch()
	if err != nil {
		return err
	}
	nextVersion := version.Parse(curVersion)
	if err := nextVersion.Increment(0, 1); err != nil {
		return err
	}
	r.Context.releaseBranch = strings.ReplaceAll(r.Context.releaseBranch, "$VERSION", curVersion)
	r.Context.startingBranch = curBranch
	r.Context.version = curVersion
	r.Context.nextVersion = nextVersion.String()
	r.Context.hasRemotes = false

	remotes := map[string]string{}
	if remotes, err = r.Vc.ListRemotes(nil); err != nil {
		return err
	} else if len(remotes) > 0 {
		r.Context.hasRemotes = true
	}
	return nil
}

func (r *Release) Do() error {
	var err error
	if err = os.Chdir(r.Project.Path); err != nil {
		return err
	}

	defer func() {
		// TODO DELETE THIS
		if err := r.Undo(); err != nil {
			panic(err.Error())
		}
	}()

	if err = r.CheckoutAndPullBranch(config.Get().BranchNames["development"]); err != nil {
		return err
	}

	if err = r.PrepareContext(); err != nil {
		return err
	}

	// stash modifications
	if err = r.StashModifications(); err != nil {
		return err
	}

	// checkout development and production branches
	if err = r.UpdateRepository(); err != nil {
		return err
	}

	// start release
	if err = r.ReleaseStart(); err != nil {
		return err
	}

	// wait for user to manually edit release
	if err = r.WaitUserToConfirm(); err != nil {
		return err
	}

	// finish release
	if err = r.ReleaseFinish(); err != nil {
		return err
	}
	if err = r.BumpVersion(); err != nil {
		return err
	}
	return nil
}

func (r *Release) Undo() error {
	errs := []error{}

	if err := os.Chdir(r.Project.Path); err != nil {
		return err
	}

	fmt.Printf("[\033[1;31m-\033[0m] Undoing release %s for '%s' ...\n", r.Context.version, r.Project.Name)
	for i := len(r.UndoActions) - 1; i >= 0; i-- {
		action := r.UndoActions[i]
		fmt.Printf("[\033[1;31m-\033[0m]\t- Undoing release step #%d: %s\n", i, action.Title)
		if err := action.Run(); err != nil {
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

func (r *Release) WaitUserToConfirm() error {
	ok, _ := ui.AskYN("Finish release now")
	if !ok {
		return errAbortRelease
	}
	return nil
}

func (r *Release) UpdateRepository() error {
	var err error
	if err = r.CheckoutAndPullBranch(r.Context.prodBranch); err != nil {
		return err
	}
	if err = r.CheckoutAndPullBranch(r.Context.devBranch); err != nil {
		return err
	}
	if err = r.PullTags(); err != nil {
		return err
	}
	return nil
}

func (r *Release) StashModifications() error {
	oldDryRun := config.Get().DryRun
	config.Get().DryRun = false
	status, err := r.Vc.Status(vcs.StatusOptions{Short: true})
	config.Get().DryRun = oldDryRun
	if err != nil {
		return err
	} else if len(status) != 0 {
		message := fmt.Sprintf("Before release %s, on branch %s", r.Context.version, r.Context.startingBranch)
		fmt.Printf("The current repository is dirty:\n%v\n\t-> stashed under '%s'\n", status, message)
		if _, err := r.Vc.Stash(vcs.StashOptions{
			Save:             true,
			IncludeUntracked: true,
			Message:          message,
		}); err != nil {
			return err
		}
		r.PushUndoAction("stash_save", r.Project.Path, r.Vc.Name(), map[string]interface{}{"name": message})
	}
	return nil
}

func (r *Release) CheckoutBranch(branch string) error {
	var err error
	if r.Context.oldBranch, err = r.Vc.CurrentBranch(); err != nil {
		return err
	} else if err = r.Vc.Checkout(branch, vcs.CheckoutOptions{CreateBranch: false}); err != nil {
		return err
	} else {
		r.PushUndoAction("checkout", r.Project.Path, r.Vc.Name(), map[string]interface{}{
			"oldBranch": r.Context.oldBranch,
			"newBranch": branch,
		})
	}
	return nil
}

func (r *Release) PullBranch() error {
	if err := r.Vc.Pull(vcs.PullOptions{All: false, ListTags: false, Force: false}); err != nil {
		return err
	}
	return nil
}

func (r *Release) CheckoutAndPullBranch(branch string) error {
	if err := r.CheckoutBranch(branch); err != nil {
		return err
	}
	if r.Context.hasRemotes {
		if err := r.PullBranch(); err != nil {
			return err
		}
	}
	return nil
}

func (r *Release) PullTags() error {
	if r.Context.hasRemotes {
		if err := r.Vc.Pull(vcs.PullOptions{All: false, ListTags: true, Force: true}); err != nil {
			return err
		}
	}
	return nil
}

func (r *Release) ReleaseStart() error {
	var err error
	if r.Context.oldBranch, err = r.Vc.CurrentBranch(); err != nil {
		return err
	} else if err = r.Vc.Checkout(r.Context.releaseBranch, vcs.CheckoutOptions{
		StartingPoint: r.Context.devBranch,
		CreateBranch:  true,
	}); err != nil {
		return err
	}

	r.PushUndoAction("create_branch", r.Project.Path, r.Vc.Name(), map[string]interface{}{
		"oldBranch": r.Context.oldBranch,
		"newBranch": r.Context.releaseBranch,
	})
	r.PushUndoAction("checkout", r.Project.Path, r.Vc.Name(), map[string]interface{}{
		"oldBranch": r.Context.oldBranch,
		"newBranch": r.Context.releaseBranch,
	})
	return nil
}

func (r *Release) ReleaseFinish() error {
	// merge release branch into prod branch
	if err := r.Vc.Merge(r.Context.releaseBranch, r.Context.prodBranch, vcs.MergeOptions{NoFastForward: true}); err != nil {
		return err
	}
	r.PushUndoAction("merge", r.Project.Path, r.Vc.Name(), map[string]interface{}{
		"source": r.Context.releaseBranch,
		"target": r.Context.prodBranch,
	})
	// tag prod branch
	if err := r.Vc.Tag(r.Context.version, vcs.TagOptions{Annotated: true, Message: fmt.Sprintf("Release %s: %s", r.Context.version, "TODO")}); err != nil {
		return err
	}

	r.PushUndoAction("create_tag", r.Project.Path, r.Vc.Name(), map[string]interface{}{
		"name": r.Context.version,
	})
	// retro merge tag into dev branch
	if err := r.Vc.Merge(r.Context.version, r.Context.devBranch, vcs.MergeOptions{NoFastForward: true}); err != nil {
		return err
	}
	r.PushUndoAction("merge", r.Project.Path, r.Vc.Name(), map[string]interface{}{
		"source": r.Context.version,
		"target": r.Context.devBranch,
	})
	return nil
}

func (r *Release) BumpVersion() error {
	r.PushUndoAction("bump_version", r.Project.Path, r.Vc.Name(), map[string]interface{}{
		"oldVersion": r.Context.version,
		"newVersion": r.Context.nextVersion,
	})
	return nil
}
