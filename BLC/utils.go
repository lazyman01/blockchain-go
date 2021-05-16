package BLC

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
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

func JSONToArray(jsonString string) []string {
	var sArr []string
	if err := json.Unmarshal([]byte(jsonString), &sArr); err != nil {
		log.Panic(err)
	}

	return sArr
}