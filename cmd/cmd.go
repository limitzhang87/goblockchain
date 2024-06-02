package cmd

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/limitzhang87/goblockchain/blockchain"
	"github.com/limitzhang87/goblockchain/utils"
	"github.com/limitzhang87/goblockchain/wallet"
	"os"
	"runtime"
	"strconv"
	"time"
)

const (
	FlagCreateBlockchain  = "createblockchain"
	FlagCreateWallet      = "createwallet"
	FlagWalletInfo        = "walletinfo"
	FlagWalletsUpdate     = "walletsupdate"
	FlagWalletsList       = "walletslist"
	FlagWalletBindRefName = "walletbindrefname"
	FlagBalance           = "balance"
	FlagBlockChainInfo    = "blockchaininfo"
	FlagSend              = "send"
	FlagSendByRefName     = "sendbyrefname"
	FlagMine              = "mine"
)

type CommandLine struct {
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Welcome to Limit's tiny blockchain system, usage is as follows:")
	fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("All you need is to first create a wallet.")
	fmt.Println("And then you can use the wallet address to create a blockchain and declare the owner.")
	fmt.Println("Make transactions to expand the blockchain.")
	fmt.Println("In addition, don't forget to run mine function after transatcions are collected.")
	fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("createwallet -refname REFNAME                       ----> Creates and save a wallet. The refname is optional.")
	fmt.Println("walletinfo -refname NAME -address Address           ----> Print the information of a wallet. At least one of the refname and address is required.")
	fmt.Println("walletsupdate                                       ----> Registrate and update all the wallets (especially when you have added an existed .wlt file).")
	fmt.Println("walletslist                                         ----> List all the wallets found (make sure you have run walletsupdate first).")
	fmt.Println("walletbindrefname                                   ----> Bind address to refname.")
	fmt.Println("createblockchain -refname NAME -address ADDRESS     ----> Creates a blockchain with the owner you input (address or refname).")
	fmt.Println("balance -refname NAME -address ADDRESS              ----> Back the balance of a wallet using the address (or refname) you input.")
	fmt.Println("blockchaininfo                                      ----> Prints the blocks in the chain.")
	fmt.Println("send -from FROADDRESS -to TOADDRESS -amount AMOUNT  ----> Make a transaction and put it into candidate block.")
	fmt.Println("sendbyrefname -from NAME1 -to NAME2 -amount AMOUNT  ----> Make a transaction and put it into candidate block using refname.")
	fmt.Println("mine                                                ----> Mine and add a block to the chain.")
	fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) Run() {
	cli.validateArgs()

	createBlockchainCmd := flag.NewFlagSet(FlagCreateBlockchain, flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet(FlagCreateWallet, flag.ExitOnError)
	walletInfoCmd := flag.NewFlagSet(FlagWalletInfo, flag.ExitOnError)
	walletBindRefNameCmd := flag.NewFlagSet(FlagWalletBindRefName, flag.ExitOnError)
	balanceCmd := flag.NewFlagSet(FlagBalance, flag.ExitOnError)
	//getBlockCmd := flag.NewFlagSet(FlagBlockChainInfo, flag.ExitOnError)
	sendCmd := flag.NewFlagSet(FlagSend, flag.ExitOnError)
	//mineCmd := flag.NewFlagSet(FlagMine, flag.ExitOnError)

	switch os.Args[1] {
	case FlagCreateBlockchain:
		refName := createBlockchainCmd.String("refname", "", "The refName refer to the owner of blockchain")
		err := createBlockchainCmd.Parse(os.Args[2:]) // 解析子参数后面的其他参数address， 并传给createBlockchainAddress(这个是一个指针)
		utils.Handle(err)

		if len(*refName) == 0 {
			fmt.Println("Please enter a valid refname")
			return
		}
		cli.createBlockchainRefName(*refName)
	case FlagCreateWallet:
		refName := createWalletCmd.String("refname", "", "The refName refer to the owner of blockchain")
		err := createWalletCmd.Parse(os.Args[2:])
		utils.Handle(err)
		cli.createWallet(*refName)
	case FlagWalletInfo:
		walletInfoAddress := walletInfoCmd.String("address", "", "The address of the wallet")
		walletInfoRefName := walletInfoCmd.String("refname", "", "The refname of the wallet")
		err := walletInfoCmd.Parse(os.Args[2:])
		utils.Handle(err)

		if len(*walletInfoRefName) == 0 && len(*walletInfoAddress) == 0 {
			fmt.Println("Please enter a valid address or refname")
			return
		}
		cli.walletInfo(*walletInfoAddress, *walletInfoRefName)
	case FlagWalletsUpdate:
		cli.walletsUpdate()
	case FlagWalletsList:
		cli.walletsList()
	case FlagWalletBindRefName:
		address := walletBindRefNameCmd.String("address", "", "The address of the wallet")
		refName := walletBindRefNameCmd.String("refname", "", "The refname of the wallet")
		err := walletBindRefNameCmd.Parse(os.Args[2:])
		utils.Handle(err)
		if len(*refName) == 0 && len(*address) == 0 {
			fmt.Println("Please enter a valid address or refname")
			return
		}
		cli.walletBindRefName(*address, *refName)
	case FlagBalance:
		address := balanceCmd.String("address", "", "The address of the wallet")
		refName := balanceCmd.String("refname", "", "Who need to get balance amount")
		err := balanceCmd.Parse(os.Args[2:])
		utils.Handle(err)
		if len(*refName) == 0 && len(*address) == 0 {
			fmt.Println("Please enter a valid refName or address")
			return
		}
		if len(*refName) != 0 {
			cli.balanceRefName(*refName)
		} else {
			cli.balance(*address)
		}
	case FlagBlockChainInfo:
		cli.getBlockchainInfo()
	case FlagSend:
		sendFromAddress := sendCmd.String("from", "", "Source address")
		sendToAddress := sendCmd.String("to", "", "Destination address")
		sendAmount := sendCmd.Int("amount", 0, "Amount to send")
		err := sendCmd.Parse(os.Args[2:])
		utils.Handle(err)
		if len(*sendFromAddress) == 0 {
			fmt.Println("Please enter a valid from address")
		}
		if len(*sendToAddress) == 0 {
			fmt.Println("Please enter a valid to address")
		}
		if *sendAmount <= 0 {
			fmt.Println("Please enter a valid amount")
		}
		cli.send(*sendFromAddress, *sendToAddress, *sendAmount)
	case FlagSendByRefName:
		sendFromRefName := sendCmd.String("from", "", "Source refName")
		sendToRefName := sendCmd.String("to", "", "Destination refName")
		sendAmount := sendCmd.Int("amount", 0, "Amount to send")
		err := sendCmd.Parse(os.Args[2:])
		utils.Handle(err)
		if len(*sendFromRefName) == 0 {
			fmt.Println("Please enter a valid from address")
		}
		if len(*sendToRefName) == 0 {
			fmt.Println("Please enter a valid to address")
		}
		if *sendAmount <= 0 {
			fmt.Println("Please enter a valid amount")
		}
		cli.sendByName(*sendFromRefName, *sendToRefName, *sendAmount)
	case FlagMine:
		cli.mine()
	default:
		cli.printUsage()
	}

}

// createBlockchain 创建区块链
func (cli *CommandLine) createBlockchain(address string) {
	newChain := blockchain.InitBlockChain(utils.Address2PubHash([]byte(address)))
	_ = newChain.Database.Close()
	fmt.Println("Created new blockchain: owner is : ", address)
}

// createBlockchainRefName
func (cli *CommandLine) createBlockchainRefName(refName string) {
	address := cli.getAddressByRefName(refName)
	cli.createBlockchain(address)
}

// createWallet 创建钱包
func (cli *CommandLine) createWallet(refName string) {
	wlt := wallet.NewWallet()
	wlt.Save()
	refList := wallet.LoadRefList()
	refList.BindRef(string(wlt.Address()), refName)
	refList.Save()
	fmt.Println("Succeed in creating wallet")
	fmt.Printf("Wallet address:%x\n", wlt.Address())
	fmt.Printf("Public Key:%x\n", wlt.PublicKey)
}

// walletInfo 查询钱包信息
func (cli *CommandLine) walletInfo(address, refName string) {
	refList := wallet.LoadRefList()
	if address != "" {
		refName = (*refList)[address]
	}
	if refName != "" {
		for k, v := range *refList {
			if v == refName {
				address = k
			}
		}
	}
	wlt := wallet.LoadWallet(address)
	fmt.Printf("Wallet address:%x\n", wlt.Address())
	fmt.Printf("Public Key:%x\n", wlt.PublicKey)
	fmt.Printf("Reference Name:%s\n", (*refList)[address])
}

// walletsUpdate 更新钱包管理文件
func (cli *CommandLine) walletsUpdate() {
	refList := wallet.LoadRefList()
	refList.Update()
	refList.Save()
}

// walletsList 钱包列表
func (cli *CommandLine) walletsList() {
	refList := wallet.LoadRefList()
	for address, refName := range *refList {
		wlt := wallet.LoadWallet(address)
		fmt.Println("--------------------------------------------------------------------------------------------------------------")
		fmt.Printf("Wallet address:%s\n", address)
		fmt.Printf("Wallet refName:%s\n", refName)
		fmt.Printf("Public Key:%x\n", wlt.PublicKey)
		fmt.Println("--------------------------------------------------------------------------------------------------------------")
		fmt.Println()
	}
}

// walletBindRefName 钱包绑定别名
func (cli *CommandLine) walletBindRefName(address, refName string) {
	_ = wallet.LoadWallet(address) // 加载钱包，用于判断address是合法的

	refList := wallet.LoadRefList()
	// 判断refName是否已经存在
	a, err := refList.FindRef(refName)
	if err == nil {
		fmt.Printf("refname had already bind to %s\n", a)
		return
	}
	refList.BindRef(address, refName)
	refList.Save()
	fmt.Printf("Bind success")
}

// balance 查询账户余额 address
func (cli *CommandLine) balance(address string) {
	chain := blockchain.ContinueBlockChain()
	defer func() {
		_ = chain.Database.Close()
	}()
	wlt := wallet.LoadWallet(address)
	balance, _ := chain.FindUTXOs(wlt.PublicKey)
	fmt.Println("Address is : ", address, "Balance : ", balance)
}

// balanceRefName 根据昵称查询余额
func (cli *CommandLine) balanceRefName(refName string) {
	address := cli.getAddressByRefName(refName)
	cli.balance(address)
}

// getBlockchainInfo 输出区块信息
func (cli *CommandLine) getBlockchainInfo() {
	chain := blockchain.ContinueBlockChain()
	iter := chain.Iterator()
	defer func() {
		_ = chain.Database.Close()
	}()
	ogPrevHash := chain.BackOgPrevHash()
	for {
		block := iter.Next()
		fmt.Println("---------------------------------------------------------------------------------------------")
		fmt.Printf("Timestamp:%s\n", time.Unix(block.Timestamp, 0).Format(time.DateTime))
		fmt.Printf("Previous hash:%x\n", block.PrevHash)
		fmt.Printf("Transactions:%v\n", block.Transactions)
		for i, tx := range block.Transactions {
			fmt.Printf("\tTransaction Index:%d\n", i)
			fmt.Println("\tInput:")
			for _, in := range tx.Inputs {
				fmt.Printf("\t\tTxID:%s\n", hex.EncodeToString(in.TxID))
				fmt.Printf("\t\tOutIdx:%d\n", in.OutIdx)
				fmt.Printf("\t\tPubKey:%x\n", in.PubKey)
				fmt.Printf("\t\tAddress:%s\n", string(utils.PubHash2Address(utils.PublicKeyHash(in.PubKey))))
			}
			fmt.Println("\tOutput:")
			for _, out := range tx.Outputs {
				fmt.Printf("\t\tPubKeyHash:%s\n", hex.EncodeToString(out.PubKeyHash))
				fmt.Printf("\t\tValue:%d\n", out.Value)
				fmt.Printf("\t\tAddress:%s\n", string(utils.PubHash2Address(out.PubKeyHash)))
			}
		}
		fmt.Printf("hash:%x\n", block.Hash)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(block.ValidatePoW()))
		fmt.Println("---------------------------------------------------------------------------------------------")
		fmt.Println()
		if bytes.Equal(ogPrevHash, block.PrevHash) {
			break
		}
	}
}

