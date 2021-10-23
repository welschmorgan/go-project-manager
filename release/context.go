package release

import (
	"time"

	"github.com/welschmorgan/go-release-manager/project/accessor"
	"github.com/welschmorgan/go-release-manager/version"
)

type State uint

const (
	ReleaseCreated        = 0
	ReleaseStarted        = 1 << iota
	ReleaseFinished       = 1 << iota
	ReleaseStartStarted   = 1 << iota
	ReleaseStartFinished  = 1 << iota
	ReleaseFinishStarted  = 1 << iota
	ReleaseFinishFinished = 1 << iota
)

var States = map[string]State{
	"created":         ReleaseCreated,
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

// func (c State) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(c.String())
// }

// func (c State) UnmarshalJSON(data []byte) error {
// 	s := ""
// 	if err := json.Unmarshal(data, &s); err != nil {
// 		return err
// 	}
// 	println(s)
// 	c = ReleaseCreated
// 	for _, part := range strings.Split(s, "|") {
// 		if v, ok := States[strings.TrimSpace(part)]; ok {
// 			c = c | v
// 		}
// 	}
// 	return nil
// }

// The release context
type Context struct {
	StartingBranch string                   `yaml:"starting_branch,omitempty" json:"startingBranch,omitempty"` // The branch the repository was on before starting release
	OldBranch      string                   `yaml:"old_branch,omitempty" json:"oldBranch,omitempty"`           // The branch before checking out the current one
	ReleaseBranch  string                   `yaml:"release_branch,omitempty" json:"releaseBranch,omitempty"`   // The release branch
	DevBranch      string                   `yaml:"dev_branch,omitempty" json:"devBranch,omitempty"`           // The development branch
	ProdBranch     string                   `yaml:"prod_branch,omitempty" json:"prodBranch,omitempty"`         // The production branch
	Date           time.Time                `yaml:"date,omitempty" json:"date,omitempty"`
	Version        version.Version          `yaml:"version,omitempty" json:"version,omitempty"`          // The version the project is in
	NextVersion    version.Version          `yaml:"next_version,omitempty" json:"nextVersion,omitempty"` // The next version the project will be in after release
	HasRemotes     bool                     `yaml:"has_remotes,omitempty" json:"hasRemotes,omitempty"`   // Wether the repository has remotes or not
	State          State                    `yaml:"state,omitempty" json:"state,omitempty"`              // The state of the release
	Accessor       accessor.ProjectAccessor `yaml:"-" json:"-"`                                          // The project accessor
}
