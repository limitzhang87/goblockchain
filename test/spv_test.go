package test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/limitzhang87/goblockchain/blockchain"
	"github.com/limitzhang87/goblockchain/merkletree"
	"github.com/limitzhang87/goblockchain/transaction"
	"strconv"
	"strings"
	"testing"
)

func GenerateTransaction(outCash int, inAccount string, toAccount string, prevTxID string, outIdx int) *transaction.Transaction {
	prevTxIDHash := sha256.Sum256([]byte(prevTxID))
	inAccountHash := sha256.Sum256([]byte(inAccount))
	toAccountHash := sha256.Sum256([]byte(toAccount))
	txIn := transaction.TxInput{TxID: prevTxIDHash[:], OutIdx: outIdx, PubKey: inAccountHash[:]}
	txOut := transaction.TxOutput{Value: outCash, PubKeyHash: toAccountHash[:]}
	tx := transaction.Transaction{ID: []byte("This is the Base Transaction!"),
		Inputs: []transaction.TxInput{txIn}, Outputs: []transaction.TxOutput{txOut}} //Whether set ID is not necessary
	tx.SetId() //Here the ID is reset to the hash of the whole transaction. Signature is skipped
	return &tx
}

var transactionTests = []struct {
	outCash   int
	inAccount string
	toAccount string
	prevTxID  string
	outIdx    int
}{
	{
		outCash:   10,
		inAccount: "LLL",
		toAccount: "CCC",
		prevTxID:  "prev1",
		outIdx:    1,
	},
	{
		outCash:   20,
		inAccount: "EEE",
		toAccount: "OOO",
		prevTxID:  "prev2",
		outIdx:    1,
	},
	{
		outCash:   30,
		inAccount: "OOO",
		toAccount: "EEE",
		prevTxID:  "prev3",
		outIdx:    0,
	},
	{
		outCash:   100,
		inAccount: "CCC",
		toAccount: "LLL",
		prevTxID:  "prev4",
		outIdx:    1,
	},
	{
		outCash:   50,
		inAccount: "AAA",
		toAccount: "OOO",
		prevTxID:  "prev5",
		outIdx:    1,
	},
	{
		outCash:   110,
		inAccount: "OOO",
		toAccount: "AAA",
		prevTxID:  "prev6",
		outIdx:    0,
	},
	{
		outCash:   200,
		inAccount: "LLL",
		toAccount: "CCC",
		prevTxID:  "prev7",
		outIdx:    1,
	},
	{
		outCash:   500,
		inAccount: "EEE",
		toAccount: "OOO",
		prevTxID:  "prev8",
		outIdx:    1,
	},
}

func GenerateBlock(txs []*transaction.Transaction, prevBlock string) *blockchain.Block {
	prevBlockHash := sha256.Sum256([]byte(prevBlock))
	testBlock := blockchain.CreateBlock(prevBlockHash[:], txs)
	return testBlock
}

var spvTests = []struct {
	txContained []int
	prevBlock   string
	findTX      []int
	wants       []bool
}{
	{
		txContained: []int{0},
		prevBlock:   "prev1",
		findTX:      []int{0, 1},
		wants:       []bool{true, false},
	},
	{
		txContained: []int{0, 1, 2, 3, 4, 5, 6, 7},
		prevBlock:   "prev2",
		findTX:      []int{3, 7, 5},
		wants:       []bool{true, true, true},
	},
	{
		txContained: []int{0, 1, 2, 3},
		prevBlock:   "prev3",
		findTX:      []int{0, 1, 5},
		wants:       []bool{true, true, false},
	},
	{
		txContained: []int{0, 3, 5, 6, 7},
		prevBlock:   "prev4",
		findTX:      []int{0, 1, 6, 7},
		wants:       []bool{true, false, true, true},
	},
	{
		txContained: []int{0, 1, 2, 4, 5, 6, 7},
		prevBlock:   "prev5",
		findTX:      []int{0, 1, 3},
		wants:       []bool{true, true, false},
	},
}

