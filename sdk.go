package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	sdk "github.com/bitmark-inc/bitmark-sdk-go"
	"github.com/bitmark-inc/bitmark-sdk-go/account"
	"github.com/bitmark-inc/bitmark-sdk-go/asset"
	"github.com/bitmark-inc/bitmark-sdk-go/bitmark"
	"golang.org/x/text/language"
)

const (
	networkTimeout = 10 * time.Second
	fingerprint    = "heartbeat@bitmark"
)

func newSdkConfig(config *Config) *sdk.Config {
	httpClient := &http.Client{
		Timeout: networkTimeout,
	}

	sdkConfig := &sdk.Config{
		APIToken:   config.SDKApiToken,
		Network:    config.Chain,
		HTTPClient: httpClient,
	}
	return sdkConfig
}

func arrayMap(strs []string, f func(string, string) string) []string {
	newStrs := make([]string, len(strs))
	for i, s := range strs {
		newStrs[i] = f(s, " ")
	}
	return newStrs
}

func restoreAccountFromRecoveryPhrase(strs []string) ([]account.Account, error) {
	var accounts []account.Account
	for _, s := range strs {
		phrases := strings.Split(s, ",")
		trimesPhrases := arrayMap(phrases, strings.Trim)
		account, err := account.FromRecoveryPhrase(trimesPhrases, language.AmericanEnglish)
		if nil != err {
			return accounts, fmt.Errorf("error recovery account from phrase: %s", err)
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func registerAsset(owner account.Account) (string, error) {
	name := "heartbeat"
	params, err := asset.NewRegistrationParams(
		name,
		map[string]string{"owner": name},
	)
	if nil != err {
		return "", err
	}

	err = params.SetFingerprint([]byte(fingerprint))
	if nil != err {
		return "", err
	}

	err = params.Sign(owner)
	if nil != err {
		return "", err
	}

	return asset.Register(params)
}

func issueAsset(issuer account.Account, assetID string) ([]string, error) {
	params := bitmark.NewIssuanceParams(
		assetID,
		1,
	)
	params.Sign(issuer)
	return bitmark.Issue(params)
}

func createIssuanceFromAccountsRandomly(accounts []account.Account) error {
	issuer := randomPickUser(accounts)
	fmt.Printf("%v: %s create issuance\n",
		time.Now(), issuer.AccountNumber())
	assetID, err := registerAsset(issuer)
	if nil != err {
		fmt.Printf("register asset error: %s", err)
		return err
	}
	fmt.Printf("assetID: %s\n", assetID)

	bitmarkIDs, err := issueAsset(issuer, assetID)
	if nil != err {
		fmt.Printf("issue asset error: %s", err)
		return nil
	}
	fmt.Printf("bitmark IDs: %v\n\n", bitmarkIDs)
	return nil
}

func randomPickUser(accounts []account.Account) account.Account {
	if 1 == len(accounts) {
		return accounts[0]
	}
	randomIndex := rand.Intn(len(accounts))
	return accounts[randomIndex]
}
