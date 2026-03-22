package utils

import (
	"encoding/binary"
	"fmt"
	"log"
	"runtime/debug"
)

type Packet struct {
	Data    []byte
	Version uint8
}

// Encapsulate adds a 4-byte header: [version][type][len_high][len_low]
func Encapsulate(version uint8, msgType uint8, data []byte) ([]byte, error) {
	size := len(data)

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

// SafeRuntime prevents goroutine panics from crashing the app
func SafeRuntime(name string, f func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[panic] %s recovered: %v", name, r)
				log.Printf("[stack] %s trace: %s", name, string(debug.Stack()))
			}
		}()
		f()
	}()
}
