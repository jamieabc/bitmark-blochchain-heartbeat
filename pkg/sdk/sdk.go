package sdk

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	bitmarkSDK "github.com/bitmark-inc/bitmark-sdk-go"
	"github.com/bitmark-inc/bitmark-sdk-go/account"
	"github.com/bitmark-inc/bitmark-sdk-go/asset"
	"github.com/bitmark-inc/bitmark-sdk-go/bitmark"
	"github.com/bxcodec/faker/v3"
	"golang.org/x/text/language"

	"github.com/jamieabc/bitmark-blochchain-heartbeat/pkg/parser"
)

const (
	networkTimeout              = 10 * time.Second
	issuanceForBlockMiner       = 2
	issuanceMakeBlockchainGoing = 1
	MaximumItemName             = 64
)

type TransitionVerbs struct {
	verb string
	adv  string
}

var (
	verbs = []TransitionVerbs{
		{"send", "to"},
		{"receive", "from"},
		{"buy", "from"},
		{"sell", "to"},
		{"transfer", "to"},
		{"register", ""},
		{"claim", ""},
		{"give", "to"},
		{"我的", ""},
		{"買", ""},
		{"賣", ""},
		{"Mua", ""},
		{"Bán", ""},
		{"Của tôi", ""},
	}
	books []string
)

func TruncateLongString(str string) string {
	if len(str) >= MaximumItemName {
		return str[0:MaximumItemName]
	}
	return str
}

func meaningfulName(item string) string {
	str := TruncateLongString(item)
	transitionVerb := verbs[rand.Intn(len(verbs))]
	if "" != transitionVerb.adv {
		return fmt.Sprintf("%s %s %s %s", transitionVerb.verb, str,
			transitionVerb.adv, faker.Name())
	}
	return fmt.Sprintf("%s %s", transitionVerb.verb, str)
}

func NewSDKConfig(config parser.Config, data []string) *bitmarkSDK.Config {
	httpClient := &http.Client{
		Timeout: networkTimeout,
	}

	books = data

	sdkConfig := &bitmarkSDK.Config{
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

func RestoreAccountFromRecoveryPhrase(strs []string) ([]account.Account, error) {
	var accounts []account.Account
	for _, s := range strs {
		phrases := strings.Split(s, ",")
		trimedPhrases := arrayMap(phrases, strings.Trim)
		acc, err := account.FromRecoveryPhrase(trimedPhrases, language.AmericanEnglish)
		if nil != err {
			return accounts, fmt.Errorf("error recovery acc from phrase: %s", err)
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func registerAsset(owner account.Account) (string, error) {
	title := meaningfulName(books[rand.Intn(len(books))])
	params, err := asset.NewRegistrationParams(
		title,
		map[string]string{
			"owner":     faker.Name(),
			"issueTime": time.Now().String(),
			"author":    faker.Name(),
			"date":      faker.Date(),
			"favorite":  faker.Sentence(),
			"sig":       faker.Password(),
		},
	)
	if nil != err {
		return "", err
	}

	err = params.SetFingerprint([]byte(time.Now().String()))
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
		//issuanceForBlockMiner,
		issuanceMakeBlockchainGoing,
	)
	err := params.Sign(issuer)
	if nil != err {
		return []string{}, err
	}
	return bitmark.Issue(params)
}

func CreateIssuanceFromAccountsRandomly(accounts []account.Account) error {
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
