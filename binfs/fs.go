package main

// Number of bytes in a sequence which must be zero to constitute a chunk
// of data, in a rolling checksum. 12 bytes is, on average, 4 MB chunks.
const CHUNK_SIZE = 2

// The number of bytes to read at once from a file. A value of "1" will give
// finer-grained chunks, but it's slower.
const BYTES_PER_CHUNK_ITEM = 8

const ADLER_MOD = 65521
const ADLER_MAX = 5552
