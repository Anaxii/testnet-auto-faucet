package service

import (
	"crypto/ecdsa"
	"testnet-autofaucet/config"
)

type FaucetService struct {
	PrivateKey    *ecdsa.PrivateKey
	PublicKey     string
	ExternalDBURI string
	Subnet        config.Subnet
	Accounts      map[string]bool
}
