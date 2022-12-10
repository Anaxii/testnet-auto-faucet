package service

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/big"
	"strconv"
	"testnet-autofaucet/api"
	db "testnet-autofaucet/embeddeddatabase"
	"time"
)

func (s *FaucetService) Start() {
	s.fetchAccounts()

	ticker := time.NewTicker(15 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				s.fetchAccounts()
			}
		}
	}()

	go s.listenForAccounts()
	log.Info("Starting faucet service")
	s.faucetService()

}

func (s *FaucetService) faucetService() {

	for {
		for k, v := range s.Accounts {
			if v {
				s.checkAccount(k)
			}
		}
		time.Sleep(time.Second * 15)
	}
}

func (s *FaucetService) checkAccount(walletAddress string) {
	lastTime, err := db.Read([]byte("accounts"), []byte(walletAddress))

	if err != nil {
		db.Write([]byte("accounts"), []byte(walletAddress), []byte("0"))
		lastTime = []byte("0")
	}

	dbTime, err := strconv.Atoi(string(lastTime))
	if err != nil {
		log.WithFields(log.Fields{
			"status": "checkAccount",
			"time":   string(lastTime),
		}).Warn("Could not convert time to int")
		db.Write([]byte("accounts"), []byte(walletAddress), []byte("0"))
	}

	if time.Since(time.Unix(int64(dbTime), 0)) <= time.Hour*24 {
		return
	}

	if s.Balance(walletAddress).Cmp(big.NewInt(1e18)) >= 0 {
		return
	}

	log.WithFields(log.Fields{
		"address": walletAddress,
	}).Info("balance under threshold, preparing to send PFN")

	if err = s.Send(walletAddress); err != nil {
		log.WithFields(log.Fields{
			"network": s.Subnet.Name,
			"err":     err,
		}).Error("Error sending funds")
	}

	api.Log(map[string]interface{}{"status": "topped off", "walletAddress": walletAddress})
	db.Write([]byte("accounts"), []byte(walletAddress), []byte(fmt.Sprintf("%v", time.Now().Unix())))

}
