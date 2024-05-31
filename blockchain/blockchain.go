package blockchain

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/limitzhang87/goblockchain/constcoe"
	"github.com/limitzhang87/goblockchain/transaction"
	"github.com/limitzhang87/goblockchain/utils"
	"runtime"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

//// CreateBlockChain 创建区块链
//func CreateBlockChain(address []byte) *BlockChain {
//	blockchain := BlockChain{}
//	blockchain.Blocks = append(blockchain.Blocks, GenesisBlock(address))
//	return &blockchain
//}

// InitBlockChain 创建区块链，首次创建并创建数据库
func InitBlockChain(address []byte) *BlockChain {
	var lashHash []byte
	if utils.FileExists(constcoe.BCFile) {
		fmt.Println("blockchain exist")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(constcoe.BCFile)
	opts.Logger = nil

	db, err := badger.Open(opts)
	utils.Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		genesis := GenesisBlock(address)
		fmt.Println("genesis block create")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		utils.Handle(err)
		err = txn.Set([]byte(constcoe.LHKey), genesis.Hash) // current block hash
		utils.Handle(err)
		err = txn.Set([]byte(constcoe.OgPrevHashKey), genesis.PrevHash) // genesis block key
		utils.Handle(err)
		lashHash = genesis.Hash
		return nil
	})
	utils.Handle(err)
	blockchain := BlockChain{lashHash, db}
	return &blockchain
}

// ContinueBlockChain 从数据库中读取区块信息创建区块链
func ContinueBlockChain() *BlockChain {
	// 判断区块链数据库是否存在
	if !utils.FileExists(constcoe.BCFile) {
		fmt.Println("blockchain does not exist, please create a new one")
		runtime.Goexit()
	}

	var lashHash []byte
	opts := badger.DefaultOptions(constcoe.BCFile)
	opts.Logger = nil
	db, err := badger.Open(opts)
	utils.Handle(err)
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(constcoe.LHKey))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			lashHash = val
			return nil
		})
		return err
	})
	utils.Handle(err)
	blockchain := BlockChain{lashHash, db}
	return &blockchain
}

// AddBlock 区块链添加区块
func (bc *BlockChain) AddBlock(block *Block) {
	//newBlock := CreateBlock(bc.Blocks[len(bc.Blocks)-1].Hash, txs)
	//bc.Blocks = append(bc.Blocks, newBlock)

	var lastHash []byte
	// 1. 验证内存中的lastHash是否等于区块链中的lh
	err := bc.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(constcoe.LHKey))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			if !bytes.Equal(val, bc.LastHash) {
				return errors.New("block hash does not match")
			}
			lastHash = val
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})

	// 2. 判断传入的区块是否是上一个区块的hash
	if !bytes.Equal(lastHash, block.PrevHash) {
		fmt.Println("block hash does not match")
		runtime.Goexit()
	}

	// 3. 存储区块
	serialized := block.Serialize()
	err = bc.Database.Update(func(txn *badger.Txn) error {
		curHash := block.Hash
		err := txn.Set(curHash, serialized)
		if err != nil {
			return err
		}
		err = txn.Set([]byte(constcoe.LHKey), curHash)
		if err != nil {
			return err
		}
		bc.LastHash = curHash
		return nil
	})
	utils.Handle(err)
}

type Iterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (bc *BlockChain) Iterator() *Iterator {
	return &Iterator{bc.LastHash, bc.Database}
}

func (bcI *Iterator) Next() *Block {
	lastHash := bcI.CurrentHash
	var block *Block
	err := bcI.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(lastHash)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			block = DeSerializeBlock(val)
			return nil
		})
		return err
	})
	utils.Handle(err)
	bcI.CurrentHash = block.PrevHash
	return block
}

func (bc *BlockChain) BackOgPrevHash() []byte {
	var ogPrevHash []byte
	err := bc.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(constcoe.OgPrevHashKey))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			ogPrevHash = val
			return nil
		})
		return err
	})
	utils.Handle(err)
	return ogPrevHash
}

// FindUnspentTransactions 根据帐号找出未使用的交易
func (bc *BlockChain) FindUnspentTransactions(address []byte) []*transaction.Transaction {
	unSpentTxs := make([]*transaction.Transaction, 0)
	spentIds := make(map[string][]int) // 用于保存已经使用的交易输出 string是交易ID， []int交易下标

	iter := bc.Iterator()
	ogPrevHash := bc.BackOgPrevHash()

	// 1 遍历整个区块链,从最后一个区块开始
	for {
		block := iter.Next()
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
		if bytes.Equal(ogPrevHash, block.PrevHash) {
			break
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
				continue Work // one transaction can only have one output referred to address
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
				continue Work // one transaction can only have one output referred to address
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

	iter := bc.Iterator()
	ogPreHash := bc.BackOgPrevHash()

all:
	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

			// 遍历交易输出，判断是否已经使用过了，没使用过就放到unSpentTxs
		IterOutputs:
			for outIdx, out := range tx.Outputs {
				if out.ToAddressRight(address) {
					for _, spentOut := range spentIds[txID] {
						if spentOut == outIdx {
							continue IterOutputs
						}
					}
					// 走到这里说明该交易输出没被使用
					unSpentTxs = append(unSpentTxs, tx)
					value += out.Value
					if value >= amount {
						break all // 钱够了，直接退出循环
					}
				}
			}

			// 遍历交易输入，将交易输入记录到spentIds中，表示该交易已经使用了
			for _, input := range tx.Inputs {
				if !input.FromAddressRight(address) {
					continue
				}

				inputTxID := hex.EncodeToString(input.TxID)
				if spentIds[inputTxID] == nil {
					spentIds[inputTxID] = make([]int, 0)
				}
				spentIds[inputTxID] = append(spentIds[inputTxID], input.OutIdx)
			}
		}

		if bytes.Equal(ogPreHash, block.PrevHash) {
			break
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
