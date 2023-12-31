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

func TestBlackjackDecisions(t *testing.T) {
	weakDealerCards := []string{"2S", "3S", "4S", "5S", "6S"}
	strongDealerCards := []string{"7S", "8S", "9S", "TS", "AS"}

	softHandsToStand := [][]string{
		{"A", "8"},
		{"A", "9"},
	}

	for _, softHand := range softHandsToStand {
		for _, strongDealerCard := range strongDealerCards {
			result := CalculateStrategyDecision(softHand, strongDealerCard)
			if result != "STAND" {
				t.Errorf("blackjack strategy expected STAND")
			}
		}

		for _, weakDealerCard := range weakDealerCards {
			result := CalculateStrategyDecision(softHand, weakDealerCard)
			if result != "STAND" {
				t.Errorf("blackjack strategy expected STAND")
			}
		}
	}

	softHandsToDecide := [][]string{
		{"A", "7"},
	}

	for _, softHand := range softHandsToDecide {
		for _, strongDealerCard := range strongDealerCards {
			result := CalculateStrategyDecision(softHand, strongDealerCard)
			dealerCardValue := CalculateBlackjackValueForCard(strongDealerCard)
			if dealerCardValue < 9 && result != "STAND" {
				t.Errorf("blackjack strategy expected STAND")
			}
			if dealerCardValue >= 9 && result != "HIT" {
				t.Errorf("blackjack strategy expected HIT")
			}
		}
		for _, weakDealerCard := range weakDealerCards {
			result := CalculateStrategyDecision(softHand, weakDealerCard)
			dealerCardValue := CalculateBlackjackValueForCard(weakDealerCard)
			if dealerCardValue < 3 && result != "STAND" {
				t.Errorf("blackjack strategy expected STAND")
			}
			if dealerCardValue >= 3 && result != "DOUBLE" {
				t.Errorf("blackjack strategy expected DOUBLE")
			}
		}
	}

	moreSoftHandsToCheck := [][]string{
		{"A", "2"},
		{"A", "3"},
		{"A", "4"},
		{"A", "5"},
		{"A", "6"},
	}

	for _, softHand := range moreSoftHandsToCheck {
		for _, strongDealerCard := range strongDealerCards {
			result := CalculateStrategyDecision(softHand, strongDealerCard)
			if result != "HIT" {
				t.Errorf("blackjack strategy expected HIT")
			}
		}
		for _, weakDealerCard := range weakDealerCards {
			handValue := CalculateBlackjackValueForCards(softHand)
			result := CalculateStrategyDecision(softHand, weakDealerCard)
			weakDealerCardValue := CalculateBlackjackValueForCard(weakDealerCard)
			switch handValue {
			case 13, 14:
				if weakDealerCardValue < 5 && result != "HIT" {
					t.Errorf("blackjack strategy expected HIT")
				}
				if weakDealerCardValue >= 5 && result != "DOUBLE" {
					t.Errorf("blackjack strategy expected DOUBLE")
				}
			case 15, 16:
				if weakDealerCardValue < 4 && result != "HIT" {
					t.Errorf("blackjack strategy expected HIT")
				}
				if weakDealerCardValue >= 4 && result != "DOUBLE" {
					t.Errorf("blackjack strategy expected DOUBLE")
				}
			case 17:
				if weakDealerCardValue < 3 && result != "HIT" {
					t.Errorf("blackjack strategy expected HIT")
				}
				if weakDealerCardValue >= 3 && result != "DOUBLE" {
					t.Errorf("blackjack strategy expected DOUBLE")
				}
			}
		}

	}

	splittableHands := [][]string{
		{"A", "A"},
		{"2", "2"},
		{"3", "3"},
		{"6", "6"},
		{"7", "7"},
		{"8", "8"},
		{"9", "9"},
	}

	for _, splittableHand := range splittableHands {
		for _, weakDealerCard := range weakDealerCards {
			result := CalculateStrategyDecision(splittableHand, weakDealerCard)
			if result != "SPLIT" {
				t.Errorf("blackjack strategy expected SPLIT")
			}
		}
	}

	strongSplits := [][]string{
		{"A", "A"},
		{"8", "8"},
	}

	for _, strongSplit := range strongSplits {
		for _, strongDealerCard := range strongDealerCards {
			result := CalculateStrategyDecision(strongSplit, strongDealerCard)
			if result != "SPLIT" {
				t.Errorf("blackjack strategy expected SPLIT")
			}
		}
	}

	nineHands := [][]string{
		{"2", "7"},
		{"5", "4"},
	}

	for _, nineHand := range nineHands {
		for _, strongDealerCard := range strongDealerCards {
			result := CalculateStrategyDecision(nineHand, strongDealerCard)
			if result != "HIT" {
				t.Errorf("blackjack strategy expected HIT")
			}
		}
		result := CalculateStrategyDecision(nineHand, "2S")
		if result != "HIT" {
			t.Errorf("blackjack strategy expected HIT")
		}
	}

	twelveHands := [][]string{
		{"2", "Q"},
		{"5", "7"},
	}

	for _, hand := range twelveHands {
		for _, card := range weakDealerCards {
			cardValue := CalculateBlackjackValueForCard(card)
			decision := CalculateStrategyDecision(hand, card)
			if cardValue == 2 || cardValue == 3 {
				if decision != "HIT" {
					t.Errorf("HIT expected")
				}
			} else {
				if decision != "STAND" {
					t.Errorf("STAND expected")
				}
			}
		}
		for _, card := range strongDealerCards {
			decision := CalculateStrategyDecision(hand, card)
			if decision != "HIT" {
				t.Errorf("HIT expected")
			}
		}
	}

	elevenHands := [][]string{
		{"2", "9"},
	}

	for _, hand := range elevenHands {
		for _, card := range weakDealerCards {
			decision := CalculateStrategyDecision(hand, card)
			if decision != "DOUBLE" {
				t.Errorf("DOUBLE expected")
			}
		}
		for _, card := range strongDealerCards {
			decision := CalculateStrategyDecision(hand, card)
			cardValue := CalculateBlackjackValueForCard(card)
			if cardValue == 11 {
				if decision != "HIT" {
					t.Errorf("HIT expected")
				}
			} else {
				if decision != "DOUBLE" {
					t.Errorf("DOUBLE expected")
				}
			}
		}
	}

	tenHands := [][]string{
		{"2", "8"},
		{"7", "3"},
	}

	for _, hand := range tenHands {
		for _, card := range weakDealerCards {
			decision := CalculateStrategyDecision(hand, card)
			if decision != "DOUBLE" {
				t.Errorf("DOUBLE expected")
			}
		}
		for _, card := range strongDealerCards {
			cardValue := CalculateBlackjackValueForCard(card)
			decision := CalculateStrategyDecision(hand, card)
			if cardValue == 10 || cardValue == 11 {
				if decision != "HIT" {
					t.Errorf("HIT expected")
				}
			} else {
				if decision != "DOUBLE" {
					t.Errorf("DOUBLE expected")
				}
			}
		}
	}

	weakDealerNineHands := [][]string{
		{"2", "7"},
		{"6", "3"},
	}

	for _, hand := range weakDealerNineHands {
		for _, card := range weakDealerCards {
			cardValue := CalculateBlackjackValueForCard(card)
			decision := CalculateStrategyDecision(hand, card)
			if cardValue == 2 {
				if decision != "HIT" {
					t.Errorf("HIT expected")
				}
			} else {
				if decision != "DOUBLE" {
					t.Errorf("DOUBLE expected")
				}
			}
		}
	}

	lowNonSplittablePlayerHands := [][]string{
		{"2", "6"},
		{"2", "5"},
		{"2", "4"},
		{"2", "3"},
	}

	for _, hand := range lowNonSplittablePlayerHands {
		for _, card := range strongDealerCards {
			decision := CalculateStrategyDecision(hand, card)
			if decision != "HIT" {
				t.Errorf("HIT expected")
			}
		}
		for _, card := range weakDealerCards {
			decision := CalculateStrategyDecision(hand, card)
			if decision != "HIT" {
				t.Errorf("HIT expected")
			}
		}
	}

	strongPlayerHands := [][]string{
		{"7", "Q"},
		{"8", "K"},
		{"9", "K"},
		{"J", "J"},
	}

	for _, hand := range strongPlayerHands {
		for _, card := range weakDealerCards {
			decision := CalculateStrategyDecision(hand, card)
			if decision != "STAND" {
				t.Errorf("STAND expected")
			}
		}
		for _, card := range strongDealerCards {
			decision := CalculateStrategyDecision(hand, card)
			if decision != "STAND" {
				t.Errorf("STAND expected")
			}
		}
	}

	weakPlayerHands := [][]string{
		{"3", "Q"},
		{"4", "K"},
		{"5", "K"},
		{"6", "J"},
	}

	for _, hand := range weakPlayerHands {
		for _, card := range weakDealerCards {
			decision := CalculateStrategyDecision(hand, card)
			if decision != "STAND" {
				t.Errorf("STAND expected")
			}
		}
		for _, card := range strongDealerCards {
			decision := CalculateStrategyDecision(hand, card)
			if decision != "HIT" {
				t.Errorf("HIT expected")
			}
		}
	}

}
