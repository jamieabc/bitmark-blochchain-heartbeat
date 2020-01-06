package block

import (
	"fmt"
	"time"

	"github.com/jamieabc/bitmarkd-broadcast-monitor/nodes/node"
)

const (
	minNextDuration = 100 * time.Millisecond
)

func DurationToNextCheck(remote node.Remote, targetDuration time.Duration) (time.Duration, error) {
	latestReceiveTime, err := latestBlockGenerationTime(remote)
	if err != nil {
		return time.Duration(0), err
	}

	now := time.Now()

	// due to incorrect of network time, treat all future time as now
	if latestReceiveTime.After(now) {
		return targetDuration, nil
	}

	durationOfLatestReceivedToNow := now.Sub(latestReceiveTime)
	fmt.Printf("now: %s, target duration: %s, latest block receive time: %s\n", now, targetDuration, latestReceiveTime)

	// long time w/o issuance
	if targetDuration < durationOfLatestReceivedToNow {
		return targetDuration, nil
	}

	return targetDuration - durationOfLatestReceivedToNow + minNextDuration, nil
}

func latestBlockGenerationTime(remote node.Remote) (time.Time, error) {
	info, err := remote.Info()
	if nil != err {
		fmt.Printf("get remote info with error: %s\n", err)
		return time.Time{}, err
	}

	resp, err := remote.BlockHeader(info.Height)
	if nil != err {
		fmt.Printf("get remote block resp of height %d with error: %s\n", info.Height, err)
		return time.Time{}, err
	}
	return time.Unix(int64(resp.Header.Timestamp), 0), nil
}
