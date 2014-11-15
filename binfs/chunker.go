package main

import (
	"io"
)

type adlerDigest struct {
	A          uint32
	B          uint32
	Components []uint32
}

func MakeDigest() adlerDigest {
	var digest adlerDigest
	digest.Components = make([]uint32, 2<<10)

	return digest
}

func (a *adlerDigest) Update(newValue uint32) {
	oldValue := a.Components[0]
	a.Components = append(a.Components[1:len(a.Components)], newValue)

	a.A += newValue - oldValue
	a.B += a.A - (32 * oldValue) - 1

	if a.B > (0xffffffff-255)/2 {
		a.A %= ADLER_MOD
		a.B %= ADLER_MOD
	}
}

func (a *adlerDigest) Digest() uint32 {
	if a.B >= ADLER_MOD {
		a.A %= ADLER_MOD
		a.B %= ADLER_MOD
	}

	return a.B<<16 | a.A
}

func (a *adlerDigest) Sum(build []byte) []byte {
	s := a.Digest()
	build = append(build, byte(s>>24))
	build = append(build, byte(s>>16))
	build = append(build, byte(s>>8))
	build = append(build, byte(s))

	return build
}

func chunkIsDone(hashed []byte) bool {
	for i := 0; i < CHUNK_SIZE; i++ {
		if hashed[i] != 0 {
			return false
		}
	}

	return true
}

func chunkData(data io.Reader, out chan []byte) {
	digest := MakeDigest()
	chunk := []byte{}
	for {
		b := make([]byte, 1)
		n, _ := data.Read(b)
		if n == 0 {
			out <- chunk
			out <- []byte{}
			return
		}

		digest.Update(uint32(b[0]))
		chunk = append(chunk, b[0])

		if chunkIsDone(digest.Sum(nil)) {
			out <- chunk
			chunk = []byte{}
		}
	}
}
