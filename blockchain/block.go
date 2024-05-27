package blockchain

import (
	"bytes"
	"crypto/sha256"
	"github.com/limitzhang87/goblockchain/transaction"
	"github.com/limitzhang87/goblockchain/utils"
	"time"
)

type Block struct {
	Timestamp    int64
	Hash         []byte
	PrevHash     []byte
	Target       []byte
	Nonce        int64
	Transactions []*transaction.Transaction
}

func GenesisBlock() *Block {
	tx := transaction.BaseTx([]byte("limitZhang"))
	return CreateBlock([]byte{}, []*transaction.Transaction{tx})
}

func CreateBlock(prevHash []byte, txs []*transaction.Transaction) *Block {
	block := &Block{
		Timestamp:    time.Now().Unix(),
		Hash:         []byte{},
		PrevHash:     prevHash,
		Target:       []byte{},
		Nonce:        0,
		Transactions: txs,
	}
	block.Target = block.GetTarget()
	block.Nonce = block.FindNonce()
	block.SetHash()
	return block
}

func (b *Block) BackTransactionSummary() []byte {
	txIDs := make([][]byte, 0, len(b.Transactions))
	for _, tx := range b.Transactions {
		txIDs = append(txIDs, tx.ID)
	}
	summary := bytes.Join(txIDs, []byte{})
	return summary
}

func (b *Block) SetHash() {
	information := bytes.Join([][]byte{
		utils.ToHexInt(b.Timestamp),
		b.PrevHash,
		b.Target,
		utils.ToHexInt(b.Nonce),
		b.BackTransactionSummary(),
	}, []byte{})
	hash := sha256.Sum256(information)
	b.Hash = hash[:]
}
