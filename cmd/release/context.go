package release

import "github.com/welschmorgan/go-release-manager/version"

type State uint

const (
	ReleaseStarted        = 1 << iota
	ReleaseFinished       = 1 << iota
	ReleaseStartStarted   = 1 << iota
	ReleaseStartFinished  = 1 << iota
	ReleaseFinishStarted  = 1 << iota
	ReleaseFinishFinished = 1 << iota
)

var States = map[string]State{
	"started":         ReleaseStarted,
	"finished":        ReleaseFinished,
	"start_started":   ReleaseStartStarted,
	"start_finished":  ReleaseStartFinished,
	"finish_started":  ReleaseFinishStarted,
	"finish_finished": ReleaseFinishFinished,
}

func (s State) String() string {
	r := ""
	for name, value := range States {
		if s&value != 0 {
			if len(r) > 0 {
				r += " | "
			}
			r += name
		}
	}
	return r
}

// The release context
type Context struct {
	startingBranch string          // The branch the repository was on before starting release
	oldBranch      string          // The branch before checking out the current one
	releaseBranch  string          // The release branch
	devBranch      string          // The development branch
	prodBranch     string          // The production branch
	version        version.Version // The version the project is in
	nextVersion    version.Version // The next version the project will be in after release
	hasRemotes     bool            // Wether the repository has remotes or not
	state          State           // The state of the release
}
