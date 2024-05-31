package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"github.com/limitzhang87/goblockchain/constcoe"
	"github.com/limitzhang87/goblockchain/utils"
	"os"
)

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	utils.Handle(err)
	// 根据椭圆算法，公钥就是椭圆的两个点
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return *privateKey, publicKey
}

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func NewWallet() *Wallet {
	privateKey, publicKey := NewKeyPair()
	return &Wallet{privateKey, publicKey}
}

// Address 公钥哈希
func (w *Wallet) Address() []byte {
	// 1.公钥转为公钥哈希
	pubKeyHash := utils.PublicKeyHash(w.PublicKey)
	// 2. 公钥哈希转为签到地址
	address := utils.PubHash2Address(pubKeyHash)
	return address
}

func (w *Wallet) Save() {
	filename := constcoe.Wallets + string(w.Address()) + ".wlt"
	//var content bytes.Buffer
	//gob.Register(elliptic.P256())
	//encoder := gob.NewEncoder(&content)
	//err := encoder.Encode(w)
	//utils.Handle(err)
	// err = os.WriteFile(filename, jsonData, 0644)

	privateKeyBytes, err := x509.MarshalECPrivateKey(&w.PrivateKey)
	utils.Handle(err)
	privateKeyFile, err := os.Create(filename)
	utils.Handle(err)
	defer func() {
		_ = privateKeyFile.Close()
	}()
	err = pem.Encode(privateKeyFile, &pem.Block{
		Bytes: privateKeyBytes,
	})
	utils.Handle(err)
}

func LoadWallet(address string) *Wallet {
	filename := constcoe.Wallets + address + ".wlt"

	privateKeyFile, err := os.ReadFile(filename)
	utils.Handle(err)
	pemBlock, _ := pem.Decode(privateKeyFile)
	privateKey, err := x509.ParseECPrivateKey(pemBlock.Bytes)
	utils.Handle(err)
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return &Wallet{
		PrivateKey: *privateKey,
		PublicKey:  publicKey,
	}

	//gob.Register(elliptic.P256())
	//decoder := gob.NewDecoder(bytes.NewReader(content))
	//var wallet Wallet
	//err = decoder.Decode(&wallet)
	//utils.Handle(err)
	//return &wallet
}
