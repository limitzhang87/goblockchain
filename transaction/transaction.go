package transaction

import (
	"bytes"
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
		TxID:        []byte{},
		OutIdx:      -1,
		FromAddress: []byte{},
	}

	output := TxOutput{
		Value:     constcoe.InitCoin,
		ToAddress: toAddress,
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
