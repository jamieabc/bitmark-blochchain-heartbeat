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
		fmt.Printf("restore accoutn error: %s", err.Error())
		return
	}
	assetID, err := registerAsset(account)
	if nil != err {
		fmt.Printf("register asset error: %s", err)
		return
	}
	bitmarkIDs, err := issueAsset(account, assetID)
	if nil != err {
		fmt.Printf("issue asset error: %s", err)
		return
	}
	fmt.Printf("assetID: %s, bitmarkIDs: %v\n", assetID, bitmarkIDs)
}
