package main

import (
	"testing"
)

func TestBlackjack(t *testing.T) {

	for _, item := range ranks {
		if item.Label == "4" {
			expected := uint8(4)
			if item.BaccaratValue != expected || item.BlackjackValue != expected {
				t.Errorf("blackjack/baccarat expected %d", expected)
			}
		}
		if item.Label == "A" {
			if item.BaccaratValue != 1 || item.BlackjackValue != 11 {
				t.Errorf("blackjack/baccarat ace failed")
			}
		}
	}
}
