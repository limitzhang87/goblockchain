package blockchain

import (
	"fmt"
	"github.com/limitzhang87/goblockchain/utils"
)

func (bc *BlockChain) RunMine() {
	pool := CreatePool()
	//err := pool.LoadFile()
	//utils.Handle(err)

	block := CreateBlock(bc.LastHash, pool.Txs)
	if !block.ValidatePoW() {
		fmt.Println("Block has invalid nonce!")
		return
	}
	bc.AddBlock(block)
	err := RemovePoolFile()
	utils.Handle(err)
}
