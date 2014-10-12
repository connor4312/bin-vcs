package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

var basePath string

func fixDir(add string) {
	if len(basePath) == 0 {
		basePath, _ = filepath.Abs("")
	}
	os.Chdir(filepath.Join(basePath, "test/fixture/", add))
	clearVcs()
}

func TestStagesBasic(t *testing.T) {
	fixDir("flat_basic")

	o := makeRelative(stageFiles(".", false, false, false))
	assert := assert.New(t)

	assert.Equal("a.txt", o[0], "")
	assert.Equal("b.txt", o[1], "")
	assert.Equal(2, len(o), "")
}

func TestStagesNested(t *testing.T) {
	fixDir("nested_basic")

	o := makeRelative(stageFiles(".", false, false, false))
	assert := assert.New(t)

	assert.Equal("a.txt", o[0], "")
	assert.Equal("b.txt", o[1], "")
	assert.Equal("c/d.txt", o[2], "")
	assert.Equal(3, len(o), "")
}

func TestDoesNotStageParents(t *testing.T) {
	fixDir("nested_basic")

	o := makeRelative(stageFiles("c", false, false, false))
	assert := assert.New(t)

	assert.Equal("c/d.txt", o[0], "")
	assert.Equal(1, len(o), "")
}

func TestStagesNestedIgnores(t *testing.T) {
	fixDir("nested_ignores")
	assert := assert.New(t)

	o := makeRelative(stageFiles(".", false, false, false))
	assert.Equal(".binignore", o[0], "staged ignore")
	assert.Equal("a.txt", o[1], "staged text file a")
	assert.Equal("b.txt", o[2], "staged text file b")
	assert.Equal("d/.binignore", o[3], "staged nexted ignore")
	assert.Equal("d/b.avi", o[4], "staged nested avi but not text file")
	assert.Equal("e/c.txt", o[5], "does not screw up ignoring nested")
	assert.Equal(6, len(o), "does not have extraneous files")
}

func TestRespectsParentIgnores(t *testing.T) {
	fixDir("test2")
	assert := assert.New(t)

	o := makeRelative(stageFiles("e", false, false, false))
	assert.Equal(0, len(o), "respects gitignore of parent")
}
