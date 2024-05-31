package blockchain

import (
	"bytes"
	"encoding/gob"
	"github.com/limitzhang87/goblockchain/constcoe"
	"github.com/limitzhang87/goblockchain/transaction"
	"github.com/limitzhang87/goblockchain/utils"
	"os"
)

type TransactionPool struct {
	Txs []*transaction.Transaction
}

func (p *TransactionPool) AddTransaction(tx *transaction.Transaction) {
	p.Txs = append(p.Txs, tx)
}

func (p *TransactionPool) SaveFile() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(p)
	utils.Handle(err)
	err = os.WriteFile(constcoe.TransactionPoolFile, buffer.Bytes(), 0644)
	utils.Handle(err)
}

func (p *TransactionPool) LoadFile() error {
	// 文件不存在直接退出
	if !utils.FileExists(constcoe.TransactionPoolFile) {
		return nil
	}

	fileContent, err := os.ReadFile(constcoe.TransactionPoolFile)
	if err != nil {
		return err
	}

	var pool TransactionPool
	decoder := gob.NewDecoder(bytes.NewBuffer(fileContent))
	err = decoder.Decode(&pool)
	if err != nil {
		return err
	}
	p.Txs = pool.Txs
	return nil
}

func CreatePool() *TransactionPool {
	pool := &TransactionPool{}
	err := pool.LoadFile()
	utils.Handle(err)
	return pool
}

func RemovePoolFile() error {
	err := os.Remove(constcoe.TransactionPoolFile)
	return err
}
