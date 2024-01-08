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

	hardHands := [][]string{
		{"T", "4"},
		{"8", "9"},
		{"2", "T"},
		{"K", "J"},
		{"4", "Q"},
		{"Q", "5"},
	}

	for _, hand := range hardHands {
		soft := CalculateIsSoft(hand)
		if soft {
			t.Errorf("soft hand not expected")
		}
	}

	softHands := [][]string{
		{"A", "4"},
		{"2", "A"},
		{"A", "A"},
		{"3", "A"},
		{"A", "T"},
		{"Q", "A"},
	}

	for _, hand := range softHands {
		soft := CalculateIsSoft(hand)
		if !soft {
			t.Errorf("hard hand not expected")
		}
	}

	blackjackHands := [][]string{
		{"A", "J"},
	}

	for _, hand := range blackjackHands {
		isBlackjack := CalculateIsBlackjack(hand)
		if !isBlackjack {
			t.Errorf("blackjack expected")
		}
	}

	nonBlackjackHands := [][]string{
		{"4", "J"},
		{"2", "2"},
		{"A", "7"},
	}

	for _, hand := range nonBlackjackHands {
		isBlackjack := CalculateIsBlackjack(hand)
		if isBlackjack {
			t.Errorf("blackjack not expected")
		}
	}

}

