package main

import (
	"blockchain-go/BLC"
)

func main()  {
	//创建区块链
	//blockchain := BLC.NewBlockChain()
	//defer blockchain.DB.Close()
	//创建CLI对象
	cli := BLC.CLI{}
	//调用CLI的Run方法
	cli.Run()
}