// send 发送交易 from,to 钱包地址
func (cli *CommandLine) send(from, to string, amount int) {
	chain := blockchain.ContinueBlockChain()
	defer func() {
		_ = chain.Database.Close()
	}()

	fromWallet := wallet.LoadWallet(from)

	tx, err := chain.CreateTransaction(fromWallet.PublicKey, utils.Address2PubHash([]byte(to)), amount, fromWallet.PrivateKey)
	if err != nil {
		fmt.Println("Create transaction error : ", err)
		return
	}

	pool := blockchain.CreatePool()
	pool.AddTransaction(tx)
	pool.SaveFile()
	fmt.Println("success")
}

// sendByName 根据名字交易
func (cli *CommandLine) sendByName(nameFrom, toFrom string, amount int) {
	refList := wallet.LoadRefList()
	addressFrom, err := refList.FindRef(nameFrom)
	utils.Handle(err)
	addressTo, err := refList.FindRef(toFrom)
	utils.Handle(err)
	cli.send(addressFrom, addressTo, amount)
}

// getAddressByRefName 根据别名查找钱包地址
func (cli *CommandLine) getAddressByRefName(refName string) string {
	refList := wallet.LoadRefList()
	address, err := refList.FindRef(refName)
	utils.Handle(err)
	return address
}

// mine 挖矿
func (cli *CommandLine) mine() {
	chain := blockchain.ContinueBlockChain()
	defer func() {
		_ = chain.Database.Close()
	}()

	chain.RunMine()
	fmt.Println("Finish Mining")
}
