package version

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
)

type Version struct {
	numParts uint8
	suffix   []byte
	parts    [][]byte
}

var Zero = NewVersion(0, 0, 0)
var FirstMajor = NewVersion(1, 0, 0)
var FirstMinor = NewVersion(0, 1, 0)
var FirstPatch = NewVersion(0, 0, 1)
var FirstBuild = NewVersion(0, 0, 0, 1)

func NewVersion(parts ...int) *Version {
	v := &Version{
		numParts: uint8(len(parts)),
		suffix:   []byte{},
		parts:    make([][]byte, len(parts)),
	}
	for i, p := range parts {
		v.parts[i] = strconv.AppendInt(v.parts[i], int64(p), 32)
	}
	return v
}

func Parse(s string) *Version {
	v := &Version{
		numParts: 0,
		suffix:   []byte{},
		parts:    [][]byte{},
	}
	suffix := bytes.SplitN([]byte(s), []byte("-"), 2)
	if len(suffix) > 1 {
		v.suffix = suffix[1]
	}
	v.parts = bytes.SplitN(suffix[0], []byte("."), 4)
	v.numParts = uint8(len(v.parts))
	return v
}

func (v *Version) GetString(i uint8) (string, error) {
	if b, err := v.GetBytes(i); err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}

func (v *Version) GetBytes(i uint8) ([]byte, error) {
	if i > uint8(len(v.parts)) {
		return nil, fmt.Errorf("%d/%d: index out of bounds", i, len(v.parts))
	}
	if i == uint8(len(v.parts)) {
		return v.suffix, nil
	}
	return v.parts[i], nil
}

func (v *Version) GetInt(i uint8) (int, error) {
	if val, err := v.GetString(i); err != nil {
		return -1, err
	} else if ret, err := strconv.ParseInt(val, 10, 32); err != nil {
		return -1, err
	} else {
		return int(ret), nil
	}
}

func (v *Version) SetString(i uint8, val string) error {
	return v.SetBytes(i, []byte(val))
}

func (v *Version) SetBytes(i uint8, val []byte) error {
	if i > uint8(len(v.parts)) {
		return fmt.Errorf("%d/%d: index out of bounds", i, len(v.parts))
	}
	if i == uint8(len(v.parts)) {
		v.suffix = val
		return nil
	}
	v.parts[i] = val
	return nil
}

func (v *Version) SetInt(i uint8, val int) error {
	return v.SetString(i, fmt.Sprintf("%d", val))
}

func (v *Version) Increment(i uint8, step int) error {
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
		if ival, err := strconv.ParseInt(val, 10, 32); err != nil {
			return err
		} else if len(prefix) > 0 {
			if err = v.SetString(i, fmt.Sprintf("%s%d", prefix, int(ival)+step)); err != nil {
				return err
			}
			for j := i + 1; j < uint8(len(v.parts)); j++ {
				if err = v.SetInt(j, 0); err != nil {
					return err
				}
			}
			return nil
		} else {
			if err = v.SetInt(i, int(ival)+step); err != nil {
				return err
			}
			for j := i + 1; j < uint8(len(v.parts)); j++ {
				if err = v.SetInt(j, 0); err != nil {
					return err
				}
			}
			return nil
		}
	}
}

func (v *Version) Decrement(i uint8, step int) error {
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
			for j := i + 1; j < uint8(len(v.parts)); j++ {
				if err = v.SetInt(j, 0); err != nil {
					return err
				}
			}
			return nil
		} else {
			if err = v.SetInt(i, int(ival)-step); err != nil {
				return err
			}
			for j := i + 1; j < uint8(len(v.parts)); j++ {
				if err = v.SetInt(j, 0); err != nil {
					return err
				}
			}
			return nil
		}
	}
}

func (v *Version) String() string {
	s := ""
	for _, p := range v.parts {
		if len(s) > 0 {
			s += "."
		}
		s += string(p)
	}
	if len(v.suffix) > 0 {
		s += "-" + string(v.suffix)
	}
	return s
}
