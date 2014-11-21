package main

import (
	"testing"
)

func TestChunking(t *testing.T) {
	file := File{Path: "fixture/file.png"}
	file.Chunk()
}
