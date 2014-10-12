package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type FileList []string
type DirList []string

// Makes a list of paths relative to the vcs.
func makeRelative(paths []string) []string {
	for index, path := range paths {
		paths[index], _ = filepath.Rel(getVcs().Base, path)
	}

	return paths
}

// Helper function that allows us to take variadic paths, or
// fall back to the cwd.
func initPath(path []string) string {
	if len(path) == 0 {
		out, _ := filepath.Abs("")
		return out
	}

	return filepath.Join(path...)
}

// Simple function to glob from the current directory.
func glob(pattern string, path ...string) (matches FileList) {
	// Search in the cwd unless we were passed a path.
	base := initPath(path)

	// Run the globbed search
	matches, err := filepath.Glob(filepath.Join(base, pattern))
	if err != nil {
		log.Panicf("%s", err)
	}

	return matches
}

// Returns a list of all parent directories of the path. Recursion is fun.
// So are variadic arguments! Takes a split path.
func getAllParents(path string) DirList {
	listing := [][]string{strings.Split(path, string(os.PathSeparator))}

	for len(listing[len(listing)-1]) > 1 {
		last := listing[len(listing)-1]
		listing = append(listing, last[0:len(last)-1])
	}

	output := DirList{}
	for _, item := range listing {
		output = append(output, "/"+filepath.Join(item...))
	}

	return output
}
