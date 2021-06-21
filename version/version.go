package version

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const DETECTION_REGEX = `^(?P<major>\d+?|)(\.(?P<minor>\d+?)|)(\.(?P<build>\d+?)|)(\.(?P<revision>\d+?)|)(\-(?P<preRelease>[\d\w]+?)|)(\+(?P<buildMetaTag>[\d\w]+?)|)$`

type Version []string

var PreReleasePrefix = "rc"
var Zero = New(0, 0, 0)
var FirstMajor = New(1, 0, 0)
var FirstMinor = New(0, 1, 0)
var FirstPatch = New(0, 0, 1)
var FirstBuild = New(0, 0, 0, 1)

func Clone(o Version) Version {
	data := []string{}
	for _, p := range o {
		data = append(data, p)
	}
	return Version(data)
}

func New(parts ...interface{}) Version {
	v := Version(make([]string, len(versionParts)))
	for i, p := range parts {
		v[i] = fmt.Sprint(p)
	}
	return v
}

func Parse(s string) Version {
	v := Version(make([]string, len(versionParts)))
	format := regexp.MustCompile(DETECTION_REGEX)
	matches := format.FindStringSubmatch(s)
	for _, vp := range versionParts {
		idx := format.SubexpIndex(vp.name)
		val := strings.TrimSpace(matches[idx])
		if len(val) > 0 {
			v[vp.id] = val
		}
	}
	return v
}

func (v Version) IsEmpty(i VersionPart) bool {
	if uint8(i) < uint8(len(v)) {
		return len(v[i]) == 0
	}
	return true
}

func (v Version) NonEmptyParts() []string {
	ret := []string{}
	for _, p := range v {
		if len(p) > 0 {
			ret = append(ret, p)
		}
	}
	return ret
}

func (v Version) Len() int {
	return len(v)
}

func (v Version) NumNonEmptyParts() int {
	return len(v.NonEmptyParts())
}

func (v Version) HasNonEmptyParts() bool {
	return len(v.NonEmptyParts()) != 0
}

func (v Version) GetString(i VersionPart) (string, error) {
	if uint8(i) > uint8(len(v)) {
		return "", fmt.Errorf("%d/%d: index out of bounds", i, len(v))
	}
	return v[i], nil
}

func (v Version) GetBytes(i VersionPart) ([]byte, error) {
	if b, err := v.GetString(i); err != nil {
		return nil, err
	} else {
		return []byte(b), nil
	}
}

func (v Version) GetInt(i VersionPart) (int, error) {
	if val, err := v.GetString(i); err != nil {
		return -1, err
	} else if ret, err := strconv.ParseInt(val, 10, 32); err != nil {
		return -1, err
	} else {
		return int(ret), nil
	}
}

func (v Version) MustGetString(i VersionPart) string {
	if val, err := v.GetString(i); err != nil {
		panic(err)
	} else {
		return val
	}
}

func (v Version) MustGetBytes(i VersionPart) []byte {
	if val, err := v.GetBytes(i); err != nil {
		panic(err)
	} else {
		return val
	}
}

func (v Version) MustGetInt(i VersionPart) int {
	if val, err := v.GetInt(i); err != nil {
		panic(err)
	} else {
		return val
	}
}

func (v Version) SetString(i VersionPart, val string) error {
	if uint8(i) > uint8(len(v)) {
		return fmt.Errorf("%d/%d: index out of bounds", i, len(v))
	}
	v[i] = string(val)
	return nil
}

func (v Version) SetBytes(i VersionPart, val []byte) error {
	return v.SetString(i, string(val))
}

func (v Version) SetInt(i VersionPart, val int) error {
	return v.SetString(i, fmt.Sprintf("%d", val))
}

