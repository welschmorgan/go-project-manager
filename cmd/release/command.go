package release

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/welschmorgan/go-project-manager/cmd/config"
	"github.com/welschmorgan/go-project-manager/models"
	"github.com/welschmorgan/go-project-manager/ui"
	"github.com/welschmorgan/go-project-manager/vcs"
)

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

func stashModifications(p *models.Project, v vcs.VersionControlSoftware) error {
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

func release(p *models.Project) (err error) {
	if err = os.Chdir(p.Path); err != nil {
		return err
	}
	var vc vcs.VersionControlSoftware
	if vc, err = vcs.Open(p.Path); err != nil {
		return err
	}
	if err = stashModifications(p, vc); err != nil {
		return err
	}
	if err = checkoutAndPullBranch(p, vc, "master"); err != nil {
		return err
	}
	if err = checkoutAndPullBranch(p, vc, "develop"); err != nil {
		return err
	}
	if err = pullTags(p, vc); err != nil {
		return err
	}
	if err = releaseStart(p, vc); err != nil {
		return err
	}
	return nil
}

func checkoutBranch(p *models.Project, v vcs.VersionControlSoftware, branch string) error {
	if err := v.Checkout(branch, vcs.CheckoutOptions{CreateBranch: false}); err != nil {
		return err
	}
	return nil
}

func pullBranch(p *models.Project, v vcs.VersionControlSoftware) error {
	if err := v.Pull(vcs.PullOptions{All: false, Tags: false, Force: false}); err != nil {
		return err
	}
	return nil
}

func checkoutAndPullBranch(p *models.Project, v vcs.VersionControlSoftware, branch string) error {
	if err := checkoutBranch(p, v, branch); err != nil {
		return err
	}
	if err := pullBranch(p, v); err != nil {
		return err
	}
	return nil
}

func pullTags(p *models.Project, v vcs.VersionControlSoftware) error {
	if err := v.Pull(vcs.PullOptions{All: false, Tags: true, Force: true}); err != nil {
		return err
	}
	return nil
}
