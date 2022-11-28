package service

import "testnet-autofaucet/externaldatabase"

func (s *FaucetService) fetchAccounts() {
	_accounts := externaldatabase.GetAccounts(s.ExternalDBURI)
	if _accounts == nil {
		_accounts = map[string]bool{}
	}
	if len(_accounts) != 0 {
		_accounts[s.PublicKey] = false
		s.Accounts = _accounts
	}
}
