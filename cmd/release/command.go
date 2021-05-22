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

type Context struct {
	releaseBranch string
	devBranch     string
	prodBranch    string
	version       string
}

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
	ctx := Context{
		releaseBranch: strings.ReplaceAll(config.Get().BranchNames["release"], "$VERSION", version),
		devBranch:     config.Get().BranchNames["development"],
		prodBranch:    config.Get().BranchNames["production"],
		version:       version,
	}

	// stash modifications
	if err = stashModifications(p, vc, &ctx); err != nil {
		return err
	}

	// checkout development and production branches
	if err = updateRepository(p, vc, &ctx); err != nil {
		return err
	}

	// start release
	if err = releaseStart(p, vc, &ctx); err != nil {
		return err
	}

	// wait for user to manually edit release

	// TODO insert code

	// finish release
	if err = releaseFinish(p, vc, &ctx); err != nil {
		return err
	}
	if err = bumpVersion(p, vc, &ctx); err != nil {
		return err
	}
	return nil
}

func abortRelease(p *config.Project, v vcs.VersionControlSoftware) error {
	println("aborting release ...")
	return nil
}

func updateRepository(p *config.Project, v vcs.VersionControlSoftware, ctx *Context) error {
	var err error
	if err = checkoutAndPullBranch(p, v, ctx.prodBranch, ctx); err != nil {
		return err
	}
	if err = checkoutAndPullBranch(p, v, ctx.devBranch, ctx); err != nil {
		return err
	}
	if err = pullTags(p, v, ctx); err != nil {
		return err
	}
	return nil
}

func stashModifications(p *config.Project, v vcs.VersionControlSoftware, ctx *Context) error {
	if status, err := v.Status(vcs.StatusOptions{Short: true}); err != nil {
		return err
	} else if len(status) != 0 {
		fmt.Printf("There is work in progress:\n%v\n", status)
		if ok, _ := ui.AskYN("Do you want to stash it"); ok {
			if out, err := v.Stash(vcs.StashOptions{
				IncludeUntracked: true,
			}); err != nil {
				return err
			} else {
				fmt.Printf("%v\n", out)
			}
		}
	}
	return nil
}

func checkoutBranch(p *config.Project, v vcs.VersionControlSoftware, branch string, ctx *Context) error {
	if err := v.Checkout(branch, vcs.CheckoutOptions{CreateBranch: false}); err != nil {
		return err
	}
	return nil
}

func pullBranch(p *config.Project, v vcs.VersionControlSoftware, ctx *Context) error {
	if err := v.Pull(vcs.PullOptions{All: false, ListTags: false, Force: false}); err != nil {
		return err
	}
	return nil
}

func checkoutAndPullBranch(p *config.Project, v vcs.VersionControlSoftware, branch string, ctx *Context) error {
	if err := checkoutBranch(p, v, branch, ctx); err != nil {
		return err
	}
	if err := pullBranch(p, v, ctx); err != nil {
		return err
	}
	return nil
}

func pullTags(p *config.Project, v vcs.VersionControlSoftware, ctx *Context) error {
	if err := v.Pull(vcs.PullOptions{All: false, ListTags: true, Force: true}); err != nil {
		return err
	}
	return nil
}

func releaseStart(p *config.Project, v vcs.VersionControlSoftware, ctx *Context) error {
	if err := v.Checkout(ctx.releaseBranch, vcs.CheckoutOptions{
		StartingPoint: ctx.devBranch,
		CreateBranch:  true,
	}); err != nil {
		return err
	}
	return nil
}

func releaseFinish(p *config.Project, v vcs.VersionControlSoftware, ctx *Context) error {
	if err := v.Merge(ctx.releaseBranch, ctx.prodBranch, vcs.MergeOptions{NoFastForward: true}); err != nil {
		return err
	}
	if err := v.Tag(ctx.version, vcs.TagOptions{Annotated: true, Message: fmt.Sprintf("Release %s: %s", ctx.version, "TODO")}); err != nil {
		return err
	}
	return nil
}

func bumpVersion(p *config.Project, v vcs.VersionControlSoftware, ctx *Context) error {
	return nil
}
