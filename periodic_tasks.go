package main

import (
	"fmt"
	"time"

	"github.com/bitmark-inc/bitmark-sdk-go/account"
	"github.com/jamieabc/bitmark-blochchain-heartbeat/block"
	"github.com/jamieabc/bitmarkd-broadcast-monitor/nodes/node"
)

const (
	blockTime = 3 * time.Minute
	retryTime = 30 * time.Second
)

type taskInfo struct {
	duration     time.Duration
	accounts     []account.Account
	shutdownChan chan struct{}
	config       *Config
}

func doPeriodicTasks(t taskInfo) {
	node.Initialise(t.shutdownChan)

	n, err := node.NewNode(
		t.config.NodeConfig,
		t.config.Keys,
		0,
		60,
	)
	if nil != err {
		fmt.Printf("new node with error: %s\n", err)
		return
	}
	timer := time.NewTimer(t.duration)
	updateTimer(timer, n.Remote(), t.duration)

	go func(time.Duration) {
		for {
			select {
			case <-timer.C:
				durationToNextCheck, err := block.DurationToNextCheck(n.Remote(), t.duration)
				if nil != err {
					fmt.Printf("get next check time with error: %s\n", err)
					timer.Reset(time.Hour)
					continue
				}

				if durationToNextCheck >= t.duration {
					err = createIssuanceFromAccountsRandomly(t.accounts)
					if err != nil {
						// error happens, retry in 1 minute
						updateTimer(timer, n.Remote(), retryTime)
					} else {
						// send command success, check for an average block time
						updateTimer(timer, n.Remote(), blockTime)
					}
				} else {
					updateTimer(timer, n.Remote(), t.duration)
				}

			case <-t.shutdownChan:
				fmt.Println("receive shutdown signal, terminate periodic tasks")
				return
			}
		}
	}(t.duration)
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
