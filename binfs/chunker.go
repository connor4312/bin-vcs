package main

import (
	"crypto/sha256"
	"io"
)

func chunkIsDone([]byte) bool {
	for i := 0; i < CHUNK_SIZE; i++ {
		if hashed[i] != 0 {
			return false
		}
	}

	return true
}

func chunk(data io.Reader) [][]byte {
	var output [][]byte
	var building []byte

	buf := make([]byte, 32)
	for data.Read(data) > 0 {
		hashed := sha256.Sum256(buf)
		if chunkIsDone(hashed) {
			output = append(output, building)
			building = []byte{}
		} else {
			building = append(building, buf[0])
			buf = buf[1:len(buf)]
		}
	}

	return output
}