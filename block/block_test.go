package block_test

import (
	"testing"
	"time"

	"github.com/jamieabc/bitmark-blochchain-heartbeat/block"
	"github.com/stretchr/testify/assert"

	"github.com/bitmark-inc/bitmarkd/blockrecord"
	"github.com/bitmark-inc/bitmarkd/merkle"

	"github.com/bitmark-inc/bitmarkd/blockdigest"
	"github.com/jamieabc/bitmarkd-broadcast-monitor/communication"

	"github.com/golang/mock/gomock"
	"github.com/jamieabc/bitmark-blochchain-heartbeat/block/mocks"
)

func newMockRemote(t *testing.T) (*gomock.Controller, *mocks.MockRemote) {
	ctl := gomock.NewController(t)
	mock := mocks.NewMockRemote(ctl)
	return ctl, mock
}

func TestNextCheckTimeWhenLatestReceiveInHour(t *testing.T) {
	ctl, mock := newMockRemote(t)
	defer ctl.Finish()

	now := time.Now()
	height := uint64(1000)

	mock.EXPECT().Info().Return(&communication.InfoResponse{
		Version: "",
		Chain:   "",
		Normal:  false,
		Height:  height,
	}, nil).Times(1)

	mock.EXPECT().BlockHeader(height).Return(&communication.BlockHeaderResponse{
		Digest: blockdigest.Digest{},
		Header: &blockrecord.Header{
			Version:          0,
			TransactionCount: 0,
			Number:           0,
			PreviousBlock:    blockdigest.Digest{},
			MerkleRoot:       merkle.Digest{},
			Timestamp:        uint64(now.Add(-20 * time.Minute).Unix()),
			Difficulty:       nil,
			Nonce:            0,
		},
	}, nil).Times(1)

	actual, err := block.DurationToNextCheck(mock, time.Hour)
	assert.Equal(t, nil, err, "wrong error")
	assert.Greater(t, 40*time.Minute.Nanoseconds(), actual.Nanoseconds(), "wrong upper bound of check time")
	assert.Less(t, 39*time.Minute.Nanoseconds(), actual.Nanoseconds(), "wrong lower bound of check time")
}

func TestNextCheckTimeWhenLatestReceiveOverHour(t *testing.T) {
	ctl, mock := newMockRemote(t)
	defer ctl.Finish()

	now := time.Now()
	height := uint64(1000)

	mock.EXPECT().Info().Return(&communication.InfoResponse{
		Version: "",
		Chain:   "",
		Normal:  false,
		Height:  height,
	}, nil).Times(1)

	mock.EXPECT().BlockHeader(height).Return(&communication.BlockHeaderResponse{
		Digest: blockdigest.Digest{},
		Header: &blockrecord.Header{
			Version:          0,
			TransactionCount: 0,
			Number:           0,
			PreviousBlock:    blockdigest.Digest{},
			MerkleRoot:       merkle.Digest{},
			Timestamp:        uint64(now.Add(-3 * time.Hour).Unix()),
			Difficulty:       nil,
			Nonce:            0,
		},
	}, nil).Times(1)

	actual, err := block.DurationToNextCheck(mock, time.Hour)
	assert.Equal(t, nil, err, "wrong error")
	assert.Equal(t, time.Hour.Nanoseconds(), actual.Nanoseconds(), "wrong next check time")
}

func TestNextCheckTimeWhenLatestReceiveInFuture(t *testing.T) {
	ctl, mock := newMockRemote(t)
	defer ctl.Finish()

	now := time.Now()
	height := uint64(1000)

	mock.EXPECT().Info().Return(&communication.InfoResponse{
		Version: "",
		Chain:   "",
		Normal:  false,
		Height:  height,
	}, nil).Times(1)

	mock.EXPECT().BlockHeader(height).Return(&communication.BlockHeaderResponse{
		Digest: blockdigest.Digest{},
		Header: &blockrecord.Header{
			Version:          0,
			TransactionCount: 0,
			Number:           0,
			PreviousBlock:    blockdigest.Digest{},
			MerkleRoot:       merkle.Digest{},
			Timestamp:        uint64(now.Add(5 * time.Minute).Unix()),
			Difficulty:       nil,
			Nonce:            0,
		},
	}, nil).Times(1)

	actual, err := block.DurationToNextCheck(mock, time.Hour)
	assert.Equal(t, nil, err, "wrong error")
	assert.Equal(t, time.Hour, actual, "wrong duration")
}
