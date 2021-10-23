package release

import (
	"fmt"
	"time"

	"github.com/welschmorgan/go-release-manager/fs"
	"github.com/welschmorgan/go-release-manager/log"
	"github.com/welschmorgan/go-release-manager/vcs"
)

type UndoAction struct {
	Id            uint                       `yaml:"id" json:"id"`
	Date          time.Time                  `yaml:"date,omitempty" json:"date,omitempty"`
	Name          string                     `yaml:"name,omitempty" json:"name,omitempty"`
	Title         string                     `yaml:"title,omitempty" json:"title,omitempty"`
	Path          fs.Path                    `yaml:"path,omitempty" json:"path,omitempty"`
	SourceControl string                     `yaml:"source_control,omitempty" json:"source_control,omitempty"`
	VC            vcs.VersionControlSoftware `yaml:"-" json:"-"`
	Params        map[string]interface{}     `yaml:"params,omitempty" json:"params,omitempty"`
	Executed      bool                       `yaml:"executed,omitempty" json:"executed,omitempty"`
}

var undoActionParams = map[string][]string{
	"stash_save":    {"name"},
	"create_branch": {"newBranch", "oldBranch"},
	"checkout":      {"newBranch", "oldBranch"},
	"merge":         {"prevHead", "nextHead", "source", "target"},
	"create_tag":    {"name"},
	"bump_version":  {"oldVersion", "newVersion"},
	"commit":        {"branch", "prevHead", "nextHead", "subject"},
	"pull_branch":   {"branch", "prevHead", "nextHead"},
}

var undoActionParamHandlers = map[string]func(*UndoAction) error{
	"stash_save": func(u *UndoAction) error {
		u.Title = fmt.Sprintf("Create stash '%s'", u.Params["name"])
		return nil
	},
	"create_branch": func(u *UndoAction) error {
		u.Title = fmt.Sprintf("Create branch %s from %s", u.Params["newBranch"], u.Params["oldBranch"])
		return nil
	},
	"pull_branch": func(u *UndoAction) error {
		u.Title = fmt.Sprintf("Pull %s (%s -> %s)", u.Params["branch"], u.Params["prevHead"], u.Params["nextHead"])
		return nil
	},
	"checkout": func(u *UndoAction) error {
		u.Title = fmt.Sprintf("Switch branch from %s -> %s", u.Params["oldBranch"], u.Params["newBranch"])
		return nil
	},
	"merge": func(u *UndoAction) error {
		u.Title = fmt.Sprintf("Merge branch %s -> %s", u.Params["source"], u.Params["target"])
		return nil
	},
	"commit": func(u *UndoAction) error {
		u.Title = fmt.Sprintf("Commit staged files on %s (%s-> %s) - %s", u.Params["branch"], u.Params["prevHead"], u.Params["nextHead"], u.Params["subject"])
		return nil
	},
	"create_tag": func(u *UndoAction) error {
		u.Title = fmt.Sprintf("Create tag %s", u.Params["name"])
		return nil
	},
	"bump_version": func(u *UndoAction) error {
		u.Title = fmt.Sprintf("Bump version %s -> %s", u.Params["oldVersion"], u.Params["newVersion"])
		return nil
	},
}

func NewUndoAction(name string, path fs.Path, vc string, params map[string]interface{}) (*UndoAction, error) {
	act := &UndoAction{
		Id:            0,
		Date:          time.Now().UTC(),
		Name:          name,
		Title:         "",
		Path:          path,
		SourceControl: vc,
		Params:        params,
	}
	if vc, err := vcs.Open(path); err != nil {
		return nil, err
	} else {
		act.VC = vc
	}
	if _, ok := undoActionParams[name]; !ok {
		return nil, fmt.Errorf("unknown action '%s'", name)
	} else {
		if handler, ok := undoActionParamHandlers[name]; !ok {
			return nil, fmt.Errorf("unknown undo action param handler for '%s'", name)
		} else if err := handler(act); err != nil {
			return nil, err
		}
	}
	return act, nil
}

