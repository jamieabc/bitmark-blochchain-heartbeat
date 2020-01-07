package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bitmark-inc/logger"

	sdk "github.com/bitmark-inc/bitmark-sdk-go"
	"github.com/jamieabc/bitmarkd-broadcast-monitor/configuration"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
)

const (
	fileName                 = "sdk.conf"
	recoveryPhrasesLength    = 12
	defaultIssueCost         = 0.001
	defaultTransferCost      = 0.002
	defaultMinSpendingPeriod = "hour"
)

var (
	cyclePeriod    = []string{"min", "hour", "day", "week", "month"}
	crypto         = []string{"ltc"}
	chain          = []string{"livenet", "testnet"}
	defaultLogging = logger.Configuration{
		Count:     100,
		Console:   false,
		Directory: "log",
		File:      "heartbeat.log",
		Levels: map[string]string{
			logger.DefaultTag: "error",
		},
		Size: 1048576,
	}
)

type Config struct {
	CyclePeriod       string                   `gluamapper:"cycle_period" json:"cycle_period"`
	Crypto            string                   `gluamapper:"crypto" json:"crypto"`
	SpendingPerCycle  float64                  `gluamapper:"spending_per_cycle" json:"spending_per_cycle"`
	MinSpendingPeriod string                   `gluamapper:"min_spending_period"`
	IssueCost         float64                  `gluamapper:"issue_cost" json:"issue_cost"`
	TransferCost      float64                  `gluamapper:"transfer_cost" json:"transfer_cost"`
	Chain             sdk.Network              `gluamapper:"chain" json:"chain"`
	SDKApiToken       string                   `gluamapper:"sdk_api_token" json:"sdk_api_token"`
	RecoveryPhrases   []string                 `gluamapper:"recovery_phrases" json:"recovery_phrases"`
	NodeConfig        configuration.NodeConfig `gluamapper:"node" json:"node"`
	Keys              configuration.Keys       `gluamapper:"keys"`
	Logging           logger.Configuration     `gluamapper:"logging"`
}

func newConfig() (*Config, error) {
	path, err := filepath.Abs(filepath.Clean(fileName))
	if nil != err {
		return nil, err
	}

	config := &Config{
		CyclePeriod:       "week",
		Crypto:            "ltc",
		SpendingPerCycle:  0.01,
		MinSpendingPeriod: defaultMinSpendingPeriod,
		IssueCost:         defaultIssueCost,
		TransferCost:      defaultTransferCost,
		Chain:             "testing",
		SDKApiToken:       "",
		Logging:           defaultLogging,
	}

	err = config.parse(path)
	if nil != err {
		return nil, err
	}

	if !config.valid() {
		return nil, fmt.Errorf("error format %v\n", config)
	}

	return config, nil
}

func contains(strs []string, str string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}

func (c *Config) valid() bool {
	if !contains(cyclePeriod, c.CyclePeriod) ||
		!contains(cyclePeriod, c.MinSpendingPeriod) ||
		!contains(crypto, c.Crypto) ||
		!contains(chain, string(c.Chain)) {
		return false
	}
	for _, s := range c.RecoveryPhrases {
		if len(strings.Split(s, ",")) != recoveryPhrasesLength {
			return false
		}
	}
	return true
}

func (c *Config) parse(path string) error {
	L := lua.NewState()
	defer L.Close()

	L.OpenLibs()

	arg := &lua.LTable{}
	arg.Insert(0, lua.LString(path))
	L.SetGlobal("arg", arg)

	if err := L.DoFile(path); nil != err {
		return err
	}

	mapperOption := gluamapper.Option{
		NameFunc: func(s string) string {
			return s
		},
		TagName: "gluamapper",
	}

	mapper := gluamapper.Mapper{Option: mapperOption}
	if err := mapper.Map(L.Get(L.GetTop()).(*lua.LTable), c); nil != err {
		return err
	}
	return nil
}
