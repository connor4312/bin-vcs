package main

import (
	"log"
	"os"
	"path/filepath"
)

type VcsSystem struct {
	Dir  string // The directory of the vcs system, including the meta-directory (.vcs by default)
	Base string // The base of the vcs system, excluding the meta directory.
}

var vcsSystem *VcsSystem

// Creates the common vcsSystem struct
func initVcs() {
	// This is how you make a VCS.
	vcsSystem = &VcsSystem{}
	// For every parent dir of the path
	path, _ := filepath.Abs("")
	for _, base := range getAllParents(path) {
		// Join the current path and the VCS_DIR to look for
		search := filepath.Join(base, VCS_DIR)

		if stat, err := os.Stat(search); err == nil && stat.IsDir() {
			vcsSystem.Dir = search
			vcsSystem.Base = base
			return
		}
	}

	log.Panicf("Cannot find vcs setup in current directory or any parent directores.")
}

func clearVcs() {
	vcsSystem = nil
}

// Tries to get the path to the vcs base directory, starting at the cwd.
// Returns a path string if found, or nil if not.
func getVcs() *VcsSystem {
	if vcsSystem == nil {
		initVcs()
	}

	return vcsSystem
}
