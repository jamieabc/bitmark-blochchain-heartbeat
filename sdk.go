package main

import (
	"fmt"
	"net/http"
	"time"

	sdk "github.com/bitmark-inc/bitmark-sdk-go"
	"github.com/bitmark-inc/bitmark-sdk-go/account"
	"github.com/bitmark-inc/bitmark-sdk-go/asset"
	"github.com/bitmark-inc/bitmark-sdk-go/bitmark"
	"github.com/bxcodec/faker/v3"
	"golang.org/x/text/language"
)

const (
	networkTimeout = 10 * time.Second
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

func restoreAccountFromRecoveryPhrase(strs []string) (account.Account, error) {
	account, err := account.FromRecoveryPhrase(strs, language.AmericanEnglish)
	if nil != err {
		return nil, fmt.Errorf("error recovery account from phrase: %s", err)
	}
	return account, nil
}

func registerAsset(owner account.Account) (string, error) {
	name := faker.Username()
	params, err := asset.NewRegistrationParams(
		name,
		map[string]string{"owner": name},
	)
	if nil != err {
		return "", err
	}

	fakeData := faker.Email()
	err = params.SetFingerprint([]byte(fakeData))
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
