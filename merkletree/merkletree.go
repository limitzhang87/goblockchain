package merkletree

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"github.com/limitzhang87/goblockchain/transaction"
	"github.com/limitzhang87/goblockchain/utils"
)

type MerkleTree struct {
	Root *MerkleNode
}

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

func CreateMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	temp := &MerkleNode{}

	if left == nil && right == nil {
		temp.Data = data
	} else {
		tmpHash := append(left.Data, right.Data...)
		hash := sha256.Sum256(tmpHash)
		temp.Data = hash[:]
	}

	temp.Left = left
	temp.Right = right
	return temp
}

// CreateMerkleTree 根据交易创建merkleTree
func CreateMerkleTree(txs []*transaction.Transaction) *MerkleTree {
	txLen := len(txs)
	if txLen%2 != 0 {
		txs = append(txs, txs[txLen-1])
	}
	nodeList := make([]*MerkleNode, 0, txLen+1)
	for _, tx := range txs {
		nodeList = append(nodeList, CreateMerkleNode(nil, nil, tx.ID))
	}

	for len(nodeList) > 1 {
		l := len(nodeList)
		tmpNodeList := make([]*MerkleNode, 0, l/2)
		// 如果是单数，将最后一个结点单独放到前面, 如果放在后面，极端情况下最后一个一直向上没有合并做hash
		if l%2 != 0 {
			tmpNodeList = append(tmpNodeList, nodeList[l-1])
		}
		for i := 0; i < l/2; i++ {
			tmpNodeList = append(tmpNodeList, CreateMerkleNode(nodeList[2*i], nodeList[2*i+1], nil))
		}
		nodeList = tmpNodeList
	}
	return &MerkleTree{Root: nodeList[0]}
}

func (mn *MerkleNode) Find(data []byte, route []int, hashRoute [][]byte) (bool, []int, [][]byte) {
	findFlag := false
	if bytes.EqualFold(mn.Data, data) {
		findFlag = true
		return findFlag, route, hashRoute
	} else {
		if mn.Left != nil {
			routeT := append(route, 0)
			hashRouteT := append(hashRoute, mn.Right.Data)
			findFlag, routeT, hashRouteT = mn.Left.Find(data, routeT, hashRouteT)
			if findFlag {
				return findFlag, routeT, hashRouteT
			} else {
				if mn.Right != nil {
					routeT = append(route, 1)
					hashRouteT = append(hashRoute, mn.Left.Data)
					findFlag, routeT, hashRouteT = mn.Right.Find(data, routeT, hashRouteT)
					if findFlag {
						return findFlag, routeT, hashRouteT
					} else {
						return findFlag, route, hashRoute
					}
				}
			}
		} else {
			return findFlag, route, hashRoute
		}
	}
	return findFlag, route, hashRoute
}

// BackValidationRoute ...
func (mt *MerkleTree) BackValidationRoute(txId []byte) ([]int, [][]byte, bool) {
	ok, route, hashRoute := mt.Root.Find(txId, []int{}, [][]byte{})
	return route, hashRoute, ok
}

func SimplePaymentValidation(txId, mtRoutHash []byte, route []int, hashRoute [][]byte) bool {
	routeLen := len(route)
	tempHash := txId

	for i := routeLen - 1; i >= 0; i-- {
		var parentHash []byte
		if route[i] == 0 {
			parentHash = append(tempHash, hashRoute[i]...)
		} else if route[i] == 1 {
			parentHash = append(hashRoute[i], tempHash...)
		} else {
			utils.Handle(errors.New("error in validation route"))
		}
		hash := sha256.Sum256(parentHash)
		tempHash = hash[:]
	}
	return bytes.Equal(mtRoutHash, tempHash)
}
