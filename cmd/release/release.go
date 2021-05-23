package release

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/project"
	"github.com/welschmorgan/go-release-manager/ui"
	"github.com/welschmorgan/go-release-manager/vcs"
	"github.com/welschmorgan/go-release-manager/version"
)

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
			state:          0,
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

func (r *Release) AcquireVersion() (string, error) {
	curVersion := "0.1"
	switch config.Get().AcquireVersionFrom {
	case "tags":
		if tags, err := r.Vc.ListTags(nil); err != nil {
			return "", err
		} else if len(tags) == 0 {
			return "", errors.New("cannot acquire version from tags, no tags yet")
		} else {
			curVersion = tags[len(tags)-1]
		}
	case "package":
		accessor, err := project.Open(r.Project.Path)
		if err != nil {
			return "", err
		}
		if curVersion, err = accessor.CurrentVersion(); err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("cannot acquire version from '%s', don't know what to do", config.Get().AcquireVersionFrom)
	}
	return curVersion, nil
}

func (r *Release) PrepareContext() error {
	// acquire current version
	var curVersion, curBranch string
	var err error
	if curVersion, err = r.AcquireVersion(); err != nil {
		return err
	}
	if curBranch, err = r.Vc.CurrentBranch(); err != nil {
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

	r.Context.state = ReleaseStarted

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

	r.Context.state |= ReleaseFinished
	return nil
}

func (r *Release) Step(n string) {
	fmt.Fprintf(os.Stderr, "[\033[1;34m*\033[0m] %s\n", n)
	config.Get().Indent = 1
}

func (r *Release) Undo() error {
	errs := []error{}

	if err := os.Chdir(r.Project.Path); err != nil {
		return err
	}

	fmt.Printf("[\033[1;31m-\033[0m] Undoing release %s for '%s' ...\n", r.Context.version, r.Project.Name)
	config.Get().Indent++
	for i := len(r.UndoActions) - 1; i >= 0; i-- {
		action := r.UndoActions[i]
		// if !action.Executed {
		fmt.Printf("%s[\033[1;34m*\033[0m] Undoing release step #%d: %s\n", strings.Repeat("\t", config.Get().Indent), i, action.Title)
		config.Get().Indent++
		if err := action.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\033[1;31merror\033[0m: %s\n", strings.Repeat("\t", config.Get().Indent), err.Error())
			errs = append(errs, err)
		}
		config.Get().Indent--
		// } else {
		// 	fmt.Printf("[\033[1;33m!\033[0m]\t- Undo of release step #%d already executed: %s\n", i, action.Title)
		// }
	}
	config.Get().Indent--
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
	if config.Get().Interactive {
		if ok, _ := ui.AskYN("Finish release now"); !ok {
			return errAbortRelease
		}
	}
	return nil
}

func (r *Release) UpdateRepository() error {
	var err error
	r.Step("Update repository")
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
	r.Step("Stash modifications")
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
	r.Step("Start release")
	r.Context.state |= ReleaseStartStarted
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
	r.Context.state |= ReleaseStartFinished
	return nil
}

func (r *Release) ReleaseFinish() error {
	// merge release branch into prod branch
	r.Step("Finish release")
	r.Context.state |= ReleaseFinishStarted
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

	// delete release branch
	if err := r.Vc.DeleteBranch(r.Context.releaseBranch, nil); err != nil {
		return err
	}
	r.Context.state |= ReleaseFinishFinished
	return nil
}

func (r *Release) BumpVersion() error {
	r.Step("Bump version")
	r.PushUndoAction("bump_version", r.Project.Path, r.Vc.Name(), map[string]interface{}{
		"oldVersion": r.Context.version,
		"newVersion": r.Context.nextVersion,
	})
	return nil
}
