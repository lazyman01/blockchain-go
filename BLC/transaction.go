package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type Transaction struct {
	// 1.交易hash
	TxHash []byte
	// 2.交易输入
	Vins []*TXInput
	// 3. 交易输出
	Vouts []*TXOutput
}

//Coinbase时的Transaction
func NewConbaseTransaction(address string) *Transaction {
	txInput := &TXInput{[]byte{}, -1, "Coinbase"}
	txOutput := &TXOutput{10, address}
	txCoinbase := &Transaction{[]byte{}, []*TXInput{txInput}, []*TXOutput{txOutput}}

	txCoinbase.HashTransaction()
	return txCoinbase
}

func (tx *Transaction) HashTransaction()  {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err:= encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	hash := sha256.Sum256(result.Bytes())

	tx.TxHash = hash[:]
}