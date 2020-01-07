package sdk_test

import (
	"math/rand"
	"testing"

	"github.com/jamieabc/bitmark-blochchain-heartbeat/pkg/sdk"
)

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func randomString(length int) string {
	str := make([]rune, length)
	for i := range str {
		str[i] = letters[rand.Intn(len(letters))]
	}
	return string(str)
}

func TestTruncateLongString(t *testing.T) {
	longStr := randomString(2 * sdk.MaximumItemName)
	shortStr := randomString(sdk.MaximumItemName - 1)
	fixture := []struct {
		str      string
		expected string
	}{
		{longStr, longStr[0:sdk.MaximumItemName]},
		{shortStr, shortStr},
	}

	for i, s := range fixture {
		actual := sdk.TruncateLongString(s.str)
		if actual != s.expected {
			t.Errorf("%d testerror truncated string, expected %s but get %s",
				i,
				s.expected,
				actual,
			)
		}
	}
}
