package finance_planner_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jamieabc/bitmark-blochchain-heartbeat/pkg/parser"

	"github.com/jamieabc/bitmark-blochchain-heartbeat/internal/finance_planner"
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
		actual := finance_planner.ConvertPeriod2Duration(s.period)
		if actual != s.expected {
			t.Errorf("error convert %s to duration, expect %d but get %d",
				s.period, s.expected, actual)
		}
	}
}

func TestActionIntervalWhenOverMinPeriod(t *testing.T) {
	c := parser.Config{
		CyclePeriod:       "week",
		Crypto:            "ltc",
		SpendingPerCycle:  0.001,
		MinSpendingPeriod: "hour",
		IssueCost:         0.001,
		TransferCost:      0.002,
		Chain:             "testnet",
		SDKApiToken:       "token",
		RecoveryPhrases:   []string{},
	}
	expected := time.Hour * 24 * 7

	p := finance_planner.NewFinancePlanner(c)
	actual := p.ActionDuration()
	assert.Equal(t, expected, actual, "wrong action interval")
}

func TestActionIntervalWhenBelowMinPeriod(t *testing.T) {
	c := parser.Config{
		CyclePeriod:       "hour",
		Crypto:            "ltc",
		SpendingPerCycle:  1000000,
		MinSpendingPeriod: "hour",
		IssueCost:         0.001,
		TransferCost:      0.002,
		Chain:             "testnet",
		SDKApiToken:       "token",
		RecoveryPhrases:   []string{},
	}

	p := finance_planner.NewFinancePlanner(c)
	actual := p.ActionDuration()
	assert.Equal(t, time.Hour, actual, "wrong action interval")
}
