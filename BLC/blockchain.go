package BLC

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"time"
)

// 数据库名
const dbFile = "blockchain.db"
// 仓库
const blockTableName = "blocks"


type BlockChain struct {
	//Blocks []*Block  //存储有序的区块
	Tip []byte //区块链里面最后一个hash
	DB *bolt.DB // 数据库
}

//新增区块
func (blockchain *BlockChain) AddBlock(txs []*Transaction)  {
	//1.创建新的区块
	//newBlock := NewBlock(data, blockchain.Tip)
	// 2. Update数据
	err := blockchain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			blockBytes := b.Get(blockchain.Tip)
			block := DeserializeBlock(blockBytes)

			newblock := NewBlock(txs, block.Height+1, block.Hash)
			err := b.Put(newblock.Hash, newblock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("1"), newblock.Hash)
			if err != nil {
				log.Panic(err)
			}
			blockchain.Tip = newblock.Hash
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

//创建一个带有创世区块的区块链
func NewBlockChain(address string)  {
	if DBExists() {
		fmt.Println("创世块已经存在...")
		os.Exit(1)
	}

	//var tip []byte //存储最后一个区块的hash
	// 1. 尝试打开或者创建数据库
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	//defer db.Close()
	// 2. db.update 更新数据
	//  2.1 表是否存在，如果不存在，需要创建表
	//  2.2 创建创世区块
	//  2.3 需要将创世区块序列化
	//  2.4 把创世区块的hash作为key, Block序列化的数据作为value存储到表里
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blockTableName))
		if err != nil {
			log.Panic(err)
		}

		//表不存在
		if b != nil {
			txCoinbase := NewConbaseTransaction(address)
			genesisBlock := NewGenesisBlock([]*Transaction{txCoinbase})

			//将创世区块存储到表中
			err := b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			//存储最新区块的hash
			err = b.Put([]byte("1"), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
}

//迭代器
func (blockchain *BlockChain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{blockchain.Tip, blockchain.DB}
}

func (blockchain *BlockChain) Printchain()  {
	//fmt.Printf("\ntip: %x\n", blockchain.Tip)
	//var blockchainIterator *BlockchainIterator
	blockchainIterator := blockchain.Iterator()

	for {
		block := blockchainIterator.GetOneBlock()
		fmt.Printf("Height:%d\n", block.Height)
		fmt.Printf("PrevBlockhash:%x\n", block.PrevBlockHash)
		fmt.Printf("Timestamp:%s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Println("Txs:")
		for _, tx := range block.Txs {
			fmt.Printf("%x\n", tx.TxHash)
			fmt.Println("Vins:")
			for _, in := range tx.Vins {
				fmt.Printf("%x\n", in.TxHash)
				fmt.Printf("%d\n", in.Vout)
				fmt.Printf("%s\n", in.ScriptSig)
			}

			fmt.Println("Vouts:")
			for _, out := range tx.Vouts {
				fmt.Println(out.Value)
				fmt.Println(out.ScriptPubKey)
			}
		}

		fmt.Println("-------------------------------------")
		var hashInt big.Int
		hashInt.SetBytes(blockchainIterator.CurrentHash[:])
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
}

func DBExists() bool  {
	if _, err := os.Stat(dbFile);os.IsNotExist(err) {
		return false
	}
	return true
}

func GetBlockchian() *BlockChain  {
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var tip []byte

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			tip = b.Get([]byte("last"))
		}
		return nil
	})

	return &BlockChain{tip, db}
}

func MineNewBlock(from []string, to []string, amount []string)  {
	fmt.Println(from, to, amount)
}