func (v Version) Increment(i VersionPart, step int) error {
	if val, err := v.GetString(i); err != nil {
		return err
	} else if step <= 0 {
		return fmt.Errorf("invalid increment step %d", step)
	} else {
		rxp := regexp.MustCompile(`(\w*)(\d+)`)
		matches := rxp.FindAllStringSubmatch(val, -1)
		prefix := ""
		if len(matches) > 0 {
			prefix = matches[0][1]
			val = matches[0][2]
		}
		if val == "" {
			val = "0"
		}
		if len(prefix) == 0 && i == PreRelease {
			prefix = PreReleasePrefix
		}
		var ival int64
		if ival, err = strconv.ParseInt(val, 10, 32); err != nil {
			return err
		}
		if err = v.SetString(i, fmt.Sprintf("%s%d", prefix, int(ival)+step)); err != nil {
			return err
		}
		// set next parts to 0
		for j := uint8(i + 1); j < uint8(len(v)); j++ {
			if !v.IsEmpty(VersionPart(j)) {
				if err = v.SetInt(VersionPart(j), 0); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (v Version) Decrement(i VersionPart, step int) error {
	if val, err := v.GetString(i); err != nil {
		return err
	} else if step <= 0 {
		return fmt.Errorf("invalid decrement step %d", step)
	} else {
		rxp := regexp.MustCompile(`(\w*)(\d+)`)
		matches := rxp.FindAllStringSubmatch(val, -1)
		prefix := ""
		if len(matches) > 0 {
			prefix = matches[0][1]
			val = matches[0][2]
		}
		if ival, err := strconv.ParseInt(val, 10, 32); err != nil {
			return err
		} else if len(prefix) > 0 {
			if err = v.SetString(i, fmt.Sprintf("%s%d", prefix, int(ival)-step)); err != nil {
				return err
			}
			for j := uint8(i + 1); j < uint8(len(v)); j++ {
				if err = v.SetInt(VersionPart(j), 0); err != nil {
					return err
				}
			}
			return nil
		} else {
			if err = v.SetInt(i, int(ival)-step); err != nil {
				return err
			}
			for j := uint8(i + 1); j < uint8(len(v)); j++ {
				if err = v.SetInt(VersionPart(j), 0); err != nil {
					return err
				}
			}
			return nil
		}
	}
}

func (v Version) String() string {
	ret := ""
	for _, vp := range versionParts {
		if !v.IsEmpty(vp.id) {
			if len(ret) > 0 {
				ret += vp.separator
			}
			ret += v[uint8(vp.id)]
		}
	}
	return ret
}

func (v Version) Major() (val string, err error) {
	val, err = v.GetString(Major)
	return
}

func (v Version) MajorInt() (val int, err error) {
	val, err = v.GetInt(Major)
	return
}

func (v Version) MajorBytes() (val []byte, err error) {
	val, err = v.GetBytes(Major)
	return
}

func (v Version) Minor() (val string, err error) {
	val, err = v.GetString(Minor)
	return
}

func (v Version) MinorInt() (val int, err error) {
	val, err = v.GetInt(Minor)
	return
}

func (v Version) MinorBytes() (val []byte, err error) {
	val, err = v.GetBytes(Minor)
	return
}

func (v Version) Build() (val string, err error) {
	val, err = v.GetString(Build)
	return
}

func (v Version) BuildInt() (val int, err error) {
	val, err = v.GetInt(Build)
	return
}

func (v Version) BuildBytes() (val []byte, err error) {
	val, err = v.GetBytes(Build)
	return
}

func (v Version) Revision() (val string, err error) {
	val, err = v.GetString(Revision)
	return
}

func (v Version) RevisionInt() (val int, err error) {
	val, err = v.GetInt(Revision)
	return
}

func (v Version) RevisionBytes() (val []byte, err error) {
	val, err = v.GetBytes(Revision)
	return
}

func (v Version) PreRelease() (val string, err error) {
	val, err = v.GetString(PreRelease)
	return
}

func (v Version) PreReleaseInt() (val int, err error) {
	val, err = v.GetInt(PreRelease)
	return
}

func (v Version) PreReleaseBytes() (val []byte, err error) {
	val, err = v.GetBytes(PreRelease)
	return
}

func (v Version) BuildMetaTag() (val string, err error) {
	val, err = v.GetString(BuildMetaTag)
	return
}

func (v Version) BuildMetaTagInt() (val int, err error) {
	val, err = v.GetInt(BuildMetaTag)
	return
}

func (v Version) BuildMetaTagBytes() (val []byte, err error) {
	val, err = v.GetBytes(BuildMetaTag)
	return
}
