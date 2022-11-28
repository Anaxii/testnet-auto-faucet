package main

import (
	"testnet-autofaucet/config"
	"testnet-autofaucet/service"
)

func main() {
	privateKey, publicKey, subnet, externalDBURI := config.GetConfig()
	faucetService := service.FaucetService{
		PrivateKey: privateKey,
		PublicKey: publicKey,
		Subnet: subnet,
		ExternalDBURI: externalDBURI,
		Accounts: map[string]bool{},
	}

	faucetService.Start()
}
