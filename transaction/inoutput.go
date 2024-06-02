package transaction

import (
	"bytes"
	"github.com/limitzhang87/goblockchain/utils"
)

// TxInput 交易输入
type TxInput struct {
	TxID   []byte
	OutIdx int
	PubKey []byte // 公钥
	Sig    []byte // 数字签名
}

// TxOutput 交易输出
type TxOutput struct {
	Value      int
	PubKeyHash []byte // 公钥hash
}

func (in *TxInput) FromAddressRight(pubKey []byte) bool {
	return bytes.Equal(in.PubKey, pubKey)
}

func (out *TxOutput) ToAddressRight(pubKey []byte) bool {
	return bytes.Equal(out.PubKeyHash, utils.PublicKeyHash(pubKey))
}
