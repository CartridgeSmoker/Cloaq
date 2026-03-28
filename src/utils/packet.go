package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
)

type Packet struct {
	Data    []byte
	Version uint8
}

func Encapsulate(data []byte, key []byte) ([]byte, error) {
	size := len(data)
	const version byte = 1
	const msgType byte = 2

	if size > 65535 {
		return nil, fmt.Errorf("payload too large: %d bytes", size)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	encryptedData := gcm.Seal(nil, nonce, data, nil)
	finalSize := len(encryptedData) + len(nonce)

	buf := make([]byte, 4+finalSize)
	buf[0] = version
	buf[1] = msgType
	binary.BigEndian.PutUint16(buf[2:4], uint16(finalSize))

	copy(buf[4:4+len(nonce)], nonce)
	copy(buf[4+len(nonce):], encryptedData)

	return buf, nil
}
