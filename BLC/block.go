package BLC

import (
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

//func (b *Block) SetHash()  {
//	//1.将时间戳长整型转字符串再转字节数组
//	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
//	//2.将除了hash以外的其他属性，以字节数组的形式全拼接起来
//	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
//	//3.将拼接起来的数据进行256hash
//	hash := sha256.Sum256(headers)
//	//4.将hash赋给hash
//	b.Hash = hash[:]
//}
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
func NewGenesisBlock() *Block  {
	return NewBlock("Ggnenis Block", []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
}