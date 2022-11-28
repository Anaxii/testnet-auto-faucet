package service

func (s *FaucetService) Start() {
	s.fetchAccounts()
}
