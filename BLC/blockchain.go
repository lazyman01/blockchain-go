package BLC

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
)

// 数据库名
const dbFile = "blockchain.db"
// 仓库
const blocksBucket = "blocks"


type BlockChain struct {
	//Blocks []*Block  //存储有序的区块
	Tip []byte //区块链里面最后一个hash
	DB *bolt.DB // 数据库
}

//新增区块
func (blockchain *BlockChain) AddBlock(data string)  {
	//1.创建新的区块
	newBlock := NewBlock(data, blockchain.Tip)
	// 2. Update数据
	err := blockchain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		//存储新区块
		err := b.Put(newBlock.Hash,newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		//更新最后hash值
		err = b.Put([]byte("last"), newBlock.Hash)
		if err != nil{
			log.Panic(err)
		}

		// 更新区块链的结构体中Tip
		blockchain.Tip = newBlock.Hash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

//创建一个带有创世区块的区块链
func NewBlockChain(data string) *BlockChain  {
	var tip []byte //存储最后一个区块的hash
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
		b := tx.Bucket([]byte(blocksBucket))

		//表不存在
		if b == nil {
			fmt.Println("no exsiting blockchain found. Creating  now")
			genesisBlock := NewGenesisBlock(data)
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}
			err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic()
			}

			err = b.Put([]byte("last"), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesisBlock.Hash
		} else {//表存在
			tip = b.Get([]byte("last"))
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	return &BlockChain{tip, db}
}

//迭代器结构
type BlockchainIterator struct {
	CunrrentHash []byte //当前正在遍历的区块hash
	DB *bolt.DB // 数据库
}

//迭代器
func (blockchain *BlockChain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{blockchain.Tip, blockchain.DB}
}

//获取下一个迭代器
func (bi *BlockchainIterator) Next() *BlockchainIterator {
	var nextHash []byte
	err := bi.DB.View(func(tx *bolt.Tx) error {
		//获取表
		t := tx.Bucket([]byte(blocksBucket))

		//通过当前的hash获取block
		currentBlockBytes := t.Get(bi.CunrrentHash)
		//反序列化
		currentBlock := DeserializeBlock(currentBlockBytes)
		nextHash = currentBlock.PrevBlockHash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &BlockchainIterator{nextHash, bi.DB}
}

func (blockchain *BlockChain) Printchain()  {
	//fmt.Printf("\ntip: %x\n", blockchain.Tip)
	var blockchainIterator *BlockchainIterator
	blockchainIterator = blockchain.Iterator()
	var hashInt big.Int
	for {
		//fmt.Printf("%x\n", blockchainIterator.CunrrentHash)
		hashInt.SetBytes(blockchainIterator.CunrrentHash[:])
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
		err := blockchainIterator.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("blocks"))
			blockBytes := b.Get(blockchainIterator.CunrrentHash)
			currentBlock := DeserializeBlock(blockBytes)
			fmt.Printf("\nPreBlockhash: %x\n",currentBlock.PrevBlockHash)
			fmt.Printf("Timestamp:     %x\n", currentBlock.Timestamp)
			fmt.Printf("Nonce:         %d\n", currentBlock.Nonce)
			fmt.Printf("Data:          %s\n", currentBlock.Data)
			fmt.Printf("Hash:          %x\n", currentBlock.Hash)

			return nil
		})
		if err != nil {
			log.Panic(err)
		}

		blockchainIterator = blockchainIterator.Next()
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
		b := tx.Bucket([]byte(blocksBucket))
		if b != nil {
			tip = b.Get([]byte("last"))
		}
		return nil
	})

	return &BlockChain{tip, db}
}