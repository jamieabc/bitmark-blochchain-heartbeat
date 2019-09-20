package main

import (
	"testing"
	"time"
)

func TestConvertPeriod2Duration(t *testing.T) {
	fixture := []struct {
		period   string
		expected time.Duration
	}{
		{"month", time.Duration(7*24*30) * time.Hour},
		{"week", time.Duration(7*24) * time.Hour},
		{"day", time.Duration(24) * time.Hour},
		{"hour", time.Duration(1) * time.Hour},
	}

	for _, s := range fixture {
		actual := convertPeriod2Duration(s.period)
		if actual != s.expected {
			t.Errorf("error convert %s to duration, expect %d but get %d",
				s.period, s.expected, actual)
		}
	}
}

func TestActionInterval(t *testing.T) {
	minDuration := time.Duration(5) * time.Second
	fixture := []struct {
		config   *Config
		expected time.Duration
	}{
		{&Config{
			CyclePeriod:      "week",
			Crypto:           "ltc",
			SpendingPerCycle: 1,
			IssueCost:        0.001,
			TransferCost:     0.002,
			Chain:            "testnet",
			SDKApiToken:      "token",
			RecoveryPhrases:  []string{}},
			time.Duration(7*24*60*60*1000*0.001) * time.Millisecond},
		{&Config{
			CyclePeriod:      "hour",
			Crypto:           "ltc",
			SpendingPerCycle: 1,
			IssueCost:        0.001,
			TransferCost:     0.002,
			Chain:            "testnet",
			SDKApiToken:      "token",
			RecoveryPhrases:  []string{}},
			minDuration},
	}

	for _, s := range fixture {
		p := newFinancePlanner(s.config)
		actual := p.actionDuration()
		if actual != s.expected {
			t.Errorf("error get action interval, expect %d but get %d",
				s.expected, actual)
		}
	}
}
