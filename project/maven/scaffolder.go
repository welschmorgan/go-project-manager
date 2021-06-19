package maven

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/welschmorgan/go-release-manager/fs"
	"github.com/welschmorgan/go-release-manager/project/accessor"
	"github.com/welschmorgan/go-release-manager/ui"
	"github.com/welschmorgan/go-release-manager/vcs"
)

type MavenScaffolder struct {
	accessor.Scaffolder
}

func NewMavenScaffolder() *MavenScaffolder {
	return &MavenScaffolder{}
}

func (s *MavenScaffolder) Name() string {
	return "Maven Archetype Scaffolder"
}

func (s *MavenScaffolder) SanitizeArtifactId(id string) string {
	rule := regexp.MustCompile(`\W+`)
	ret := rule.ReplaceAllString(id, "-")
	return ret
}

func (s *MavenScaffolder) SanitizeGroupId(id string) string {
	rule := regexp.MustCompile(`\s+`)
	ret := rule.ReplaceAllString(id, "-")
	return ret
}

func (s *MavenScaffolder) Scaffold(ctx *accessor.FinalizationContext) error {
	var artifactId string
	var groupId string
	var version string
	if ok, err := ctx.PA.Detect(ctx.Project.Path); err != nil || !ok {
		if ctx.UserWantsScaffolding, _ = ui.AskYN(ctx.PA.AccessorName() + " have no " + ctx.PA.DescriptionFile() + " file, create basic project scaffolding"); ctx.UserWantsScaffolding {
			fmt.Printf("Initializing %s project...\n", ctx.PA.AccessorName())
			println("------[ Maven POM ]-------")
			var ans string
			var err error
			// if ans, err = ui.Ask("\tModelVersion", DefaultPOMModel); err != nil {
			// 	return err
			// } else {
			// 	a.pom.Root.SetModelVersion(ParseModelVersion(ans))
			// }
			if ans, err = ui.Ask("\tArtifactId", ctx.Project.Name); err != nil {
				return err
			} else {
				artifactId = s.SanitizeArtifactId(ans)
			}
			if ans, err = ui.Ask("\tGroupId", "com."); err != nil {
				return err
			} else {
				groupId = s.SanitizeGroupId(ans)
			}
			if ans, err = ui.Ask("\tVersion", DefaultPOMVersion); err != nil {
				return err
			} else {
				version = ans
			}
			archetypeDir := filepath.Join(os.TempDir(), "maven-archetype")
			if _, err = os.Stat(archetypeDir); err == nil || os.IsExist(err) {
				if err = os.RemoveAll(archetypeDir); err != nil {
					return err
				}
			}
			println(archetypeDir)
			println(artifactId)
			_, stdout, _, _ := vcs.RunCommand("mvn", "-v")
			for _, line := range stdout {
				println(line)
			}
			exit, stdout, stderr, err := vcs.RunCommand(
				"mvn", "archetype:generate",
				"-B",
				"-DgroupId="+groupId,
				"-DartifactId="+artifactId,
				"-Dversion="+version,
				"-Durl="+ctx.Project.Url,
				"-DoutputDirectory="+archetypeDir,
				"-Dmaven.compiler.source="+DefaultPOMJavaVersion,
				"-DarchetypeGroupId=org.apache.maven.archetypes",
				"-DarchetypeArtifactId=maven-archetype-quickstart",
				"-DinteractiveMode=false",
			)
			fmt.Printf("done generating archetype, exit code = %d, %d stdout, %d stderr\n", exit, len(stdout), len(stderr))
			for _, line := range stdout {
				println(line)
			}
			vcs.DumpCommandErrors(exit, stderr)
			if err != nil {
				return fmt.Errorf("failed to generate maven project archetype, %s", err.Error())
			}
			println("copydir")
			if err = fs.CopyDir(filepath.Join(archetypeDir, artifactId), ctx.PA.Path(), false); err != nil {
				return err
			}
			if ctx.RepositoryInitialized {
				if err = ctx.VC.Stage(vcs.StageOptions{All: true}); err != nil {
					return err
				}

				commitOpts := vcs.CommitOptions{}
				if !ctx.InitialCommitExists {
					commitOpts.Message = "Initial commit"
					commitOpts.AllowEmpty = true
				} else {
					commitOpts.Message = fmt.Sprintf("Create %s project file", ctx.PA.AccessorName())
				}
				if err = ctx.VC.Commit(commitOpts); err != nil {
					return err
				}
				ctx.InitialCommitExists = true
			}
		}
	}
	return nil
}
