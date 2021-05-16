package BLC

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

type BlockchainIterator struct{
	CurrentHash []byte
	DB *bolt.DB
}

func (BlockchainIterator *BlockchainIterator) GetOneBlock() *Block  {
	var block *Block
	err := BlockchainIterator.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b!= nil {
			currentBlockBytes := b.Get(BlockchainIterator.CurrentHash)
			fmt.Println(currentBlockBytes)
			block = DeserializeBlock(currentBlockBytes)
			BlockchainIterator.CurrentHash = block.PrevBlockHash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return block
}
