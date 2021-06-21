package vcs

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeFakeGitRepoWithVCS(t *testing.T, parentDirPrefix string, commits ...string) (g *Git, dir string, err error) {
	// make temp dir
	if dir, err = os.MkdirTemp("", fmt.Sprintf("test-%s-*", parentDirPrefix)); err != nil {
		return
	}
	g = &Git{}
	// t.Logf("Init git in: %s", dir)
	if err = g.Initialize(dir, InitOptions{}); err != nil {
		return
	}
	// t.Logf("-> path is now: %s", g.path)
	for _, commit := range commits {
		if err = g.Commit(CommitOptions{
			AllowEmpty: true,
			Message:    commit,
		}); err != nil {
			return
		}
	}
	return
}

func makeFakeGitRepoWithFolders(t *testing.T, parentDirPrefix string) (dir string, err error) {
	// make temp dir
	if dir, err = os.MkdirTemp("", fmt.Sprintf("test-%s-*", parentDirPrefix)); err != nil {
		return
	}
	// init git repo
	dotGitDir := filepath.Join(dir, ".git")
	if err = os.Mkdir(dotGitDir, 0755); err != nil {
		return
	}
	return
}

func TestGitShouldInitializeEmptyFolder(t *testing.T) {
	// make temp dir
	git, dir, err := makeFakeGitRepoWithVCS(t, "grlm-git-should-initialize-empty-folder")
	if err != nil {
		t.Fatal(err)
	}
	// check path has been assigned
	assert.Equal(t, dir, git.path)
	// check repository ok
	dotGitDir := filepath.Join(git.path, ".git")
	if fi, err := os.Stat(dotGitDir); err != nil {
		t.Fatal(err)
	} else {
		assert.Equal(t, true, fi.IsDir())
	}
}

func TestGitShouldDetectValidFolder(t *testing.T) {
	// make temp dir
	dir, err := makeFakeGitRepoWithFolders(t, "grlm-git-should-detect-valid-folder")
	if err != nil {
		t.Fatal(err)
	}
	// init git repo
	g := &Git{}
	if err = g.Detect(dir); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, dir, g.path)
}

func TestGitShouldOpenValidFolder(t *testing.T) {
	// make temp dir
	g, dir, err := makeFakeGitRepoWithVCS(t, "grlm-git-should-open-valid-folder")
	if err != nil {
		t.Fatal(err)
	}
	// init git repo
	if err = g.Open(dir); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, dir, g.path)
}

func TestGitCommand_CommitAndExtractLogShouldWork(t *testing.T) {
	// make temp dir
	git, dir, err := makeFakeGitRepoWithVCS(t, "grlm-git-command-commit-and-extract-log-should-work", "first commit", "second commit")
	if err != nil {
		t.Fatal(err)
	}
	// init git repo
	if err = git.Open(dir); err != nil {
		t.Fatal(err)
	}
	var log []string
	if log, err = git.ExtractLog(ExtractLogOptions{
		Limit:  2,
		Format: "%s",
	}); err != nil {
		t.Fatal(err)
	}
	assert.Len(t, log, 2)
	assert.Equal(t, "second commit", log[0])
	assert.Equal(t, "first commit", log[1])
}

func TestGitCommand_CheckoutThenGetCurrentBranchWorks(t *testing.T) {
	// make temp dir
	git, dir, err := makeFakeGitRepoWithVCS(t, "grlm-git-command-checkout-then-get-current-branch-works")
	if err != nil {
		t.Fatal(err)
	}
	// init git repo
	if err = git.Open(dir); err != nil {
		t.Fatal(err)
	}
	// create root commit
	if err = git.Commit(CommitOptions{AllowEmpty: true, Message: "initial commit"}); err != nil {
		t.Error(err)
	}

	var branch string
	// checkout develop
	if err = git.Checkout("develop", CheckoutOptions{CreateBranch: true}); err != nil {
		t.Error(err)
	}
	if branch, err = git.CurrentBranch(); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "develop", branch)
	}
	// checkout master
	if err = git.Checkout("master", CheckoutOptions{CreateBranch: false}); err != nil {
		t.Error(err)
	}
	if branch, err = git.CurrentBranch(); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "master", branch)
	}
}

func TestGitCommand_CreatingAndListingTagsShouldWork(t *testing.T) {
	// make temp dir
	git, dir, err := makeFakeGitRepoWithVCS(t, "grlm-git-command-creating-and-listing-tags-should-work")
	if err != nil {
		t.Fatal(err)
	}
	// init git repo
	if err = git.Open(dir); err != nil {
		t.Fatal(err)
	}
	// create root commit
	if err = git.Commit(CommitOptions{AllowEmpty: true, Message: "initial commit"}); err != nil {
		t.Error(err)
	}
	if err = git.Tag("test-tag", TagOptions{Annotated: true, Message: "Tag testing"}); err != nil {
		t.Error(err)
	}
	var tags []string
	if tags, err = git.ListTags(nil); err != nil {
		t.Error(tags)
	}
	assert.Len(t, tags, 1)
	assert.Equal(t, "test-tag", tags[0])

}

func TestGitCommand_CreatingAndListingStashesShouldWork(t *testing.T) {
	// make temp dir
	git, dir, err := makeFakeGitRepoWithVCS(t, "grlm-git-command-creating-and-listing-stashes-should-work")
	if err != nil {
		t.Fatal(err)
	}
	// init git repo
	if err = git.Open(dir); err != nil {
		t.Fatal(err)
	}
	// create root commit
	if err = git.Commit(CommitOptions{AllowEmpty: true, Message: "initial commit"}); err != nil {
		t.Error(err)
	}
	var f *os.File
	if f, err = os.Create(filepath.Join(dir, "hello.txt")); err != nil {
		t.Error(f)
	}
	var stashes []string
	if _, err = git.Stash(StashOptions{Save: true, IncludeUntracked: true, Message: "test hello stash"}); err != nil {
		t.Error(err)
	}
	if stashes, err = git.ListStashes(); err != nil {
		t.Error(err)
	}

	assert.Len(t, stashes, 1)
	assert.Equal(t, "stash@{0}: On master: test hello stash", stashes[0])
}
