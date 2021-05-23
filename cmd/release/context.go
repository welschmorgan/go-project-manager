package release

// The release context
type Context struct {
	startingBranch string // The branch the repository was on before starting release
	oldBranch      string // The branch before checking out the current one
	releaseBranch  string // The release branch
	devBranch      string // The development branch
	prodBranch     string // The production branch
	version        string // The version the project is in
	nextVersion    string // The next version the project will be in after release
	hasRemotes     bool   // Wether the repository has remotes or not
}
