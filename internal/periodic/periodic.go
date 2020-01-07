package periodic

import (
	"fmt"
	"time"

	"github.com/bitmark-inc/bitmark-sdk-go/account"
	"github.com/jamieabc/bitmarkd-broadcast-monitor/nodes/node"

	"github.com/jamieabc/bitmark-blochchain-heartbeat/block"
	"github.com/jamieabc/bitmark-blochchain-heartbeat/pkg/parser"
	"github.com/jamieabc/bitmark-blochchain-heartbeat/pkg/sdk"
)

const (
	blockTime = 3 * time.Minute
	retryTime = 30 * time.Second
)

type Periodic interface {
	Do()
}

func NewPeriodic(duration time.Duration, accounts []account.Account, shutdownChan chan struct{}, config parser.Config) Periodic {
	return &periodic{
		duration:     duration,
		accounts:     accounts,
		shutdownChan: shutdownChan,
		config:       config,
	}
}

type periodic struct {
	duration     time.Duration
	accounts     []account.Account
	shutdownChan chan struct{}
	config       parser.Config
}

func (p *periodic) Do() {
	node.Initialise(p.shutdownChan)

	n, err := node.NewNode(
		p.config.NodeConfig,
		p.config.Keys,
		0,
		60,
	)
	if nil != err {
		fmt.Printf("new node with error: %s\n", err)
		return
	}
	timer := time.NewTimer(p.duration)
	updateTimer(timer, n.Remote(), p.duration)

	go func(time.Duration) {
		for {
			select {
			case <-timer.C:
				durationToNextCheck, err := block.DurationToNextCheck(n.Remote(), p.duration)
				if nil != err {
					fmt.Printf("get next check time with error: %s\n", err)
					timer.Reset(time.Hour)
					continue
				}

				if durationToNextCheck >= p.duration {
					err = sdk.CreateIssuanceFromAccountsRandomly(p.accounts)
					if err != nil {
						// error happens, retry in 1 minute
						updateTimer(timer, n.Remote(), retryTime)
					} else {
						// send command success, check for an average block time
						updateTimer(timer, n.Remote(), blockTime)
					}
				} else {
					updateTimer(timer, n.Remote(), p.duration)
				}

			case <-p.shutdownChan:
				fmt.Println("receive shutdown signal, terminate periodic tasks")
				return
			}
		}
	}(p.duration)
}

func updateTimer(timer *time.Timer, remote node.Remote, targetDuration time.Duration) {
	durationToNextCheck, err := block.DurationToNextCheck(remote, targetDuration)
	if nil != err {
		fmt.Printf("get next check time with error: %s\n", err)
		timer.Reset(time.Hour)
		return
	}

	// add a one minute buffer to avoid some check is not consistent for 1 hour
	durationToNextCheck = durationToNextCheck + time.Minute

	fmt.Printf("duration to next check: %s\n", durationToNextCheck)
	timer.Reset(durationToNextCheck)
}
