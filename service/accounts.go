package service

import (
	log "github.com/sirupsen/logrus"
	"testnet-autofaucet/externaldatabase"
)

func (s *FaucetService) fetchAccounts() {
	_accounts := externaldatabase.GetAccounts(s.ExternalDBURI)
	if _accounts == nil {
		_accounts = map[string]bool{}
	}
	if len(_accounts) != 0 {
		_accounts[s.PublicKey] = false
		s.Accounts = _accounts
	}
	log.Info("Refreshed wallet list from DB")
}

func (s *FaucetService) listenForAccounts() {
	accountChannel := make(chan string)
	go externaldatabase.ListenForChanges(s.ExternalDBURI, accountChannel)
	for {
		select {
			case k := <- accountChannel:
				s.Accounts[k] = true
				log.WithFields(log.Fields{
					"address": k,
				}).Info("New wallet added to faucet list")
		}
	}
}
