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

func (cli *CLI) addBlock(txs []*Transaction)  {
	if DBExists() == false {
		fmt.Println("数据不存在")
		os.Exit(1)
	}

	blockchain := GetBlockchian()
	blockchain.AddBlock(txs)
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

func (cli *CLI) createGenesisBlockchain(address string)  {
	NewBlockChain(address)
}

// 转账
func (cli *CLI) send(from []string,to []string,amount []string)  {

	MineNewBlock(from,to,amount)
}

func (cli *CLI) Run()  {
	sendBlockCmd := flag.NewFlagSet("send",flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain",flag.ExitOnError)
	//addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	flagFrom := sendBlockCmd.String("from","","转账源地址......")
	flagTo := sendBlockCmd.String("to","","转账目的地地址......")
	flagAmount := sendBlockCmd.String("amount","","转账金额......")

	//addBlockData := addBlockCmd.String("data", "", "Block Data")
	createBlockchainData := createBlockchainCmd.String("address", "Genesis Block", "创世区块地址")
	isValidargs()
	switch os.Args[1] {
	//case "addblock":
	//	err := addBlockCmd.Parse(os.Args[2:])
	//	if err != nil {
	//		log.Panic(err)
	//	}
	case "send":
		err := sendBlockCmd.Parse(os.Args[2:])
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

	//if addBlockCmd.Parsed() {
	//	if *addBlockData == "" {
	//		printUsage()
	//		return
	//	}
	//	cli.addBlock(*addBlockData)
	//}
	if sendBlockCmd.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flagAmount == "" {
			printUsage()
			os.Exit(1)
		}

		from := JSONToArray(*flagFrom)
		to := JSONToArray(*flagTo)
		amount := JSONToArray(*flagAmount)
		cli.send(from,to,amount)
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