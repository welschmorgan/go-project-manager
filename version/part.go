package version

import "errors"

type VersionPart uint8

const (
	Major VersionPart = iota
	Minor
	Build
	Revision
	PreRelease
	BuildMetaTag
)

type versionPartData struct {
	id        VersionPart
	name      string
	separator string
}

var versionParts = []versionPartData{
	{Major, "major", ""},
	{Minor, "minor", "."},
	{Build, "build", "."},
	{Revision, "revision", "."},
	{PreRelease, "preRelease", "-"},
	{BuildMetaTag, "buildMetaTag", "+"},
}

func getData(vp VersionPart) *versionPartData {
	for _, p := range versionParts {
		if p.id == vp {
			return &p
		}
	}
	return nil
}

func ParsePart(s string) (VersionPart, error) {
	for _, v := range versionParts {
		if s == v.name {
			return v.id, nil
			break
		}
	}
	return Major, errors.New("failed to parse part from '" + s + "'")
}
func (p VersionPart) String() string {
	return p.Name()
}

func (p VersionPart) Id() uint8 {
	return uint8(versionParts[p].id)
}

func (p VersionPart) Name() string {
	return versionParts[p].name
}

func (p VersionPart) Separator() string {
	return versionParts[p].separator
}
