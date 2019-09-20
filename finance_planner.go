package main

import (
	"time"
)

const (
	minActionInterval = time.Hour
)

type FinancePlanner struct {
	period              time.Duration
	spendingPerCycle    float64
	costPerAction       float64
	minDurationInterval time.Duration
}

func newFinancePlanner(c *Config) *FinancePlanner {
	planner := &FinancePlanner{
		period:              convertPeriod2Duration(c.CyclePeriod),
		spendingPerCycle:    c.SpendingPerCycle,
		costPerAction:       c.IssueCost,
		minDurationInterval: minActionInterval,
	}
	return planner
}

func (p *FinancePlanner) actionDuration() time.Duration {
	actionCount := p.spendingPerCycle / p.costPerAction
	duration := time.Duration(float64(p.period)/actionCount) * time.Nanosecond
	if duration < p.minDurationInterval {
		return p.minDurationInterval
	}
	return duration
}

func convertPeriod2Duration(period string) time.Duration {
	var duration time.Duration
	switch period {
	case "month":
		duration = time.Duration(7*24*30) * time.Hour
	case "week":
		duration = time.Duration(7*24) * time.Hour
	case "day":
		duration = time.Duration(24) * time.Hour
	case "hour":
		duration = time.Duration(1) * time.Hour
	default:
		duration = time.Duration(7*24) * time.Hour
	}
	return duration
}
