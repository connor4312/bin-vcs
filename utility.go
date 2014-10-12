package main

// Checks to see if the string exists in a slice of strings.
func inSlice(needle string, haystack []string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}

	return false
}
