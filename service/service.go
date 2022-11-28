package service

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/big"
	"strconv"
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
			"time":  string(lastTime),
		}).Warn("Could not convert time to int")
		db.Write([]byte("accounts"), []byte(walletAddress), []byte("0"))
	}

	if time.Since(time.Unix(int64(dbTime), 0)) > time.Hour * 24 {
		bal := s.Balance(walletAddress)
		if bal.Cmp(big.NewInt(1e18)) < 0 {
			log.WithFields(log.Fields{
				"address": walletAddress,
			}).Info("balance under threshold, preparing to send PFN")
			s.Send(walletAddress)
			db.Write([]byte("accounts"), []byte(walletAddress), []byte(fmt.Sprintf("%v", time.Now().Unix())))
		} else {
			//log.WithFields(log.Fields{
			//	"address": walletAddress,
			//}).Info("balance over threshold")
		}
	}

}
