package transaction

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/gob"
	"github.com/limitzhang87/goblockchain/constcoe"
	"github.com/limitzhang87/goblockchain/utils"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

func (tx *Transaction) TxHash() []byte {
	var encoded bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	utils.Handle(err)
	hash = sha256.Sum256(encoded.Bytes())
	return hash[:]
}

func (tx *Transaction) SetId() {
	tx.ID = tx.TxHash()
}

func BaseTx(toAddress []byte) *Transaction {
	input := TxInput{
		TxID:   []byte{},
		OutIdx: -1,
		PubKey: []byte{},
	}

	output := TxOutput{
		Value:      constcoe.InitCoin,
		PubKeyHash: toAddress,
	}
	tx := Transaction{
		ID:      []byte("This is the Base Transaction!"),
		Inputs:  []TxInput{input},
		Outputs: []TxOutput{output},
	}
	return &tx
}

func (tx *Transaction) IsBase() bool {
	return len(tx.Inputs) == 1 && tx.Inputs[0].OutIdx == -1
}

// Sign 对交易进行签名
func (tx *Transaction) Sign(priKey ecdsa.PrivateKey) {
	if tx.IsBase() {
		return
	}
	/*
		对一个交易进行签名时，需要分别针对每一笔交易输入进行签名，
		但是为了防止交易输入被移除或者增加，每一笔交易输入签名时，其签名数据也包含了全部的交易输入，
		只不过其他的交易输入只需要包含索引和金额，不需要公钥
	*/
	for idx, in := range tx.Inputs {
		plainHash := tx.PlainHash(idx, in.PubKey)
		signature := utils.Sign(plainHash, priKey)
		tx.Inputs[idx].Sig = signature
	}
}

// PlainHash 加密前数据hash
func (tx *Transaction) PlainHash(inIdx int, prevPubKey []byte) []byte {
	txCopy := tx.PlainCopy()
	txCopy.Inputs[inIdx].PubKey = prevPubKey
	return txCopy.TxHash()
}

// PlainCopy 复制交易数据，交易输入中的公钥需要排除掉
func (tx *Transaction) PlainCopy() *Transaction {
	input := make([]TxInput, 0, len(tx.Inputs))
	output := make([]TxOutput, 0, len(tx.Outputs))

	for _, in := range tx.Inputs {
		input = append(input, TxInput{TxID: in.TxID, OutIdx: in.OutIdx})
	}

	for _, out := range tx.Outputs {
		output = append(output, TxOutput{Value: out.Value, PubKeyHash: out.PubKeyHash})
	}
	return &Transaction{
		ID:      tx.ID,
		Inputs:  input,
		Outputs: output,
	}
}

// Verity 验证签名
func (tx *Transaction) Verity() bool {
	for i, in := range tx.Inputs {
		txHash := tx.PlainHash(i, in.PubKey)
		ok := utils.Verity(txHash, in.PubKey, in.Sig)
		if !ok {
			return false
		}
	}
	return true
}
