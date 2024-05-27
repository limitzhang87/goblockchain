package main

import (
	"fmt"
	"github.com/limitzhang87/goblockchain/blockchain"
	"github.com/limitzhang87/goblockchain/transaction"
)

// 转载https://zhuanlan.zhihu.com/p/413585427

func main() {
	txPool := make([]*transaction.Transaction, 0)
	var tempTx *transaction.Transaction
	var err error
	var property int
	chain := blockchain.CreateBlockChain()
	property, _ = chain.FindUTXOs([]byte("limitZhang"))
	fmt.Println("Balance of limit zhang", property)

	tempTx, err = chain.CreateTransaction([]byte("limitZhang"), []byte("zhang"), 100)
	if err == nil {
		txPool = append(txPool, tempTx)
	}
	chain.Mine(txPool)
	txPool = make([]*transaction.Transaction, 0)

	property, _ = chain.FindUTXOs([]byte("limitZhang"))
	fmt.Println("Balance of limit zhang", property)

	tempTx, err = chain.CreateTransaction([]byte("zhang"), []byte("li"), 200)
	if err == nil {
		txPool = append(txPool, tempTx)
	}

	tempTx, err = chain.CreateTransaction([]byte("zhang"), []byte("li"), 50)
	if err == nil {
		txPool = append(txPool, tempTx)
	}

	tempTx, err = chain.CreateTransaction([]byte("limitZhang"), []byte("li"), 100)
	if err == nil {
		txPool = append(txPool, tempTx)
	}
	chain.Mine(txPool)
	txPool = make([]*transaction.Transaction, 0)

	property, _ = chain.FindUTXOs([]byte("limitZhang"))
	fmt.Println("Balance of limit zhang: ", property)
	property, _ = chain.FindUTXOs([]byte("zhang"))
	fmt.Println("Balance of zhang: ", property)
	property, _ = chain.FindUTXOs([]byte("li"))
	fmt.Println("Balance of li: ", property)

	fmt.Println("===================================================")

	for _, block := range chain.Blocks {
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("hash: %x\n", block.Hash)
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("nonce: %d\n", block.Nonce)
		fmt.Println("Proof of Work validation:", block.ValidatePoW())
	}
	fmt.Println("===================================================")

	//展示当前版本的BUG， zhang经过上面的交易，只剩下50
	// 但是在接下来的直接创建两次交易，zhang分别转给两个人各30(总共60)，结果能够执行下去，这是有问题的
	// 这是因为第一笔交易 zhang->li:30了之后还没有入块，本来只剩下20，
	// 但zhang->limitZhang:30时，还是直接从区块链中取余额，余额还是50
	tempTx, err = chain.CreateTransaction([]byte("zhang"), []byte("li"), 30)
	if err == nil {
		txPool = append(txPool, tempTx)
	}

	tempTx, err = chain.CreateTransaction([]byte("zhang"), []byte("limitZhang"), 30)
	if err == nil {
		txPool = append(txPool, tempTx)
	}

	chain.Mine(txPool)
	txPool = make([]*transaction.Transaction, 0)
	for _, block := range chain.Blocks {
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("hash: %x\n", block.Hash)
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("nonce: %d\n", block.Nonce)
		fmt.Println("Proof of Work validation:", block.ValidatePoW())
	}
	fmt.Println("===================================================")
	property, _ = chain.FindUTXOs([]byte("limitZhang"))
	fmt.Println("Balance of limit zhang: ", property)
	property, _ = chain.FindUTXOs([]byte("zhang"))
	fmt.Println("Balance of zhang: ", property)
	property, _ = chain.FindUTXOs([]byte("li"))
	fmt.Println("Balance of li: ", property)
}
