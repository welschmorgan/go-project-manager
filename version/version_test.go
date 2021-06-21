package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionParsingSimple(t *testing.T) {
	v := Parse("1.0.0")
	assert.Equal(t, 3, v.NumNonEmptyParts())
	assert.Equal(t, []string{"1", "0", "0"}, v.NonEmptyParts())
}

func TestVersionParsingSemverPreRelease(t *testing.T) {
	v := Parse("1.0.0-rc0")
	assert.Equal(t, 4, v.NumNonEmptyParts())
	assert.Equal(t, []string{"1", "0", "0", "rc0"}, v.NonEmptyParts())
}
func TestVersionParsingSemverBuildMetaTag(t *testing.T) {
	hash := "sadfeed44"
	v := Parse("1.0.0+" + hash)
	assert.Equal(t, 4, v.NumNonEmptyParts())
	assert.Equal(t, []string{"1", "0", "0", hash}, v.NonEmptyParts())
}

func TestVersionIncrementPreRelease(t *testing.T) {
	v := New(0, 0, 0, 1, "rc1")
	if err := v.Increment(PreRelease, 1); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "0.0.0.1-rc2", v.String())
	}
}

func TestVersionIncrementBuildMetaTag(t *testing.T) {
	v := New(0, 0, 0, 1, 0, "", "234234sdf")
	if err := v.Increment(BuildMetaTag, 1); err == nil {
		t.Fatalf("PreRelease part of version should not be incrementable")
	}
}

func TestVersionIncrementRevision(t *testing.T) {
	v := New(0, 0, 0, 1)
	if err := v.Increment(Revision, 1); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "0.0.0.2", v.String())
	}
}

func TestVersionIncrementBuild(t *testing.T) {
	v := New(0, 0, 1, 1)
	if err := v.Increment(Build, 1); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "0.0.2.0", v.String())
	}
}

func TestVersionIncrementMinor(t *testing.T) {
	v := New(0, 1, 1, 1)
	if err := v.Increment(Minor, 1); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "0.2.0.0", v.String())
	}
}

func TestVersionIncrementMajor(t *testing.T) {
	v := New(1, 1, 1, 1)
	if err := v.Increment(Major, 1); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "2.0.0.0", v.String())
	}
}

func TestVersionIncrementSuffix(t *testing.T) {
	v := New(1, 0, 0, 0, "rc0")
	if err := v.Increment(4, 1); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "1.0.0.0-rc1", v.String())
	}
}

func TestVersionDecrementRevision(t *testing.T) {
	v := New(0, 0, 0, 2)
	if err := v.Decrement(3, 1); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "0.0.0.1", v.String())
	}
}

func TestVersionDecrementBuild(t *testing.T) {
	v := New(0, 0, 2, 1)
	if err := v.Decrement(2, 1); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "0.0.1.0", v.String())
	}
}

func TestVersionDecrementMinor(t *testing.T) {
	v := New(0, 2, 1, 1)
	if err := v.Decrement(1, 1); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "0.1.0.0", v.String())
	}
}

func TestVersionDecrementMajor(t *testing.T) {
	v := New(2, 1, 1, 1)
	if err := v.Decrement(0, 1); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "1.0.0.0", v.String())
	}
}

func TestVersionDecrementSuffix(t *testing.T) {
	v := New(1, 0, 0, 0, "rc1")
	if err := v.Decrement(4, 1); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "1.0.0.0-rc0", v.String())
	}
}

func TestVersionDecrementPreRelease(t *testing.T) {
	v := New(0, 0, 0, 1, "rc2")
	if err := v.Decrement(PreRelease, 1); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "0.0.0.1-rc1", v.String())
	}
}

func TestVersionDecrementBuildMetaTag(t *testing.T) {
	v := New(0, 0, 0, 1, 0, "", "234234sdf")
	if err := v.Decrement(BuildMetaTag, 1); err == nil {
		t.Fatalf("PreRelease part of version should not be decrementable")
	}
}
