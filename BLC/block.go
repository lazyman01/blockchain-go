package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	//块高度
	Height int64

	//时间戳,创建区块是的时间
	Timestamp int64

	//上一个区块的hash,父hash
	PrevBlockHash []byte

	//交易数据
	Txs []*Transaction

	//Hash, 当前区块hash
	Hash []byte

	//Nonce 随机数
	Nonce int
}

func (block *Block) HashTransaction() []byte  {
	var txHashes[][]byte
	var txHash [32]byte

	for _, tx := range block.Txs {
		txHashes = append(txHashes, tx.TxHash)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

//将Block对象序列化成[]byte
func (block *Block) Serialize() []byte  {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
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

func NewBlock(txs []*Transaction, height int64, prevBlockHash []byte) *Block  {
	//创建区块
	block := &Block{height, time.Now().Unix(), prevBlockHash,
		txs, []byte{}, 0}
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
func NewGenesisBlock(txs []*Transaction) *Block  {
	return NewBlock(txs, 1, []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
}