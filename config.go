package main

import (
	"fmt"
	"path/filepath"

	sdk "github.com/bitmark-inc/bitmark-sdk-go"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
)

const (
	fileName                = "sdk.conf"
	recovery_phrases_length = 12
)

var (
	cyclePeriod = []string{"hour", "day", "week", "month"}
	crypto      = []string{"ltc"}
	chain       = []string{"bitmark", "testnet"}
)

type Config struct {
	CyclePeriod      string      `gluamapper:"cycle_period" json:"cycle_period"`
	Crypto           string      `gluamapper:"crypto" json:"crypto"`
	SpendingPerCycle float64     `gluamapper:"spending_in_cycle" json:"spending_in_cycle"`
	Chain            sdk.Network `gluamapper:"chain" json:"chain"`
	SDKApiToken      string      `gluamapper:"sdk_api_token" json:"sdk_api_token"`
	RecoveryPhrases  []string    `gluamapper:"recovery_phrases" json:"recovery_phrases"`
}

func newConfig() (*Config, error) {
	path, err := filepath.Abs(filepath.Clean(fileName))
	if nil != err {
		return nil, err
	}

	config := &Config{
		CyclePeriod:      "week",
		Crypto:           "ltc",
		SpendingPerCycle: 0.01,
		Chain:            "testing",
		SDKApiToken:      "",
	}

	_ = config.parse(path)

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
		!contains(crypto, c.Crypto) ||
		!contains(chain, string(c.Chain)) {
		return false
	}
	if len(c.RecoveryPhrases) != recovery_phrases_length {
		return false
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
