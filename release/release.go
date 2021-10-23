package release

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/fs"
	"github.com/welschmorgan/go-release-manager/log"
	"github.com/welschmorgan/go-release-manager/project"
	"github.com/welschmorgan/go-release-manager/ui"
	"github.com/welschmorgan/go-release-manager/vcs"
	"github.com/welschmorgan/go-release-manager/version"
	"gopkg.in/yaml.v2"
)

var errAbortRelease = errors.New("release aborted")

func ListUndos() (map[string][]UndoAction, error) {
	dir := config.Get().Workspace.Path.Join(".grlm", "releases").Expand()
	undoActions := map[string][]UndoAction{}
	var readDir func(dir string) (err error)

	readFile := func(name, path string, fi os.FileInfo) error {
		if content, err := os.ReadFile(path); err != nil {
			return fmt.Errorf("failed to load undo %s, %s", path, err.Error())
		} else {
			var releaseUndoActions = []UndoAction{}
			if err = yaml.Unmarshal(content, &releaseUndoActions); err != nil {
				return fmt.Errorf("failed to load undo %s, %s", path, err.Error())
			}
			undoActions[name] = releaseUndoActions
		}
		return nil
	}

	readDir = func(dir string) (err error) {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("failed to read directory: " + err.Error())
		}
		path := ""
		var fi os.FileInfo
		for _, e := range entries {
			path = filepath.Join(dir, e.Name())
			err = nil
			if fi, err = e.Info(); err != nil {
				log.Errorf("failed to retrieve file infos for %s, %s", path, err.Error())
			} else if fi.IsDir() {
				err = readDir(path)
			} else {
				err = readFile(e.Name(), path, fi)
			}
			if err != nil {
				return err
			}
		}
		return nil
	}

	if err := readDir(dir); err != nil {
		return nil, err
	}

	return undoActions, nil
}

type Release struct {
	Project     *config.Project            `yaml:"project,omitempty" json:"project,omitempty"`
	Vc          vcs.VersionControlSoftware `yaml:"-" json:"-"`
	Context     Context                    `yaml:"context,omitempty" json:"context,omitempty"`
	UndoActions []*UndoAction              `yaml:"undo_actions,omitempty" json:"undoActions,omitempty"`
}

func NewRelease(p *config.Project) (r *Release, err error) {
	r = &Release{
		Project: p,
		Context: Context{
			Date:           time.Now().UTC(),
			StartingBranch: "",
			ReleaseBranch:  config.Get().BranchNames["release"],
			DevBranch:      config.Get().BranchNames["development"],
			ProdBranch:     config.Get().BranchNames["production"],
			Version:        nil,
			NextVersion:    nil,
			HasRemotes:     false,
			State:          0,
			Accessor:       nil,
		},
		UndoActions: []*UndoAction{},
	}

	if r.Context.Accessor, err = project.Open(r.Project.Path); err != nil {
		return
	}

	if r.Vc, err = vcs.Open(r.Project.Path); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Release) PushUndoAction(name string, path fs.Path, vc string, params map[string]interface{}) error {
	if act, err := NewUndoAction(name, path, vc, params); err != nil {
		return err
	} else {
		act.Id = uint(len(r.UndoActions))
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
		if v, err = r.Context.Accessor.ReadVersion(); err != nil {
			return
		}
	default:
		return nil, fmt.Errorf("cannot acquire version from '%s', don't know what to do", config.Get().AcquireVersionFrom)
	}
	return v, nil
}

func (r *Release) PrepareContext() (err error) {
	// acquire current version

	if err = r.Project.Path.Chdir(); err != nil {
		return err
	}

	remotes := map[string]string{}
	if remotes, err = r.Vc.ListRemotes(nil); err != nil {
		return err
	} else if len(remotes) > 0 {
		r.Context.HasRemotes = true
	}
	if r.Context.HasRemotes {
		if err = r.Vc.FetchIndex(vcs.FetchIndexOptions{All: true, Tags: true, Force: true}); err != nil {
			return err
		}
	}

	if err = r.CheckoutAndPullBranch(config.Get().BranchNames["development"]); err != nil {
		return err
	}

	var curVersion version.Version
	var curBranch string
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
	fs.PutPathEnv("version", curVersion.String)
	fs.PutPathEnv("next_version", curVersion.String)
	r.Context.ReleaseBranch = fs.ExpandPath(r.Context.ReleaseBranch)
	r.Context.StartingBranch = curBranch
	r.Context.Version = curVersion
	r.Context.NextVersion = nextVersion

	tags := []string{}
	if tags, err = r.Vc.ListTags(nil); err != nil {
		return err
	}
	for _, tag := range tags {
		if strings.EqualFold(r.Context.Version.String(), strings.ToLower(tag)) {
			return fmt.Errorf("current version '%s' already tagged", r.Context.Version)
		}
		if strings.EqualFold(r.Context.Version.String(), strings.ToLower(tag)) {
			return fmt.Errorf("next version '%s' already tagged", r.Context.NextVersion)
		}
	}
	return nil
}

func (r *Release) Do() error {
	var err error
	if err = r.Project.Path.Chdir(); err != nil {
		return err
	}

	r.Context.State = ReleaseStarted

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

	if err = r.Write(); err != nil {
		return err
	}
	r.Context.State |= ReleaseFinished
	return nil
}

