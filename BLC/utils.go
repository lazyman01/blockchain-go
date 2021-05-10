package BLC

import (
	"bytes"
	"encoding/binary"
	"log"
)

func IntToHex(d int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, d)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}