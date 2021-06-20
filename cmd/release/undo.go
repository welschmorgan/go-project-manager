package release

import (
	"fmt"
	"os"

	"github.com/welschmorgan/go-release-manager/vcs"
)

type UndoAction struct {
	Name          string
	Title         string
	Path          string
	SourceControl string
	VC            vcs.VersionControlSoftware
	Params        map[string]interface{}
	Executed      bool
}

var undoActionParams = map[string][]string{
	"stash_save":    {"name"},
	"create_branch": {"newBranch", "oldBranch"},
	"checkout":      {"newBranch", "oldBranch"},
	"merge":         {"source", "target"},
	"create_tag":    {"name"},
	"bump_version":  {"oldVersion", "newVersion"},
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
	"checkout": func(u *UndoAction) error {
		u.Title = fmt.Sprintf("Switch branch from %s -> %s", u.Params["oldBranch"], u.Params["newBranch"])
		return nil
	},
	"merge": func(u *UndoAction) error {
		u.Title = fmt.Sprintf("Merge branch %s -> %s", u.Params["source"], u.Params["target"])
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

func NewUndoAction(name, path, vc string, params map[string]interface{}) (*UndoAction, error) {
	act := &UndoAction{
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
	if err := os.Chdir(u.Path); err != nil {
		return err
	}
	switch u.Name {
	case "stash_save":
		return u.undoStashSave()
	case "create_branch":
		return u.undoCreateBranch()
	case "checkout":
		return u.undoCheckout()
	case "merge":
		return u.undoMerge()
	case "create_tag":
		return u.undoTag()
	case "bump_version":
		return u.undoBumpVersion()
	default:
		return fmt.Errorf("unknown undo action '%s'", u.Name)
	}
}

func (u *UndoAction) undoStashSave() error {
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
			if err := u.VC.DeleteBranch(newBranch, nil); err != nil {
				return err
			}
		} else {
			fmt.Fprintf(os.Stderr, "\t\tBranch '%s' has already been deleted\n", newBranch)
		}
	}
	return nil
}

func (u *UndoAction) undoCheckout() error {
	oldBranch := u.Params["oldBranch"].(string)
	// newBranch := u.Params["newBranch"].(string)
	return u.VC.Checkout(oldBranch, nil)
}

func (u *UndoAction) undoMerge() error {
	// source := u.Params["source"].(string)
	target := u.Params["target"].(string)
	if err := u.VC.Checkout(target, nil); err != nil {
		return err
	} else if err := u.VC.Reset(vcs.ResetOptions{Commit: "HEAD~1", Hard: true}); err != nil {
		return err
	}
	return nil
}

func (u *UndoAction) undoTag() error {
	name := u.Params["name"].(string)
	return u.VC.Tag(name, vcs.TagOptions{
		Delete: true,
	})
}

func (u *UndoAction) undoBumpVersion() error {
	// oldVersion := u.Params["oldVersion"].(string)
	// newVersion := u.Params["newVersion"].(string)
	return nil
}
