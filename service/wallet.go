package service

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

func (s *FaucetService) Send(walletAddress string) error {
	client, err := ethclient.Dial(s.Subnet.RpcURL)
	if err != nil {
		return err
	}

	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(s.PublicKey))
	if err != nil {
		return err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	toAddress := common.HexToAddress(walletAddress)
	value := big.NewInt(1500000000000000000)
	gasLimit := uint64(21000)

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), s.PrivateKey)
	if err != nil {
		return err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())

	return nil
}

