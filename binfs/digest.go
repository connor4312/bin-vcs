package main

// Represents an adler digest. This is a lot like Go's default adler32
// function, except we support rolling hashes - don't have to recompute
// every element every time, only new elements.
type adlerDigest struct {
	Components []uint32
	S1         uint32
	S2         uint32
}

func MakeDigest() adlerDigest {
	var digest adlerDigest
	digest.S1, digest.S2 = 0, 0
	digest.Components = []uint32{}

	return digest
}

// "rolls" the adler digest with the new value, popping the first item
// off of the hash if we have a full 32 bytes.
func (a *adlerDigest) Update(newValue uint32) {
	if len(a.Components) < 32 {
		a.Components = append(a.Components, newValue)
		a.S1 += newValue
		a.S2 += a.S1
	} else {
		oldValue := a.Components[0]
		a.Components = append(a.Components[1:len(a.Components)], newValue)

		a.S1 += newValue - oldValue
		a.S2 += a.S1 - 31*oldValue
	}
}

// Finalizes an adler digest and returns a uint sum.
func (a *adlerDigest) Digest() uint32 {
	return (a.S2 << 16) | (a.S1 & 0xffff)
}

// Returns the adler digest sum as a byte array.
func (a *adlerDigest) Sum(build []byte) []byte {
	s := a.Digest()
	return append(build, byte(s>>24), byte(s>>16), byte(s>>8), byte(s))
}
