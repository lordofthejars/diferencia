package plain

// NoiseOperation struct
type NoiseOperation struct {
	index int
	noise bool
}

// ContainsNoise method
func (nd NoiseOperation) ContainsNoise() bool {
	return nd.noise
}

// Detect Noise between documents
func (nd *NoiseOperation) Detect(primary, secondary []byte) {

	for index := range primary {
		if index >= len(secondary) {
			nd.index = index
			nd.noise = true
			break
		}

		if primary[index] != secondary[index] {
			nd.index = index
			nd.noise = true
			break
		}

	}

	if !nd.ContainsNoise() {
		if len(secondary) > len(primary) {
			nd.index = len(primary) - 1
			nd.noise = true
		}
	}

}

// Remove noise from primary and candidate documents
func (nd *NoiseOperation) Remove(primary, candidate []byte) ([]byte, []byte) {

	// It is different from first char, to avoid false substring, we just return the original one
	if nd.ContainsNoise() {
		if nd.index == 0 {
			return primary, candidate
		} else {
			return primary[0:nd.index], candidate[0:nd.index]
		}
	}

	return primary, candidate

}
