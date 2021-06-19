package node

import (
	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/project/accessor"
	"github.com/welschmorgan/go-release-manager/vcs"
)

type NodeScaffolder struct {
	accessor.Scaffolder
}

func NewNodeScaffolder() *NodeScaffolder {
	return &NodeScaffolder{}
}

func (s *NodeScaffolder) Name() string {
	return "Node Scaffolder"
}

func (s *NodeScaffolder) Scaffold(workspace *config.Workspace, project *config.Project, v vcs.VersionControlSoftware, firstInit, vcsEnabled bool) error {
	return nil
}