func TestSPV(t *testing.T) {
	var primeTXs []*transaction.Transaction
	for _, tx := range transactionTests {
		tx := GenerateTransaction(tx.outCash, tx.inAccount, tx.toAccount, tx.prevTxID, tx.outIdx)
		primeTXs = append(primeTXs, tx)
	}

	fmt.Println("TestSPV Begin...")
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	for idx, test := range spvTests {
		fmt.Println("Current test No: ", idx)
		fmt.Println("Merkle Tree is like:")
		mtGraphPaint(test.txContained)
		var txs []*transaction.Transaction
		for _, txIdx := range test.txContained {
			txs = append(txs, primeTXs[txIdx])
		}
		testBlock := GenerateBlock(txs, test.prevBlock)
		fmt.Println("------------------------------------------------------------------")
		for num, findIdx := range test.findTX {
			fmt.Println("Find transaction:", findIdx)
			fmt.Printf("Transaction ID: %x\n", primeTXs[findIdx].ID)
			route, hashRoute, ok := testBlock.MTree.BackValidationRoute(primeTXs[findIdx].ID)
			fmt.Println(route, changeByteToStr(hashRoute))
			if ok {
				fmt.Println("Route is like:")
				routePaint(route)
			} else {
				fmt.Println("Has not found the referred transaction")
			}
			spvRes := merkletree.SimplePaymentValidation(primeTXs[findIdx].ID, testBlock.MTree.Root.Data, route, hashRoute)
			fmt.Println("SPV result: ", spvRes, ", Want result: ", test.wants[num])
			if spvRes != test.wants[num] {
				t.Errorf("test %d find %d: SPV is not right", idx, findIdx)
			}
			fmt.Println("------------------------------------------------------------------")
		}
		fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	}
}

func mtGraphPaint(txContained []int) {
	mtLayer := make([][]string, 0)
	bottomLayer := make([]string, 0)
	for i := 0; i < len(txContained); i++ {
		bottomLayer = append(bottomLayer, strconv.Itoa(txContained[i]))
	}
	if len(bottomLayer)%2 == 1 {
		bottomLayer = append(bottomLayer, bottomLayer[len(bottomLayer)-1])
	}
	mtLayer = append(mtLayer, bottomLayer)

	for len(mtLayer[len(mtLayer)-1]) != 1 {
		tempLayer := make([]string, 0)
		if len(mtLayer[len(mtLayer)-1])%2 == 1 {
			tempLayer = append(tempLayer, mtLayer[len(mtLayer)-1][len(mtLayer[len(mtLayer)-1])-1])
			mtLayer[len(mtLayer)-1][len(mtLayer[len(mtLayer)-1])-1] = "->"
		}
		for i := 0; i < len(mtLayer[len(mtLayer)-1])/2; i++ {
			tempLayer = append(tempLayer, mtLayer[len(mtLayer)-1][2*i]+mtLayer[len(mtLayer)-1][2*i+1])
		}

		mtLayer = append(mtLayer, tempLayer)
	}

	layers := len(mtLayer)
	fmt.Println(strings.Repeat(" ", layers-1), mtLayer[layers-1][0])
	foreSpace := 0
	for i := layers - 2; i >= 0; i-- {
		var str1, str2 string
		str1 += strings.Repeat(" ", foreSpace)
		str2 += strings.Repeat(" ", foreSpace)

		for j := 0; j < len(mtLayer[i]); j++ {
			str1 += strings.Repeat(" ", i+1)
			if j%2 == 0 {
				if mtLayer[i][j] == "->" {
					foreSpace += (i+1)*2 + 1
					str1 = strings.Repeat(" ", foreSpace) + str1
					str2 = strings.Repeat(" ", foreSpace) + str2
				} else {
					str1 += "/"
				}

			} else {
				str1 += "\\"
			}
			str1 += strings.Repeat(" ", len(mtLayer[i][j])-1)
			str2 += strings.Repeat(" ", i+1)
			str2 += mtLayer[i][j]
		}
		fmt.Println(str1)
		fmt.Println(str2)
	}

}

func routePaint(route []int) {
	probe := len(route)
	fmt.Println(strings.Repeat(" ", probe) + "*")
	for i := 0; i < len(route); i++ {
		var str1, str2 string
		str1 += strings.Repeat(" ", probe)
		if route[i] == 0 {
			str1 += "/"
			probe -= 1
		} else {
			str1 += "\\"
			probe += 1
		}
		str2 += strings.Repeat(" ", probe) + "*"
		fmt.Println(str1)
		fmt.Println(str2)
	}
}

func changeByteToStr(data [][]byte) []string {
	l := make([]string, len(data))
	for i, datum := range data {
		l[i] = hex.EncodeToString(datum)
	}
	return l
}
