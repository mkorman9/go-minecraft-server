package types

import (
	"math/rand"
)

func GetVarIntSize(value int) int {
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

func GetRandomUUID() UUID {
	return UUID{
		Upper: rand.Int63(),
		Lower: rand.Int63(),
	}
}
