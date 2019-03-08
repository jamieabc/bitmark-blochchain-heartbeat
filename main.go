package main

import (
	"fmt"

	sdk "github.com/bitmark-inc/bitmark-sdk-go"
)

func main() {
	config, err := newConfig()
	if nil != err {
		fmt.Printf("parse config: %s", err.Error())
		return
	}

	sdk.Init(newSdkConfig(config))
	account, err := restoreAccountFromRecoveryPhrase(config.RecoveryPhrases)
	if nil != err {
		fmt.Printf("%s", err.Error())
		return
	}
	fmt.Printf("account: %v", account)
}
