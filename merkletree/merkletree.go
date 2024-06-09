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

	if left != nil && right != nil {
		tmpHash := append(left.Data, right.Data...)
		hash := sha256.Sum256(tmpHash)
		temp.Data = hash[:]
	} else {
		temp.Data = data
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

func (mn *MerkleNode) Find2(data []byte, route []int, hashRoute [][]byte) (bool, []int, [][]byte) {
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

// Find SPV算法中，返回待验证数据的路径和路径hash， 使用深度有限算法
func (mn *MerkleNode) Find(data []byte, route []int, hashRoute [][]byte) (bool, []int, [][]byte) {
	flag := false
	if bytes.Equal(mn.Data, data) {
		flag = true
		return flag, route, hashRoute
	}
	if mn.Left != nil {
		routeT := append(route, 0)                     // 路径中0， 代表当前深度有限算法中，走的是左边
		hashRouteT := append(hashRoute, mn.Right.Data) // 放入右边的数据，为了验证，从底部往上回溯时左边的数据和这个右边的数据合并一起
		flag, routeT, hashRouteT = mn.Left.Find2(data, routeT, hashRouteT)
		if flag {
			return flag, routeT, hashRouteT
		}
	}
	if mn.Right != nil {
		routeT := append(route, 1)                    // 路径中1， 代表当前深度有限算法中，走的是右边
		hashRouteT := append(hashRoute, mn.Left.Data) // 放入左边的结点，这是为了后面验证时，
		flag, routeT, hashRouteT = mn.Right.Find2(data, routeT, hashRouteT)
		if flag {
			return flag, routeT, hashRouteT
		}
	}
	return flag, route, hashRoute
}

// BackValidationRoute ...
func (mt *MerkleTree) BackValidationRoute(txId []byte) ([]int, [][]byte, bool) {
	ok, route, hashRoute := mt.Root.Find(txId, []int{}, [][]byte{})
	return route, hashRoute, ok
}

// SimplePaymentValidation SPV算法验算数据是否是块中数据，从下往上验算
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
