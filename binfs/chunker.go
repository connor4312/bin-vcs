package main

import (
	"io"
)

// Checks to see if the chunk with the hashed bytes is "done" - if its first
// `CHUNK_SIZE` bits are all zeroes.
func chunkIsDone(hashed []byte) bool {
	for i := 0; i < CHUNK_SIZE; i++ {
		if hashed[i] != 0 {
			return false
		}
	}

	return true
}

// Chunks out data from the reader, emitting byte arrays down the channel.
// After the reader's entire content has been processed, a zero-length
// byte array is sent down to signal completion.
func chunkData(data io.Reader, out chan []byte) {
	digest := MakeDigest()
	chunk := []byte{}

	for {
		// Read up to four bytes off of the io.Reader
		bytes := make([]byte, BYTES_PER_CHUNK_ITEM)
		n, _ := data.Read(bytes)
		// If we didn't read anything, we reached the end of input. Send back
		// the final chunk then an empty byte array to signal completion.
		if n == 0 {
			out <- chunk
			out <- []byte{}
			return
		}

		// Otherwise go ahead and roll the bytes onto the digest and
		// add them to the chunk.
		for _, k := range bytes {
			digest.Update(uint32(k))
		}
		chunk = append(chunk, bytes...)

		// If the chunk is ready to send back, emit it down the channel then
		// clear it.
		if chunkIsDone(digest.Sum(nil)) {
			out <- chunk
			chunk = []byte{}
		}
	}
}
