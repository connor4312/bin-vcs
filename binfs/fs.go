package main

// Number of bytes in a sequence which must be zero to constitute a chunk
// of data, in a rolling checksum. 12 bytes is, on average, 4 MB chunks.
const CHUNK_SIZE = 12

const ADLER_MOD = 65521
