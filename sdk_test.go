package main

import (
	"math/rand"
	"testing"
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
	longStr := randomString(2 * maximumItemName)
	shortStr := randomString(maximumItemName - 1)
	fixture := []struct {
		str      string
		expected string
	}{
		{longStr, longStr[0:maximumItemName]},
		{shortStr, shortStr},
	}

	for i, s := range fixture {
		actual := truncateLongString(s.str)
		if actual != s.expected {
			t.Errorf("%d testerror truncated string, expected %s but get %s",
				i,
				s.expected,
				actual,
			)
		}
	}
}
