package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"github.com/limitzhang87/goblockchain/merkletree"
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
	MTree        *merkletree.MerkleTree
}

func GenesisBlock(address []byte) *Block {
	tx := transaction.BaseTx(address)
	genesis := CreateBlock([]byte("limitZhang is awesome!"), []*transaction.Transaction{tx})
	genesis.SetHash()
	return genesis
}

func CreateBlock(prevHash []byte, txs []*transaction.Transaction) *Block {
	block := &Block{
		Timestamp:    time.Now().Unix(),
		Hash:         []byte{},
		PrevHash:     prevHash,
		Target:       []byte{},
		Nonce:        0,
		Transactions: txs,
		MTree:        merkletree.CreateMerkleTree(txs),
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
		b.MTree.Root.Data,
	}, []byte{})
	hash := sha256.Sum256(information)
	b.Hash = hash[:]
}

func (b *Block) Serialize() []byte {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(b)
	utils.Handle(err)
	return buf.Bytes()
}

func DeSerializeBlock(data []byte) *Block {
	block := new(Block) // new 返回一个对象指针，并将指针指向一块内存
	// var block *Block 不能使用这种语法，声明一个指针，结果没有给指针分配内存，导致指针指向为空，后面使用时会报错
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(block)
	utils.Handle(err)
	return block
}
