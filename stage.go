package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type StagedData struct {
	Path       string   // Base path we are staging from
	Force      bool     // Whether we are force-adding files
	Files      FileList // All files that can be staged
	Ignored    FileList // All ignored files
	Interested FileList // The set of files excluding ignores
	Stages     FileList // The set of interested files exluding unmodified
	Vcs        *VcsSystem
}

// Gathers all files that can be staged.
func (s *StagedData) GatherFiles() {
	// If the path doesn't exist, we have nothing to do...
	stat, err := os.Stat(s.Path)
	if err != nil {
		LogVerbose("Path %s doesn't exist, nothing to stage.", s.Path)
		s.Files = []string{}
		return
	}

	if stat.IsDir() {
		// If the path is a directory, walk over it and add the files inside.
		filepath.Walk(s.Path, func(file string, info os.FileInfo, err error) error {
			if !info.IsDir() && !strings.HasSuffix(s.Path, VCS_DIR) {
				LogVerbose("Discovered %s in the tree.", file)
				s.Files = append(s.Files, file)
			}
			return nil
		})
	} else {
		// Otherwise, just send back the path itself
		s.Files = append(s.Files, s.Path)
		LogVerbose("Discovered %s in the tree.", s.Path)
	}
}

func (s *StagedData) GatherIgnores() {
	// If we're forcing adding files, no need to try and ignore some!
	if s.Force {
		LogVerbose("Force adding files, not loading ignores...")
		return
	}
	// Find all vcs ignore files in this path. Some will be available for
	// staging, so we can just run through the files we picked up before
	// and check them out.
	files := FileList{}
	for _, file := range s.Files {
		if strings.HasSuffix(file, VCS_IGNORE_FILE) {
			files = append(files, file)
		}
	}
	// We should also get the ignores of any parents, father up the
	// directory tree. Loop through parent directories...
	to_length := len(getVcs().Base)
	for _, parent := range getAllParents(s.Path) {
		// If the length of the parent path is larger than the vcs install
		// path, it must be nested!
		if len(parent) >= to_length {
			// Add the ignore file, if one exists.
			file := filepath.Join(parent, VCS_IGNORE_FILE)
			if _, err := os.Stat(file); err == nil {
				files = append(files, file)
			}
		}
	}

	LogVerbose("Loaded ignores: %s", files)

	// If we didn't get any files, chill out, we're good!
	if len(files) == 0 {
		return
	}

	// Make an output channel and spin off routines to go search for
	// ignored files we don't want.
	ch := make(chan FileList)
	for _, file := range files {
		go loadIgnoresFor(file, ch)
	}

	// Gather all the data as it's pushed back down the channel.
	for i := len(files); i > 0; i-- {
		s.Ignored = append(s.Ignored, <-ch...)
	}
}

// Get the files we're interested in -- that is, all files minus the
// ignored files.
func (s *StagedData) GatherInterested() {
	for _, file := range s.Files {
		if strings.HasSuffix(file, VCS_IGNORE_FILE) || !inSlice(file, s.Ignored) {
			s.Interested = append(s.Interested, file)
		}
	}
}

// Loads all globbed ignored files in the directory...
func loadIgnoresFor(path string, out chan FileList) {
	// Open the ignore file, or fail if we can't.
	file, err := os.Open(path)
	if err != nil {
		log.Panicf("Error opening file %s: %s", path, err)
	}
	defer file.Close()

	// Look throught every line and add the matches
	scanner := bufio.NewScanner(file)
	output := FileList{}
	for scanner.Scan() {
		g := filepath.Join(filepath.Dir(path), scanner.Text())
		files, _ := filepath.Glob(g)
		output = append(output, files...)
	}

	out <- output
}

func stageFiles(path string, dry bool, update bool, force bool) FileList {
	// Fix the absolute path
	path, _ = filepath.Abs(path)

	stage := &StagedData{Vcs: getVcs(), Path: path, Force: force}
	stage.GatherFiles()
	stage.GatherIgnores()
	stage.GatherInterested()

	return stage.Interested
}

func Stage() {
	stageFiles(*cmdAddGlob, *cmdAddDry, *cmdAddUpdate, *cmdAddForce)
}
