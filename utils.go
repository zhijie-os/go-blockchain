package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	// big endian
	err := binary.Write(buff, binary.BigEndian, num)

	// if error, log panic
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
