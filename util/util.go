package util

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
	"math/rand"
)

func GenerateECDSAKey(pkey string) (string, *ecdsa.PrivateKey) {
	privateKey, err := crypto.HexToECDSA(pkey)
	if err != nil {
		log.Println(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	_publicKey := hexutil.Encode(hash.Sum(nil)[12:])
	_privateKey := privateKey

	return _publicKey, _privateKey
}

func RandStringRunes(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}