package blockchain

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/limitzhang87/goblockchain/transaction"
	"github.com/limitzhang87/goblockchain/utils"
	"log"
)

func (bc *BlockChain) RunMine() {
	pool := CreatePool()
	if !bc.VerityTransaction(pool.Txs) {
		log.Println("falls in transactions verification")
		err := RemovePoolFile()
		utils.Handle(err)
		return
	}

	block := CreateBlock(bc.LastHash, pool.Txs)
	if !block.ValidatePoW() {
		fmt.Println("Block has invalid nonce!")
		return
	}
	bc.AddBlock(block)
	err := RemovePoolFile()
	utils.Handle(err)
}

// VerityTransaction 验证交易是否有效
func (bc *BlockChain) VerityTransaction(txs []*transaction.Transaction) bool {
	// 1. 交易输入不能重复使用
	// 2. 交易输入要有效
	// 3. 输入金额要等于输出金额

	inAmount, outAmount := 0, 0
	spentInput := make(map[string]int)
	for _, tx := range txs {
		if len(tx.Inputs) == 0 {
			fmt.Println("Inputs is empty!")
			return false
		}
		pubKey := tx.Inputs[0].PubKey
		unspentTx := bc.FindUnspentTransactions(pubKey)

		// 1. 交易输入不能重复使用
		for _, input := range tx.Inputs {
			inTxId := hex.EncodeToString(input.TxID)

			// 交易输入重复使用，直接返回失败
			if txId, ok := spentInput[inTxId]; ok && txId == input.OutIdx {
				fmt.Println("input had already spent")
				return false
			}
			v, ok := bc.isInputRight(unspentTx, input)
			if !ok {
				fmt.Println("input not right")
				return false
			}
			inAmount += v
			spentInput[inTxId] = input.OutIdx
		}

		for _, output := range tx.Outputs {
			outAmount += output.Value
		}
	}

	if inAmount != outAmount {
		fmt.Println("inAmount != outAmount")
		return false
	}
	return true
}

// isInputRight 判断交易输入是否有效，是否是当前用于在块中的未花费交易输出
func (bc *BlockChain) isInputRight(txs []*transaction.Transaction, input transaction.TxInput) (int, bool) {
	for _, tx := range txs {
		if bytes.Equal(tx.ID, input.TxID) && input.OutIdx < len(tx.Outputs) {
			return tx.Outputs[input.OutIdx].Value, true
		}
	}
	return 0, false
}
