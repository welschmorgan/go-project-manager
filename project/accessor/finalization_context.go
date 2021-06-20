package accessor

import (
	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/vcs"
)

type FinalizationContext struct {
	UserWantsVCSInit        bool                       // user answered yes to initialize repository
	UserWantsScaffolding    bool                       // user answered yes to generate scaffolding
	UserWantsFolderCreation bool                       // user answered yes to create project folder
	RepositoryInitialized   bool                       // wether repository is already initialized
	InitialCommitExists     bool                       // initial commit found
	DevelopExists           bool                       // development branch exists
	MasterExists            bool                       // production branch exists
	VC                      vcs.VersionControlSoftware // project version control software
	PA                      ProjectAccessor            // project accessor
	Project                 *config.Project            // project infos
	Workspace               *config.Workspace          // project workspace
}

func NewFinalizationContext(vc vcs.VersionControlSoftware, pa ProjectAccessor, project *config.Project, workspace *config.Workspace) *FinalizationContext {
	return &FinalizationContext{
		VC:        vc,
		PA:        pa,
		Project:   project,
		Workspace: workspace,
	}
}