func TestBlackjackHandCreation(t *testing.T) {

	value4Hands := [][]string{
		{"2", "2"},
	}

	value5Hands := [][]string{
		{"3", "2"},
	}

	value6Hands := [][]string{
		{"4", "2"},
		{"3", "3"},
	}

	value7Hands := [][]string{
		{"5", "2"},
		{"4", "3"},
	}

	value8Hands := [][]string{
		{"6", "2"},
		{"5", "3"},
		{"4", "4"},
	}

	value9Hands := [][]string{
		{"7", "2"},
		{"6", "3"},
		{"5", "4"},
	}

	value10Hands := [][]string{
		{"8", "2"},
		{"7", "3"},
		{"6", "4"},
	}

	value11Hands := [][]string{
		{"9", "2"},
		{"8", "3"},
		{"7", "4"},
	}

	value12Hands := [][]string{
		{"A", "A"},
		{"9", "3"},
		{"8", "4"},
	}

	value13Hands := [][]string{
		{"A", "2"},
		{"9", "4"},
		{"8", "5"},
	}

	value14Hands := [][]string{
		{"A", "3"},
		{"Q", "4"},
		{"9", "5"},
	}

	value15Hands := [][]string{
		{"A", "4"},
		{"Q", "5"},
		{"9", "6"},
	}

	value16Hands := [][]string{
		{"A", "5"},
		{"Q", "6"},
		{"9", "7"},
	}

	value17Hands := [][]string{
		{"A", "6"},
		{"Q", "7"},
		{"9", "8"},
	}

	value18Hands := [][]string{
		{"A", "7"},
		{"Q", "8"},
		{"9", "9"},
	}

	value19Hands := [][]string{
		{"A", "8"},
		{"Q", "9"},
		{"9", "T"},
	}

	value20Hands := [][]string{
		{"A", "9"},
		{"Q", "K"},
		{"T", "J"},
	}

	for _, hand := range value4Hands {
		if CalculateBlackjackValueForCards(hand) != 4 {
			t.Errorf("4 expected")
		}
	}

	for _, hand := range value5Hands {
		if CalculateBlackjackValueForCards(hand) != 5 {
			t.Errorf("5 expected")
		}
	}

	for _, hand := range value6Hands {
		if CalculateBlackjackValueForCards(hand) != 6 {
			t.Errorf("6 expected")
		}
	}

	for _, hand := range value7Hands {
		if CalculateBlackjackValueForCards(hand) != 7 {
			t.Errorf("7 expected")
		}
	}

	for _, hand := range value8Hands {
		if CalculateBlackjackValueForCards(hand) != 8 {
			t.Errorf("8 expected")
		}
	}

	for _, hand := range value9Hands {
		if CalculateBlackjackValueForCards(hand) != 9 {
			t.Errorf("9 expected")
		}
	}

	for _, hand := range value10Hands {
		if CalculateBlackjackValueForCards(hand) != 10 {
			t.Errorf("10 expected")
		}
	}

	for _, hand := range value11Hands {
		if CalculateBlackjackValueForCards(hand) != 11 {
			t.Errorf("11 expected")
		}
	}

	for _, hand := range value12Hands {
		if CalculateBlackjackValueForCards(hand) != 12 {
			t.Errorf("12 expected")
		}
	}

	for _, hand := range value13Hands {
		if CalculateBlackjackValueForCards(hand) != 13 {
			t.Errorf("13 expected")
		}
	}

	for _, hand := range value14Hands {
		if CalculateBlackjackValueForCards(hand) != 14 {
			t.Errorf("14 expected")
		}
	}

	for _, hand := range value15Hands {
		if CalculateBlackjackValueForCards(hand) != 15 {
			t.Errorf("15 expected")
		}
	}

	for _, hand := range value16Hands {
		if CalculateBlackjackValueForCards(hand) != 16 {
			t.Errorf("16 expected")
		}
	}

	for _, hand := range value17Hands {
		if CalculateBlackjackValueForCards(hand) != 17 {
			t.Errorf("17 expected")
		}
	}

	for _, hand := range value18Hands {
		if CalculateBlackjackValueForCards(hand) != 18 {
			t.Errorf("18 expected")
		}
	}

	for _, hand := range value19Hands {
		if CalculateBlackjackValueForCards(hand) != 19 {
			t.Errorf("19 expected")
		}
	}

	for _, hand := range value20Hands {
		if CalculateBlackjackValueForCards(hand) != 20 {
			t.Errorf("20 expected")
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

func TestBaccaratNaturals(t *testing.T) {

	naturalHands := [][]string{
		{"2", "7"},
		{"A", "8"},
		{"4", "5"},
		{"6", "3"},
		{"5", "4"},
		{"8", "A"},
		{"3", "6"},
		{"7", "2"},
	}

	for _, hand := range naturalHands {
		isNatural := CalculateIsNatural(hand)
		if !isNatural {
			t.Errorf("natural expected")
		}
	}

	nonNaturalHands := [][]string{
		{"4", "K"},
		{"2", "2"},
		{"5", "7"},
		{"5", "7", "7"},
		{"K", "Q", "9"},
		{"2", "4", "3"},
		{"A", "8", "9"},
	}

	for _, hand := range nonNaturalHands {
		if CalculateIsNatural(hand) {
			t.Errorf("natural not expected")
		}
	}

}

func TestBaccaratHandCreation(t *testing.T) {

	value4Hands := [][]string{
		{"2", "2"},
		{"A", "A", "2"},
	}

	value5Hands := [][]string{
		{"3", "2"},
	}

	value6Hands := [][]string{
		{"4", "2"},
		{"3", "3"},
	}

	value7Hands := [][]string{
		{"5", "2"},
		{"4", "3"},
	}

	value8Hands := [][]string{
		{"6", "2"},
		{"5", "3"},
		{"4", "4"},
	}

	value9Hands := [][]string{
		{"7", "2"},
		{"6", "3"},
		{"5", "4"},
	}

	value0Hands := [][]string{
		{"8", "2"},
		{"7", "3"},
		{"6", "4"},
		{"5", "5"},
	}

	value1Hands := [][]string{
		{"9", "2"},
		{"8", "3"},
		{"7", "4"},
		{"6", "5"},
	}

	value2Hands := [][]string{
		{"A", "A"},
		{"K", "2"},
		{"9", "3"},
		{"8", "4"},
		{"7", "5"},
		{"6", "6"},
	}

	value3Hands := [][]string{
		{"A", "2"},
		{"K", "3"},
		{"9", "4"},
		{"8", "5"},
		{"7", "6"},
	}

	for _, hand := range value0Hands {
		if CalculateBaccaratValueForCards(hand) != 0 {
			t.Errorf("0 expected")
		}
	}

	for _, hand := range value1Hands {
		if CalculateBaccaratValueForCards(hand) != 1 {
			t.Errorf("1 expected")
		}
	}

	for _, hand := range value2Hands {
		if CalculateBaccaratValueForCards(hand) != 2 {
			t.Errorf("2 expected")
		}
	}

	for _, hand := range value3Hands {
		if CalculateBaccaratValueForCards(hand) != 3 {
			t.Errorf("3 expected")
		}
	}

	for _, hand := range value4Hands {
		if CalculateBaccaratValueForCards(hand) != 4 {
			t.Errorf("4 expected")
		}
	}

	for _, hand := range value5Hands {
		if CalculateBaccaratValueForCards(hand) != 5 {
			t.Errorf("5 expected")
		}
	}

	for _, hand := range value6Hands {
		if CalculateBaccaratValueForCards(hand) != 6 {
			t.Errorf("6 expected")
		}
	}

	for _, hand := range value7Hands {
		if CalculateBaccaratValueForCards(hand) != 7 {
			t.Errorf("7 expected")
		}
	}

	for _, hand := range value8Hands {
		if CalculateBaccaratValueForCards(hand) != 8 {
			t.Errorf("8 expected")
		}
	}

	for _, hand := range value9Hands {
		if CalculateBaccaratValueForCards(hand) != 9 {
			t.Errorf("9 expected")
		}
	}

}

// func TestPoker(t *testing.T) {
// 	cards := []string{"2S", "3S", "4S", "5S", "6S"}

// 	flushRanks := UpdateFlushKeyRanks(cards, flushSuit)
// }