func (r *Release) Write() error {
	if data, err := yaml.Marshal(r); err != nil {
		return err
	} else {
		undosDir := config.Get().Workspace.Path.Join(".grlm", "releases").Expand()
		os.MkdirAll(undosDir, 0755)
		path := filepath.Join(undosDir, r.Context.Version.String()+".yaml")
		if err = os.WriteFile(path, data, 0755); err != nil {
			return err
		}
		log.Infof("Release saved: %s\n", path)
	}
	return nil
}

func (r *Release) Read() error {
	undosDir := config.Get().Workspace.Path.Join(".grlm", "releases").Expand()
	os.MkdirAll(undosDir, 0755)
	path := filepath.Join(undosDir, r.Context.Version.String()+".yaml")
	if data, err := os.ReadFile(path); err != nil {
		return err
	} else if err = yaml.Unmarshal(data, r); err != nil {
		return err
	}
	log.Infof("Release loaded: %s\n", path)
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

	if err := r.Project.Path.Chdir(); err != nil {
		return err
	}

	log.Debugf("[\033[1;31m-\033[0m] Undoing release %s for '%s' ...\n", r.Context.Version, r.Project.Name)
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
	if err = r.CheckoutAndPullBranch(r.Context.ProdBranch); err != nil {
		return err
	}
	if err = r.CheckoutAndPullBranch(r.Context.DevBranch); err != nil {
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
		message := fmt.Sprintf("Before release %s, on branch %s", r.Context.Version, r.Context.StartingBranch)
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
	if r.Context.OldBranch, err = r.Vc.CurrentBranch(); err != nil {
		return err
	}
	log.Debugf("Checkout %s at %s - %s", branch, hash, subj)
	if err = r.Vc.Checkout(branch, vcs.CheckoutOptions{CreateBranch: false}); err != nil {
		return err
	} else {
		r.PushUndoAction("checkout", r.Project.Path, r.Vc.Name(), map[string]interface{}{
			"oldBranch": r.Context.OldBranch,
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
	if r.Context.HasRemotes {
		if err := r.PullBranch(); err != nil {
			return err
		}
	}
	return nil
}

func (r *Release) PullTags() error {
	if r.Context.HasRemotes {
		if err := r.Vc.Pull(vcs.PullOptions{All: false, ListTags: true, Force: true}); err != nil {
			return err
		}
	}
	return nil
}

func (r *Release) ReleaseStart() error {
	var err error
	r.Step("Start release")
	r.Context.State |= ReleaseStartStarted
	if r.Context.OldBranch, err = r.Vc.CurrentBranch(); err != nil {
		return err
	} else if err = r.Vc.Checkout(r.Context.ReleaseBranch, vcs.CheckoutOptions{
		StartingPoint: r.Context.DevBranch,
		CreateBranch:  true,
	}); err != nil {
		return err
	}
	r.PushUndoAction("create_branch", r.Project.Path, r.Vc.Name(), map[string]interface{}{
		"oldBranch": r.Context.OldBranch,
		"newBranch": r.Context.ReleaseBranch,
	})
	r.PushUndoAction("checkout", r.Project.Path, r.Vc.Name(), map[string]interface{}{
		"oldBranch": r.Context.OldBranch,
		"newBranch": r.Context.ReleaseBranch,
	})
	r.Context.State |= ReleaseStartFinished
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
	r.Context.State |= ReleaseFinishStarted

	r.SubStep("Merge " + r.Context.ReleaseBranch + " into " + r.Context.ProdBranch)
	if err := r.Merge(r.Context.ReleaseBranch, r.Context.ProdBranch, vcs.MergeOptions{NoFastForward: true}); err != nil {
		return err
	}

	// tag prod branch
	r.SubStep("Tag " + r.Context.Version.String())
	if err := r.Vc.Tag(r.Context.Version.String(), vcs.TagOptions{Annotated: true, Message: fmt.Sprintf("Release %s: %s", r.Context.Version, "TODO")}); err != nil {
		return err
	}

	r.PushUndoAction("create_tag", r.Project.Path, r.Vc.Name(), map[string]interface{}{
		"name": r.Context.Version.String(),
	})

	r.SubStep("Merge tag " + r.Context.Version.String() + " into " + r.Context.DevBranch)
	// retro merge tag into dev branch
	if err := r.Merge(r.Context.Version.String(), r.Context.DevBranch, vcs.MergeOptions{NoFastForward: true}); err != nil {
		return err
	}

	// delete release branch
	r.SubStep("Delete release branch " + r.Context.ReleaseBranch)
	if err := r.Vc.DeleteBranch(r.Context.ReleaseBranch, nil); err != nil {
		return err
	}
	r.Context.State |= ReleaseFinishFinished
	return nil
}

func (r *Release) PrepareForNextSprint() (err error) {
	r.Step("Prepare for next sprint: %s", r.Context.NextVersion.String())
	r.SubStep("Checkout %s", r.Context.DevBranch)
	if err = r.CheckoutBranch(r.Context.DevBranch); err != nil {
		return
	}
	r.SubStep("Bump version: %s -> %s", r.Context.Version, r.Context.NextVersion)
	if err = r.Context.Accessor.WriteVersion(&r.Context.NextVersion); err != nil {
		return
	}
	r.PushUndoAction("bump_version", r.Project.Path, r.Vc.Name(), map[string]interface{}{
		"oldVersion": r.Context.Version,
		"newVersion": r.Context.NextVersion,
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
	subject := "Prepare for next sprint: " + r.Context.NextVersion.String()
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
