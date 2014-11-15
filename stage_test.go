package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
)

var basePath string

// Copy file from https://www.socketloop.com/tutorials/golang-copy-directory-including-sub-directories-files
func CopyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}

	}

	return
}

func fixDir(add string) {
	if basePath != "" {
		os.Chdir(basePath)
	}

	target := "test_tmp"

	os.RemoveAll(target)
	err := filepath.Walk(filepath.Join(basePath, "test/fixture"), func(file string, info os.FileInfo, err error) error {
		p := filepath.Join(target, file)
		if err != nil {
			w, _ := filepath.Abs("")
			log.Panicf("%s", w)
			log.Panicf("Error searching for %s: %s", target, err)
		}
		if info.IsDir() {
			return os.MkdirAll(p, info.Mode())
		} else {
			return CopyFile(file, p)
		}
	})
	if err != nil {
		log.Panicf("%s", err)
	}

	basePath, _ := filepath.Abs("")
	os.Chdir(filepath.Join(basePath, target, "test/fixture", add))
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
