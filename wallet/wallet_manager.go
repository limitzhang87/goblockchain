package wallet

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/limitzhang87/goblockchain/constcoe"
	"github.com/limitzhang87/goblockchain/utils"
	"os"
	"path/filepath"
	"strings"
)

type RefList map[string]string

func (r *RefList) Save() {
	filename := constcoe.WalletsRefList + "ref_list.data"
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(r)
	utils.Handle(err)
	err = os.WriteFile(filename, buffer.Bytes(), 0644)
	utils.Handle(err)
}

// Update 用于遍历钱包存在refList
func (r *RefList) Update() {
	err := filepath.Walk(constcoe.Wallets, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		filename := f.Name()
		if strings.Compare(filename[len(filename)-4:], ".wlt") == 0 {
			_, ok := (*r)[filename[:len(filename)-4]]
			if !ok {
				(*r)[filename[:len(filename)-4]] = ""
			}
		}
		return nil
	})
	utils.Handle(err)
}

// BindRef 别名绑定地址
func (r *RefList) BindRef(address, refName string) {
	(*r)[address] = refName
}

// FindRef 通过别名找地址
func (r *RefList) FindRef(refName string) (string, error) {
	tmp := ""
	for k, v := range *r {
		if v == refName {
			tmp = k
		}
	}
	if tmp == "" {
		return "", errors.New("refName not found")
	}
	return tmp, nil

}

// LoadRefList 加载地址与别名
func LoadRefList() *RefList {
	filename := constcoe.WalletsRefList + "ref_list.data"
	var refList RefList
	if utils.FileExists(filename) {
		fileContent, err := os.ReadFile(filename)
		utils.Handle(err)
		decoder := gob.NewDecoder(bytes.NewBuffer(fileContent))
		err = decoder.Decode(&refList)
		utils.Handle(err)
	} else {
		refList = make(RefList)
		refList.Update()
	}
	return &refList
}
