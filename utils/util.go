package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"github.com/limitzhang87/goblockchain/constcoe"
	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
	"log"
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

func Address2PubHash(address []byte) []byte {
	decodeHash := Base58Decode(address)
	pubKeyHash := decodeHash[1 : len(decodeHash)-constcoe.ChecksumLength]
	return pubKeyHash
}