func (u *UndoAction) Run() error {
	defer func() { u.Executed = true }()
	if err := u.Path.Chdir(); err != nil {
		return err
	}
	switch u.Name {
	case "stash_save":
		return u.undoStashSave()
	case "create_branch":
		return u.undoCreateBranch()
	case "pull_branch":
		return u.undoPullBranch()
	case "checkout":
		return u.undoCheckout()
	case "merge":
		return u.undoMerge()
	case "create_tag":
		return u.undoTag()
	case "commit":
		return u.undoCommit()
	case "bump_version":
		return u.undoBumpVersion()
	default:
		return fmt.Errorf("unknown undo action '%s'", u.Name)
	}
}

func (u *UndoAction) undoPullBranch() (err error) {
	branch := u.Params["branch"].(string)
	prevHead := u.Params["prevHead"].(string)
	nextHead := u.Params["nextHead"].(string)
	log.Debugf("Checkout %s", branch)
	if err = u.VC.Checkout(branch, nil); err != nil {
		return err
	}
	log.Debugf("Reset HEAD from %s to %s", nextHead, prevHead)
	err = u.VC.Reset(vcs.ResetOptions{
		Hard:   true,
		Commit: prevHead,
	})
	return err
}

func (u *UndoAction) undoStashSave() error {
	log.Debugln("Pop stash")
	_, err := u.VC.Stash(vcs.StashOptions{
		Pop: true,
	})
	return err
}

func (u *UndoAction) undoCreateBranch() error {
	// oldBranch := u.Params["oldBranch"].(string)
	newBranch := u.Params["newBranch"].(string)
	if branches, err := u.VC.ListBranches(nil); err != nil {
		return err
	} else {
		alreadyDeleted := true
		for _, b := range branches {
			if b == newBranch {
				alreadyDeleted = false
				break
			}
		}
		if !alreadyDeleted {
			log.Debugf("Delete branch '%s'")
			if err := u.VC.DeleteBranch(newBranch, nil); err != nil {
				return err
			}
		} else {
			log.Errorf("\t\tBranch '%s' has already been deleted\n", newBranch)
		}
	}
	return nil
}

func (u *UndoAction) undoCheckout() error {
	oldBranch := u.Params["oldBranch"].(string)
	// newBranch := u.Params["newBranch"].(string)
	log.Debugf("Checkout %s", oldBranch)
	return u.VC.Checkout(oldBranch, nil)
}

func (u *UndoAction) undoMerge() error {
	// source := u.Params["source"].(string)
	target := u.Params["target"].(string)
	prevHead := u.Params["prevHead"].(string)
	log.Debugf("Checkout %s", target)
	if err := u.VC.Checkout(target, nil); err != nil {
		return err
	} else {
		log.Debugf("Reset HEAD to %s", prevHead)
		if err := u.VC.Reset(vcs.ResetOptions{Commit: prevHead, Hard: true}); err != nil {
			return err
		}
	}
	return nil
}

func (u *UndoAction) undoTag() error {
	name := u.Params["name"].(string)
	log.Debugf("Delete tag %s", name)
	return u.VC.Tag(name, vcs.TagOptions{
		Delete: true,
	})
}

func (u *UndoAction) undoCommit() error {
	branch := u.Params["branch"].(string)
	subject := u.Params["subject"].(string)
	prevHead := u.Params["prevHead"].(string)
	nextHead := u.Params["nextHead"].(string)
	log.Debugf("Delete commit on %s %s: %s", branch, nextHead, subject)
	return u.VC.Reset(vcs.ResetOptions{Hard: true, Commit: prevHead})
}

func (u *UndoAction) undoBumpVersion() error {
	// oldVersion := u.Params["oldVersion"].(string)
	// newVersion := u.Params["newVersion"].(string)
	return nil
}
