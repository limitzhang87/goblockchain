package transaction

import "bytes"

type TxOutput struct {
	Value     int
	ToAddress []byte
}

type TxInput struct {
	TxID        []byte
	OutIdx      int
	FromAddress []byte
}

func (in *TxInput) FromAddressRight(address []byte) bool {
	return bytes.Equal(address, in.FromAddress)
}

func (out *TxOutput) ToAddressRight(address []byte) bool {
	return bytes.Equal(address, out.ToAddress)
}
