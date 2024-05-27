package blockchain

import (
	"encoding/hex"
	"errors"
	"github.com/limitzhang87/goblockchain/transaction"
	"github.com/limitzhang87/goblockchain/utils"
)

type BlockChain struct {
	Blocks []*Block
}

// AddBlock 区块链添加区块
func (bc *BlockChain) AddBlock(txs []*transaction.Transaction) {
	newBlock := CreateBlock(bc.Blocks[len(bc.Blocks)-1].Hash, txs)
	bc.Blocks = append(bc.Blocks, newBlock)
}

// CreateBlockChain 创建区块链
func CreateBlockChain() *BlockChain {
	blockchain := BlockChain{}
	blockchain.Blocks = append(blockchain.Blocks, GenesisBlock())
	return &blockchain
}

// FindUnspentTransactions 根据帐号找出未使用的交易
func (bc *BlockChain) FindUnspentTransactions(address []byte) []*transaction.Transaction {
	unSpentTxs := make([]*transaction.Transaction, 0)
	spentIds := make(map[string][]int) // 用于保存已经使用的交易输出 string是交易ID， []int交易下标

	// 1 遍历整个区块链,从最后一个区块开始
	for i := len(bc.Blocks) - 1; i >= 0; i-- {
		block := bc.Blocks[i]
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

			// 2. 先遍历交易输出， 判断交易输出是否在已花费
		IterOutputs:
			for outIdx, output := range tx.Outputs {
				if !output.ToAddressRight(address) {
					continue
				}

				if spentIds[txID] != nil {
					for _, spentOut := range spentIds[txID] {
						if spentOut == outIdx {
							continue IterOutputs
						}
					}
				}
				if output.ToAddressRight(address) {
					unSpentTxs = append(unSpentTxs, tx)
				}
			}

			// 3. 将交易输入添加到已花费的交易中
			if tx.IsBase() { // 第一个区块没有交易输入
				continue
			}

			for _, input := range tx.Inputs {
				if !input.FromAddressRight(address) {
					continue
				}
				inputTxID := hex.EncodeToString(input.TxID)
				if len(spentIds[inputTxID]) == 0 {
					spentIds[inputTxID] = make([]int, 0)
				}
				spentIds[inputTxID] = append(spentIds[inputTxID], input.OutIdx)
			}
		}

	}

	return unSpentTxs
}

func (bc *BlockChain) FindUTXOs(address []byte) (int, map[string]int) {
	unspentOuts := make(map[string]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if out.ToAddressRight(address) {
				accumulated += out.Value
				unspentOuts[txID] = outIdx
				continue Work // one transaction can only have one output referred to adderss
			}
		}
	}
	return accumulated, unspentOuts
}

func (bc *BlockChain) FindSpendableOutputs(address []byte, amount int) (int, map[string]int) {
	unspentOuts := make(map[string]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if out.ToAddressRight(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = outIdx
				if accumulated >= amount {
					break Work
				}
				continue Work // one transaction can only have one output referred to adderss
			}
		}
	}
	return accumulated, unspentOuts
}

// FindUnspentTransactions2 根据帐号找出未使用的交易
func (bc *BlockChain) FindUnspentTransactions2(address []byte, amount int) (int, []*transaction.Transaction) {
	unSpentTxs := make([]*transaction.Transaction, 0)
	spentIds := make(map[string][]int) // 用于保存已经使用的交易输出 string是交易ID， []int交易下标
	value := 0

	// 1 遍历整个区块链,从最后一个区块开始
	for i := len(bc.Blocks); i >= 0; i-- {
		block := bc.Blocks[i]
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

			// 2. 先遍历交易输出， 判断交易输出是否在已花费
		IterOutputs:
			for outIdx, output := range tx.Outputs {
				if !output.ToAddressRight(address) {
					continue
				}

				if spentIds[txID] != nil {
					for _, spentOut := range spentIds[txID] {
						if spentOut == outIdx {
							continue IterOutputs
						}
					}
				}
				if output.ToAddressRight(address) {
					unSpentTxs = append(unSpentTxs, tx)
					value = value + output.Value
					if value >= amount {
						return value, unSpentTxs
					}
				}
			}

			// 3. 将交易输入添加到已花费的交易中
			if tx.IsBase() { // 第一个区块没有交易输入
				continue
			}

			for _, input := range tx.Inputs {
				if !input.FromAddressRight(address) {
					continue
				}
				inputTxID := hex.EncodeToString(input.TxID)
				if len(spentIds[inputTxID]) == 0 {
					spentIds[inputTxID] = make([]int, 0)
				}
				spentIds[inputTxID] = append(spentIds[inputTxID], input.OutIdx)
			}
		}

	}

	return value, unSpentTxs
}

// CreateTransaction 创建交易
func (bc *BlockChain) CreateTransaction(from, to []byte, amount int) (*transaction.Transaction, error) {
	input := make([]transaction.TxInput, 0)
	output := make([]transaction.TxOutput, 0)

	// 获取余额
	value, validOutputs := bc.FindSpendableOutputs(from, amount)
	if value < amount {
		return nil, errors.New("not enough funds")
	}

	// 足够的金额给to, 剩余的给from
	for idx, outIdx := range validOutputs {
		idxByte, err := hex.DecodeString(idx)
		utils.Handle(err)
		input = append(input, transaction.TxInput{
			TxID:        idxByte,
			OutIdx:      outIdx,
			FromAddress: from,
		})
	}

	output = append(output, transaction.TxOutput{
		Value:     amount,
		ToAddress: to,
	})
	if value > amount {
		output = append(output, transaction.TxOutput{
			Value:     value - amount,
			ToAddress: from,
		})
	}

	tx := transaction.Transaction{
		Inputs:  input,
		Outputs: output,
	}
	tx.SetId()
	return &tx, nil
}

func (bc *BlockChain) Mine(txs []*transaction.Transaction) {
	bc.AddBlock(txs)
}
