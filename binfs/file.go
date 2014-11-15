package main

import (
	"log"
	"os"
)

type File struct {
	Name       string
	Path       string
	RevisionOf *File
}

func (f *File) Chunk() {
	data, err := os.Open(f.Path)
	if err != nil {
		log.Printf("Error opening %s: %s", f.Path, err)
		return
	}

	ch := make(chan []byte)
	go chunkData(data, ch)

	for {
		chunk := <-ch

		if len(chunk) == 0 {
			return
		} else {
			print(chunk)
		}
	}
}
