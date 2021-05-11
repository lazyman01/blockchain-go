package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	//时间戳,创建区块是的时间
	Timestamp int64
	//上一个区块的hash,父hash
	PrevBlockHash []byte
	//data 交易数据
	Data []byte
	//Hash, 当前区块hash
	Hash []byte
	//Nonce 随机数
	Nonce int
}

//将Block对象序列化成[]byte
func (b *Block) Serialize() []byte  {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

//将字节数组反序列化成Block
func DeserializeBlock(d []byte) *Block  {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

func NewBlock(data string, prevBlockHash []byte) (*Block)  {
	//创建区块
	block := &Block{time.Now().Unix(), prevBlockHash,
		[]byte(data), []byte{}, 0}
	//创建一个pow对象
	pow := NewProofOfWork(block)
	//Run()执行一次工作量证明
	nonce, hash := pow.Run()
	//设置区块hash
	block.Hash = hash[:]
	//设置Nonce
	block.Nonce = nonce

	isValid := pow.Validate()

	if isValid {
		return block
	} else {
		log.Printf("block valid error: %v\n", block)
		return nil
	}

}

//生成创世块
func NewGenesisBlock(data string) *Block  {
	return NewBlock(data, []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
}