package packets

func getVarIntSize(value int) int {
	size := 0

	for {
		if (value & ^SegmentBits) == 0 {
			size++
			break
		}

		size++
		value >>= 7
	}

	return size
}
