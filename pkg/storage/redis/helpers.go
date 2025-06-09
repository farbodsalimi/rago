package redisutils

import (
	"encoding/binary"
	"math"
)

func floatsToBytes(fs []float64) ([]byte, error) {
	buf := make([]byte, len(fs)*4)

	for i, f := range fs {
		float32Val := float32(f)
		u := math.Float32bits(float32Val)
		binary.NativeEndian.PutUint32(buf[i*4:], u)
	}

	return buf, nil
}
