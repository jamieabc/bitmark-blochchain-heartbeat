package main

import (
	"time"

	"github.com/bitmark-inc/bitmark-sdk-go/account"
)

func runBackground(interval time.Duration, accounts []account.Account) {
	go func(time.Duration) {
		for {
			select {
			case <-time.After(interval):
				createIssuanceFromAccountsRandomly(accounts)
			}
		}
	}(interval)
}
