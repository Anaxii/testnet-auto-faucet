package service

import (
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
	select {

	}

}
