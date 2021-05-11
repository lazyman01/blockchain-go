package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
	Blockchain *BlockChain
}

func isValidargs()  {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

func printUsage()  {
	fmt.Println("Usage:")
	fmt.Println("\taddblock -data DATA -- 交易数据.")
	fmt.Println("\tprintchain -- 输出区块信息.")
}

func (cli *CLI) addBlock(data string)  {
	blockchain := GetBlockchian()
	blockchain.AddBlock(data)
	defer blockchain.DB.Close()
}

func (cli *CLI) printchain()  {
	if DBExists() == false {
		fmt.Println("数据库文件不存在...")
		return
	}
	blockchain := GetBlockchian()

	defer blockchain.DB.Close()

	blockchain.Printchain()
}

func (cli *CLI) createGenesisBlockchain(data string)  {
	NewBlockChain(data)
}

func (cli *CLI) Run()  {
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block Data")
	createBlockchainData := createBlockchainCmd.String("data", "Genesis Block", "gennesis block")
	isValidargs()
	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		fmt.Println("\nNo addblock and printchain!")
		return
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			printUsage()
			return
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printchain()
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainData == "" {
			fmt.Println("交易数据不能为空")
			printUsage()
			return
		}
		cli.createGenesisBlockchain(*createBlockchainData)
	}


}