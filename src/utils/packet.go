package utils

import (
	"encoding/binary"
	"fmt"
)

type Packet struct {
	Data    []byte
	Version uint8
}

// Encapsulate adds a 4-byte header: [version][type][len_high][len_low]
func Encapsulate(data []byte) ([]byte, error) {
	size := len(data)
	var version byte
	var msgType byte

	if size > 65535 {
		return nil, fmt.Errorf("payload too large: %d bytes", size)
	}

	buf := make([]byte, 4+size)

	buf[0] = version
	buf[1] = msgType

	binary.BigEndian.PutUint16(buf[2:4], uint16(size))

	copy(buf[4:], data)

	return buf, nil
}
