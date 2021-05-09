package main

import (
	"blockchain-go/BLC"
	"fmt"
	"time"
)
func main()  {
	//创建创世区块
	//block := BLC.NewGenesisBlock()
	blockchain := BLC.NewBlockChain()
	blockchain.AddBlock("Send 20 BTC To HaoLin")
	blockchain.AddBlock("Send 30 BTC To LaoWang")
	blockchain.AddBlock("Send 10 BTC To XiaoSao")

	for _, block := range blockchain.Blocks {
		fmt.Println("timestamp:", time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("data:%s\n", block.Data)
		fmt.Printf("preHash:%x\n", block.PrevBlockHash)
		fmt.Printf("Hash:%x\n", block.Hash)
	}
}
