package config

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/big"
	"os"
	"testnet-autofaucet/util"
)

type ConfigStruct struct {
	PrivateKey string `json:"private_key"`
	ExternalDBURI string `json:"external_db_uri"`
	Subnet     Subnet `json:"subnet"`
}

type Subnet struct {
	Name             string   `json:"string"`
	RpcURL           string   `json:"rpc_url"`
	WSURL            string   `json:"ws_url"`
	ChainId          *big.Int `json:"chain_id"`
	KYCAddress       string   `json:"kyc_address"`
	BridgeAddress    string   `json:"bridge_address"`
	BlockRequirement int      `json:"block_requirement"`
}

func GetConfig() (*ecdsa.PrivateKey, string, Subnet, string) {
	jsonFile, err := os.Open("config.json")
	if err != nil {

		privateKey, err := crypto.GenerateKey()
		if err != nil {
			log.Fatal(err)
		}
		privateKeyBytes := crypto.FromECDSA(privateKey)

		file, _ := json.MarshalIndent(ConfigStruct{
			PrivateKey: fmt.Sprintf("%v", hexutil.Encode(privateKeyBytes)[2:]),
			Subnet: Subnet{
				Name:             "puffin",
				RpcURL:           "https://node.thepuffin.network/ext/bc/273dwzFtrR6JQzLncTAbN5RBtiqdysVfKTJKBvYHhtUHBnrYWe/rpc",
				WSURL:            "ws://52.35.42.217:9650/ext/bc/273dwzFtrR6JQzLncTAbN5RBtiqdysVfKTJKBvYHhtUHBnrYWe/ws",
				ChainId:          big.NewInt(43113114),
				KYCAddress:       "0x0200000000000000000000000000000000000002",
				BridgeAddress:    "0xd3E11DeF6d34E231ab410e5aA187e1f2d9fF19E1",
				BlockRequirement: 0,
			},
		}, "", "  ")
		_ = ioutil.WriteFile("config.json", file, 0644)
		jsonFile, err = os.Open("config.json")
		log.Fatal("Edit mongodb URI")
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal("Config JSON invalid", err)
	}

	var config ConfigStruct
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		log.Fatal("Could not parse config", err)
	}

	key := os.Getenv("private_key")
	if key != "" {
		config.PrivateKey = key
	}

	log.Info("Config initialized")
	_publicKey, _privateKey := util.GenerateECDSAKey(config.PrivateKey)

	return _privateKey, _publicKey, config.Subnet, config.ExternalDBURI

}
