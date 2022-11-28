package service

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"math/big"
)

func (s *FaucetService) Balance(walletAddress string) *big.Int {
	conn, err := ethclient.Dial(s.Subnet.RpcURL)
	if err != nil {
		log.Println("wallet/wallet.go:Balance(): Failed to connect to the Ethereum client:", err)
		log.WithFields(log.Fields{
			"network":   s.Subnet.Name,
			"err": err,
		}).Error("wallet/wallet.go:Balance(): Failed to connect to the Ethereum client:")
		return big.NewInt(0)
	}

	balance, err := conn.BalanceAt(context.Background(), common.HexToAddress(walletAddress), nil)
	if err != nil {
		log.WithFields(log.Fields{
			"network":   s.Subnet.Name,
			"err": err,
		}).Error("wallet/wallet.go:Balance(): Failed to call balance")
		return big.NewInt(0)
	}

	return balance
}

