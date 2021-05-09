package BLC

import (
	"bytes"
	"crypto/sha256"
	"strconv"
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
}

func (b *Block) SetHash()  {
	//1.将时间戳长整型转字符串再转字节数组
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	//2.将除了hash以外的其他属性，以字节数组的形式全拼接起来
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	//3.将拼接起来的数据进行256hash
	hash := sha256.Sum256(headers)
	//4.将hash赋给hash
	b.Hash = hash[:]
}
func NewBlock(data string, prevBlockHash []byte) *Block  {
	//创建区块
	block := &Block{time.Now().Unix(), prevBlockHash,
		[]byte(data), []byte{}}
	//设置当前区块hash值
	block.SetHash()
	//返回hash
	return block
}

func NewGenesisBlock() *Block  {
	return NewBlock("Ggnenis Block", []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
}