package release

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/log"
	"github.com/welschmorgan/go-release-manager/project"
	"github.com/welschmorgan/go-release-manager/ui"
	"github.com/welschmorgan/go-release-manager/vcs"
	"github.com/welschmorgan/go-release-manager/version"
	"gopkg.in/yaml.v2"
)

var errAbortRelease = errors.New("release aborted")

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
			version:        nil,
			nextVersion:    nil,
			hasRemotes:     false,
			state:          0,
			accessor:       nil,
		},
		UndoActions: []*UndoAction{},
	}

	if r.Context.accessor, err = project.Open(r.Project.Path); err != nil {
		return
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

func (r *Release) AcquireVersion() (v version.Version, err error) {
	switch config.Get().AcquireVersionFrom {
	case "tags":
		if tags, err := r.Vc.ListTags(nil); err != nil {
			return nil, err
		} else if len(tags) == 0 {
			return nil, errors.New("cannot acquire version from tags, no tags yet")
		} else {
			tag := tags[len(tags)-1]
			if v = version.Parse(tag); v == nil {
				return nil, fmt.Errorf("failed to parse version from '%s'", tag)
			}
		}
	case "package":
		if v, err = r.Context.accessor.ReadVersion(); err != nil {
			return
		}
	default:
		return nil, fmt.Errorf("cannot acquire version from '%s', don't know what to do", config.Get().AcquireVersionFrom)
	}
	return v, nil
}

func (r *Release) PrepareContext() error {
	// acquire current version
	var curVersion version.Version
	var curBranch string
	var err error
	if curVersion, err = r.AcquireVersion(); err != nil {
		return err
	}
	if curBranch, err = r.Vc.CurrentBranch(); err != nil {
		return err
	}
	// duplicate current version
	nextVersion := version.Clone(curVersion)
	// increment it
	if err := nextVersion.Increment(config.Get().ReleaseType, 1); err != nil {
		return err
	}
	r.Context.releaseBranch = strings.ReplaceAll(r.Context.releaseBranch, "$VERSION", curVersion.String())
	r.Context.startingBranch = curBranch
	r.Context.version = curVersion
	r.Context.nextVersion = nextVersion
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

	if err = r.PrepareForNextSprint(); err != nil {
		return err
	}

	if err = r.WriteUndos(); err != nil {
		return err
	}
	r.Context.state |= ReleaseFinished
	return nil
}

func (r *Release) WriteUndos() error {
	if data, err := yaml.Marshal(r.UndoActions); err != nil {
		return err
	} else {
		dir := filepath.Join(config.Get().Workspace.Path(), ".grlm-undos")
		numFiles := 0
		if dirEntries, err := os.ReadDir(dir); err != nil {
			return err
		} else {
			for _, de := range dirEntries {
				if !de.IsDir() {
					numFiles++
				}
			}
		}
		os.MkdirAll(dir, 0755)
		path := filepath.Join(dir, fmt.Sprintf("%04d", numFiles)+"-"+r.Context.version.String()+".yaml")
		if err = os.WriteFile(path, data, 0755); err != nil {
			return err
		}
		log.Infof("Undo actions written at: %s\n", path)
	}
	return nil
}
func (r *Release) Step(fmt string, a ...interface{}) {
	log.Infof("[\033[1;34m*\033[0m] "+fmt+"\n", a...)
	config.Get().Indent = 1
}

func (r *Release) SubStep(fmt string, a ...interface{}) {
	log.Infof("[\033[1;34m**\033[0m] "+fmt+"\n", a...)
	config.Get().Indent++
}

func (r *Release) Undo() error {
	errs := []error{}

	if err := os.Chdir(r.Project.Path); err != nil {
		return err
	}

	log.Debugf("[\033[1;31m-\033[0m] Undoing release %s for '%s' ...\n", r.Context.version, r.Project.Name)
	config.Get().Indent++
	for i := len(r.UndoActions) - 1; i >= 0; i-- {
		action := r.UndoActions[i]
		// if !action.Executed {
		log.Debugf("[\033[1;34m*\033[0m] Undoing release step #%d: %s - path = %s\n", i, action.Title, action.Path)
		config.Get().Indent++
		if err := action.Run(); err != nil {
			log.Errorf("\033[1;31merror\033[0m: %s\n", err.Error())
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
	hash := ""
	subj := ""

	if hash, subj, err = r.Vc.CurrentCommit(vcs.CurrentCommitOptions{ShortHash: true}); err != nil {
		return err
	}
	if r.Context.oldBranch, err = r.Vc.CurrentBranch(); err != nil {
		return err
	}
	log.Debugf("Checkout %s at %s - %s", branch, hash, subj)
	if err = r.Vc.Checkout(branch, vcs.CheckoutOptions{CreateBranch: false}); err != nil {
		return err
	} else {
		r.PushUndoAction("checkout", r.Project.Path, r.Vc.Name(), map[string]interface{}{
			"oldBranch": r.Context.oldBranch,
			"newBranch": branch,
		})
	}
	return nil
}

func (r *Release) PullBranch() (err error) {
	branch := ""
	prevHead, nextHead := "", ""
	if branch, err = r.Vc.CurrentBranch(); err != nil {
		return err
	}
	if prevHead, _, err = r.Vc.CurrentCommit(vcs.CurrentCommitOptions{ShortHash: true}); err != nil {
		return err
	}
	if err := r.Vc.Pull(vcs.PullOptions{All: false, ListTags: false, Force: false}); err != nil {
		return err
	}
	if nextHead, _, err = r.Vc.CurrentCommit(vcs.CurrentCommitOptions{ShortHash: true}); err != nil {
		return err
	}
	log.Debugf("Pull %s (was %s, now %s)", branch, prevHead, nextHead)
	r.PushUndoAction("pull_branch", r.Project.Path, r.Vc.Name(), map[string]interface{}{
		"branch":   branch,
		"prevHead": prevHead,
		"nextHead": nextHead,
	})

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

func (r *Release) Merge(source, dest string, options vcs.MergeOptions) error {
	prevHead, nextHead := "", ""
	var err error
	if prevHead, _, err = r.Vc.CurrentCommit(vcs.CurrentCommitOptions{ShortHash: true}); err != nil {
		return err
	}
	if err := r.Vc.Merge(source, dest, options); err != nil {
		return err
	}
	if nextHead, _, err = r.Vc.CurrentCommit(vcs.CurrentCommitOptions{ShortHash: true}); err != nil {
		return err
	}
	r.PushUndoAction("merge", r.Project.Path, r.Vc.Name(), map[string]interface{}{
		"prevHead": prevHead,
		"nextHead": nextHead,
		"source":   source,
		"target":   dest,
	})
	return nil
}

func (r *Release) ReleaseFinish() error {
	// merge release branch into prod branch
	r.Step("Finish release")
	r.Context.state |= ReleaseFinishStarted

	r.SubStep("Merge " + r.Context.releaseBranch + " into " + r.Context.prodBranch)
	if err := r.Merge(r.Context.releaseBranch, r.Context.prodBranch, vcs.MergeOptions{NoFastForward: true}); err != nil {
		return err
	}

	// tag prod branch
	r.SubStep("Tag " + r.Context.version.String())
	if err := r.Vc.Tag(r.Context.version.String(), vcs.TagOptions{Annotated: true, Message: fmt.Sprintf("Release %s: %s", r.Context.version, "TODO")}); err != nil {
		return err
	}

	r.PushUndoAction("create_tag", r.Project.Path, r.Vc.Name(), map[string]interface{}{
		"name": r.Context.version.String(),
	})

	r.SubStep("Merge tag " + r.Context.version.String() + " into " + r.Context.devBranch)
	// retro merge tag into dev branch
	if err := r.Merge(r.Context.version.String(), r.Context.devBranch, vcs.MergeOptions{NoFastForward: true}); err != nil {
		return err
	}

	// delete release branch
	r.SubStep("Delete release branch " + r.Context.releaseBranch)
	if err := r.Vc.DeleteBranch(r.Context.releaseBranch, nil); err != nil {
		return err
	}
	r.Context.state |= ReleaseFinishFinished
	return nil
}

func (r *Release) PrepareForNextSprint() (err error) {
	r.Step("Prepare for next sprint: %s", r.Context.nextVersion.String())
	r.SubStep("Checkout %s", r.Context.devBranch)
	if err = r.CheckoutBranch(r.Context.devBranch); err != nil {
		return
	}
	r.SubStep("Bump version: %s -> %s", r.Context.version, r.Context.nextVersion)
	if err = r.Context.accessor.WriteVersion(&r.Context.nextVersion); err != nil {
		return
	}
	r.PushUndoAction("bump_version", r.Project.Path, r.Vc.Name(), map[string]interface{}{
		"oldVersion": r.Context.version,
		"newVersion": r.Context.nextVersion,
	})
	r.SubStep("Stage & Commit")
	if err = r.Vc.Stage(vcs.StageOptions{All: true}); err != nil {
		return
	}
	branch := ""
	prevHead := ""
	nextHead := ""
	if branch, err = r.Vc.CurrentBranch(); err != nil {
		return
	}
	if prevHead, _, err = r.Vc.CurrentCommit(vcs.CurrentCommitOptions{ShortHash: true}); err != nil {
		return
	}
	subject := "Prepare for next sprint: " + r.Context.nextVersion.String()
	if err = r.Vc.Commit(vcs.CommitOptions{Message: subject, AllowEmpty: true}); err != nil {
		return
	}
	if nextHead, _, err = r.Vc.CurrentCommit(vcs.CurrentCommitOptions{ShortHash: true}); err != nil {
		return
	}
	r.PushUndoAction("commit", r.Project.Path, r.Vc.Name(), map[string]interface{}{
		"branch":   branch,
		"subject":  subject,
		"prevHead": prevHead,
		"nextHead": nextHead,
	})

	return nil
}
