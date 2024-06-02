package utils

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"github.com/limitzhang87/goblockchain/constcoe"
	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
	"log"
	"math/big"
	"os"
)

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToHexInt(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func FileExists(fileAddr string) bool {
	if _, err := os.Stat(fileAddr); os.IsNotExist(err) {
		return false
	}
	return true
}

// PublicKeyHash 生成公钥哈希
func PublicKeyHash(public []byte) []byte {
	hashedPublicKey := sha256.Sum256(public)
	hasher := ripemd160.New()
	_, err := hasher.Write(hashedPublicKey[:])
	Handle(err)
	publicRipeMd := hasher.Sum(nil)
	return publicRipeMd
}

func Checksum(ripeMdHash []byte) []byte {
	hash1 := sha256.Sum256(ripeMdHash)
	hash2 := sha256.Sum256(hash1[:])
	return hash2[:constcoe.ChecksumLength]
}

func Base58Encode(input []byte) []byte {
	return []byte(base58.Encode(input))
}

func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	Handle(err)
	return decode
}

// PubHash2Address 公钥哈希转为签到地址
func PubHash2Address(publicKeyHash []byte) []byte {
	networkVersionHash := append([]byte{constcoe.NetworkVersion}, publicKeyHash...)
	checksum := Checksum(networkVersionHash)
	finalHash := append(networkVersionHash, checksum...)
	return Base58Encode(finalHash)
}

// Address2PubHash 地址转为公钥哈希
func Address2PubHash(address []byte) []byte {
	decodeHash := Base58Decode(address)
	pubKeyHash := decodeHash[1 : len(decodeHash)-constcoe.ChecksumLength]
	return pubKeyHash
}

// Sign 根据私钥对数据进行签名
func Sign(msg []byte, privateKey ecdsa.PrivateKey) []byte {
	r, s, err := ecdsa.Sign(rand.Reader, &privateKey, msg)
	Handle(err)
	signature := append(r.Bytes(), s.Bytes()...)
	return signature
}

// Verity 验证签名
func Verity(msg []byte, pubKey []byte, signature []byte) bool {
	curve := elliptic.P256()
	r := big.Int{}
	s := big.Int{}
	sigLen := len(signature)

	r.SetBytes(signature[:sigLen/2])
	s.SetBytes(signature[sigLen/2:])

	x, y := big.Int{}, big.Int{}
	keyLen := len(pubKey)
	x.SetBytes(pubKey[:keyLen/2])
	y.SetBytes(pubKey[keyLen/2:])

	rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
	return ecdsa.Verify(&rawPubKey, msg, &r, &s)
}
