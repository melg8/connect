package crypt

import (
	"encoding/binary"
	"fmt"
)

func Checksum(data []byte) (uint32, error) {
	if len(data) < 4 {
		return 0, fmt.Errorf("data is too small")
	}
	if len(data)%4 != 0 {
		return 0, fmt.Errorf("data is not multiple of 4")
	}

	checksum := uint32(0)
	for i := 0; i < len(data); i += 4 {
		checksum ^= binary.BigEndian.Uint32(data[i : i+4])
	}
	return checksum, nil
}
