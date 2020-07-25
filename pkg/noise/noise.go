package noise

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

const chunkSize = 8

type Noise struct {
	filled bool
	left   *bytes.Buffer
	right  *bytes.Buffer
}

func (noise Noise) Stream(samples [][2]float64) (n int, ok bool) {
	buffer := make([]byte, 8)

	for i := range samples {
		// Read the left bit.
		_, err := noise.left.Read(buffer)

		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}

			break
		}

		samples[i][0] = float64frombytes(buffer)

		// Read the right channel
		_, err = noise.right.Read(buffer)

		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}

			break
		}

		samples[i][1] = float64frombytes(buffer)
	}

	return len(samples), true
}

func (noise Noise) Err() error {
	return nil
}

func float64frombytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func float64bytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func New(data []byte) *Noise {
	var leftBytes, rightBytes []byte

	for i, b := range data {
		if i%2 == 0 {
			leftBytes = append(leftBytes, b)
		} else {
			rightBytes = append(rightBytes, b)
		}
	}

	lenLeft := len(leftBytes)
	lenRight := len(rightBytes)

	modLen := lenLeft % chunkSize
	modRight := lenRight % chunkSize

	if modLen != 0 || modRight != 0 {
		if modLen != 0 {
			for i := 0; i < chunkSize-modLen; i++ {
				leftBytes = append(leftBytes, 255)
			}
		}
		if modRight != 0 {
			for i := 0; i <= chunkSize-modRight; i++ {
				rightBytes = append(rightBytes, 255)
			}
		}
	}

	if lenLeft != lenRight {
		if lenLeft < lenRight {
			leftBytes = append(leftBytes, 255, 255, 255, 255, 255, 255, 255, 255)
		} else {
			rightBytes = append(rightBytes, 255, 255, 255, 255, 255, 255, 255, 255)
		}
	}

	return &Noise{
		left:   bytes.NewBuffer(leftBytes),
		right:  bytes.NewBuffer(rightBytes),
		filled: true,
	}

}
