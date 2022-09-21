package blc

import (
	"blockDemo/blc/utils"
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct{}

func (cli *CLI) Run() {
	isValid()

	nodeId := os.Getenv("NODE_ID")
	if nodeId == "" {
		fmt.Println("Not set nodeId")
		os.Exit(1)
	}
	fmt.Printf("nodeId=%s\n", nodeId)

	// 新建命令行
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	getwalletlistCmd := flag.NewFlagSet("getwalletlist", flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)
	testCmd := flag.NewFlagSet("test", flag.ExitOnError)

	switch os.Args[1] {
	case "send":
		{
			flagFrom := sendBlockCmd.String("from", "", "转账源地址")
			flagTo := sendBlockCmd.String("to", "", "转账目的地址")
			flagAmount := sendBlockCmd.String("amount", "", "转账金额")
			flagMine := sendBlockCmd.Bool("mine", false, "是否在当前节点中立即验证")

			err := sendBlockCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}

			if *flagFrom == "" || *flagTo == "" || *flagAmount == "" {
				printUsage()
				os.Exit(1)
			}
			//fmt.Printf("%s\n", *flagAddBlockData)
			//cli.addBlock([]*Transaction{})

			// fmt.Println(utils.JsonToArray(*flagFrom))
			// fmt.Println(utils.JsonToArray(*flagTo))
			// fmt.Println(utils.JsonToArray(*flagAmount))

			from := utils.JsonToArray(*flagFrom)
			to := utils.JsonToArray(*flagTo)

			for idx, addr := range from {
				if !IsValidForAddress([]byte(addr)) || !IsValidForAddress([]byte(to[idx])) {
					log.Println("地址无效...")
					os.Exit(1)
				}
			}
			amount := utils.JsonToArray(*flagAmount)
			cli.send(from, to, amount, nodeId, *flagMine)
		}
	case "printchain":
		{
			err := printChainCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}

			cli.printChain(nodeId)
		}
	case "createblockchain":
		{
			flagCreateBlockChain := createBlockChainCmd.String("address", "genesis block", "创世纪区块")

			err := createBlockChainCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}

			if !IsValidForAddress([]byte(*flagCreateBlockChain)) {
				log.Println("地址无效...")
				printUsage()
				os.Exit(1)
			}
			cli.createBlockChainWithGenesis(*flagCreateBlockChain, nodeId)
		}
	case "getbalance":
		{
			flagBalance := getBalanceCmd.String("address", "0", "查询账号余额")
			err := getBalanceCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
			if !IsValidForAddress([]byte(*flagBalance)) {
				log.Println("地址无效...")
				printUsage()
				os.Exit(1)
			}

			cli.getBalance(*flagBalance, nodeId)
		}
	case "createwallet":
		{
			err := createWalletCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
			cli.CreateWallet(nodeId)
		}
	case "getwalletlist":
		{
			err := getwalletlistCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
			cli.GetWalletlist(nodeId)
		}
	case "startnode":
		{
			flagMinerAdd := startNodeCmd.String("miner", "", "定义挖坑奖励的地址")
			err := startNodeCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
			cli.StartNode(nodeId, *flagMinerAdd)
		}
	case "test":
		{
			err := testCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
			cli.TestMethod(nodeId)
		}
	default:
		{
			printUsage()
			os.Exit(1)
		}

	}
}

func printUsage() {
	fmt.Println("Usage:")

	fmt.Println("\tcreatewallet  --创建钱包")
	fmt.Println("\tcreateblockchain -address DATA --交易数据")
	fmt.Println("\tsend -from From -to To -amount Amount --交易明细")
	fmt.Println("\tgetbalance -address --查询余额")
	fmt.Println("\tprintchain --输出区块信息")
	fmt.Println("\tgetwalletlist  --输出所有钱包")
	fmt.Println("\tstartnode -address --启动节点服务器,并且指定挖矿奖励地址")
	fmt.Println("\ttest  --测试")
}

func isValid() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}
