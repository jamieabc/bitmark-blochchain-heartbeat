package finance_planner

import (
	"time"

	"github.com/jamieabc/bitmark-blochchain-heartbeat/pkg/parser"
)

type FinancePlanner struct {
	period              time.Duration
	spendingPerCycle    float64
	costPerAction       float64
	minDurationInterval time.Duration
}

func NewFinancePlanner(c parser.Config) *FinancePlanner {
	planner := &FinancePlanner{
		period:              ConvertPeriod2Duration(c.CyclePeriod),
		spendingPerCycle:    c.SpendingPerCycle,
		costPerAction:       c.IssueCost,
		minDurationInterval: ConvertPeriod2Duration(c.MinSpendingPeriod),
	}
	return planner
}

func (p *FinancePlanner) ActionDuration() time.Duration {
	actionCount := p.spendingPerCycle / p.costPerAction
	duration := time.Duration(float64(p.period)/actionCount) * time.Nanosecond
	if duration < p.minDurationInterval {
		return p.minDurationInterval
	}
	return duration
}

func ConvertPeriod2Duration(period string) time.Duration {
	var duration time.Duration
	switch period {
	case "month":
		duration = time.Duration(7*24*30) * time.Hour
	case "week":
		duration = time.Duration(7*24) * time.Hour
	case "day":
		duration = time.Duration(24) * time.Hour
	case "hour":
		duration = time.Hour
	case "min":
		duration = time.Minute
	default:
		duration = time.Duration(7*24) * time.Hour
	}
	return duration
}
