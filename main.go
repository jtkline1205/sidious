package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gorilla/mux"
)

// Structs

type StrategyBody struct {
	Cards  []string `json:"cards"`
	UpCard string   `json:"upCard"`
}

type MapRuneBoolResponseBody struct {
	Entries map[rune]bool `json:"entries"`
}

type CardsBody struct {
	Cards []string `json:"cards"`
}

type ValuesBody struct {
	Values []uint8 `json:"values"`
}

type SingleStringBody struct {
	SingleString string `json:"singleString"`
}

type StringsAndStringBody struct {
	Strings []string `json:"strings"`
	String  string   `json:"string"`
}

type StringsBody struct {
	Strings []string `json:"strings"`
}

type Rank struct {
	Name           string `json:"name"`
	BlackjackValue uint8  `json:"blackjackValue"`
	BaccaratValue  uint8  `json:"baccaratValue"`
	Label          string `json:"label"`
	Order          uint8  `json:"order"`
}

type Suit struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

type PokerHandType struct {
	Name     string `json:"name"`
	Strength uint8  `json:"strength"`
}

type Card struct {
	RankLabel string
	Suit      string
}

type Deck struct {
	Cards []Card
}

type Shoe struct {
	Decks []*Deck
}

// Package Variables

var ranks []Rank
var suits []Suit
var pokerHandTypes []PokerHandType
var SizeToShoeMap = make(map[int]*Shoe)
var SizeToSequencedCardsMap = make(map[int][]Card)
var mutex = &sync.Mutex{}

// Functions

func NewDeck() *Deck {
	var cards []Card

	for _, suit := range []string{"H", "S", "D", "C"} {
		for _, rank := range []string{"2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K", "A"} {
			cards = append(cards, Card{RankLabel: rank, Suit: suit})
		}
	}

	return &Deck{Cards: cards}
}

func (s *Shoe) CardsLeft() int {
	totalCards := 0
	for _, deck := range s.Decks {
		totalCards += len(deck.Cards)
	}
	return totalCards
}

func (d *Deck) DrawCard() Card {
	if len(d.Cards) == 0 {
		panic("deck is empty")
	}

	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(d.Cards))
	drawnCard := d.Cards[index]
	d.Cards = append(d.Cards[:index], d.Cards[index+1:]...)
	return drawnCard
}

func (d *Deck) HasCards() bool {
	return len(d.Cards) > 0
}

func (card *Card) String() string {
	var result string
	result += fmt.Sprintf("%s%s", card.RankLabel, card.Suit)
	return result
}

func (d *Deck) String() string {
	var result string
	for _, card := range d.Cards {
		result += fmt.Sprintf("%s%s\n", card.RankLabel, card.Suit)
	}
	return result
}

func NewShoe(numDecks int) *Shoe {
	var decks []*Deck
	for i := 0; i < numDecks; i++ {
		decks = append(decks, NewDeck())
	}
	return &Shoe{Decks: decks}
}

func (s *Shoe) DrawCard() Card {
	if len(s.Decks) == 0 {
		panic("shoe has no decks")
	}

	shoeSize := len(s.Decks)
	fmt.Println("drawing a card from shoeSize = " + strconv.Itoa(shoeSize))
	sequencedCards := SizeToSequencedCardsMap[shoeSize]
	if len(sequencedCards) == 0 {
		fmt.Println("sequenced cards has len = 0")
		// Choose a random deck from the shoe
		rand.Seed(time.Now().UnixNano())
		deckIndex := rand.Intn(len(s.Decks))
		if s.Decks[deckIndex].HasCards() {
			return s.Decks[deckIndex].DrawCard()
		} else {
			ResetShoe(shoeSize)
			return SizeToShoeMap[shoeSize].DrawCard()
		}

	} else {
		fmt.Println("sequenced cards has len != 0")
		drawnSequencedCard := sequencedCards[0]
		fmt.Println("chose drawnSequencedCard")
		fmt.Println(drawnSequencedCard)
		SizeToSequencedCardsMap[shoeSize] = sequencedCards[1:]
		return drawnSequencedCard
	}
}

func (s *Shoe) HasCards() bool {
	for _, deck := range s.Decks {
		if deck.HasCards() {
			return true
		}
	}
	return false
}

func (s *Shoe) String() string {
	var result string
	for i, deck := range s.Decks {
		result += fmt.Sprintf("Deck %d:\n%s", i+1, deck)
	}
	return result
}

func CalculateIsNatural(cards []string) bool {
	return (CalculateBaccaratValueForCards(cards) == 9 || CalculateBaccaratValueForCards(cards) == 8) && len(cards) == 2
}

func CalculateIsSoft(cards []string) bool {
	var value, acesAs11s int

	for _, card := range cards {
		value += CalculateBlackjackValueForCard(card)
		if card[0] == 'A' {
			acesAs11s++
		}
	}

	for i := 0; i < acesAs11s; i++ {
		if value > 21 {
			value -= 10
		}
	}

	return acesAs11s > 0 && value <= 21
}

func CalculateIsBlackjack(cards []string) bool {
	if len(cards) != 2 {
		return false
	}

	dealerUpCard := cards[0]
	dealerHoleCard := cards[1]

	if dealerUpCard[0] == 'A' {
		return CalculateBlackjackValueForCard(dealerHoleCard) == 10
	} else if CalculateBlackjackValueForCard(dealerUpCard) == 10 {
		return dealerHoleCard[0] == 'A'
	}

	return false
}

func CalculateBlackjackValueForCard(card string) int {
	if len(card) == 2 {
		for _, item := range ranks {
			if item.Label[0] == card[0] {
				return int(item.BlackjackValue)
			}
		}
	} else {
		for _, item := range ranks {
			if item.Label == card {
				return int(item.BlackjackValue)
			}
		}
	}

	return -1
}

func CalculateBaccaratValueForCard(card string) int {
	if len(card) == 2 {
		for _, item := range ranks {
			if item.Label[0] == card[0] {
				return int(item.BaccaratValue)
			}
		}
	} else {
		for _, item := range ranks {
			if item.Label == card {
				return int(item.BaccaratValue)
			}
		}
	}

	return -1
}

func CalculateBlackjackValueForCards(cards []string) int {
	returnValue := 0
	for _, card := range cards {
		returnValue += CalculateBlackjackValueForCard(card)
	}

	acesAs11s := 0
	for _, card := range cards {
		if card[0] == 'A' {
			acesAs11s++
		}
	}

	for i := 0; i < acesAs11s; i++ {
		if returnValue > 21 {
			returnValue -= 10
		}
	}

	return returnValue
}

func CalculateBaccaratValueForCards(cards []string) int {
	returnValue := 0
	for _, card := range cards {
		returnValue += CalculateBaccaratValueForCard(card)
	}

	return returnValue % 10
}

func MakeCardsFromStrings(cardStrings []string) []Card {
	var cards []Card

	for _, cardString := range cardStrings {
		if len(cardString) >= 2 {
			card := Card{
				RankLabel: string(cardString[0]),
				Suit:      string(cardString[1]),
			}
			cards = append(cards, card)
		}
	}

	return cards
}

func SetCardsInShoe(shoeSize int, cards []string) bool {
	println("setting cards in shoe with shoeSize = " + fmt.Sprint(shoeSize))
	println("cards = " + fmt.Sprint(cards))
	cardStructs := MakeCardsFromStrings(cards)
	println("cardStructs = " + fmt.Sprint(cardStructs))

	SizeToSequencedCardsMap[shoeSize] = cardStructs

	fmt.Println(SizeToSequencedCardsMap)

	return true
}

func ResetShoe(shoeSize int) bool {
	SizeToShoeMap[shoeSize] = NewShoe(shoeSize)
	return true
}

func FindOrderForRank(rank string) int {
	for _, item := range ranks {
		if item.Label == rank {
			return int(item.Order)
		}
	}
	return 0
}

func UpdateFlushKeyRanks(cards []string, flushSuit rune) map[rune]bool {
	flushKeyRanks := map[rune]bool{'N': true}

	for _, card := range cards {
		if rune(card[1]) == flushSuit && !flushKeyRanks[rune(card[0])] {
			minFlushKeyRank := 'K'
			for rank := range flushKeyRanks {
				if FindOrderForRank(string(rank)) < FindOrderForRank(string(minFlushKeyRank)) && rank != 'A' {
					minFlushKeyRank = rank
				}
			}

			if rune(card[0]) == 'A' || FindOrderForRank(string(rune(card[0]))) > FindOrderForRank(string(minFlushKeyRank)) {
				delete(flushKeyRanks, minFlushKeyRank)
				flushKeyRanks[rune(card[0])] = true
			}
		}
	}

	return flushKeyRanks
}

func CalculateIsFlush(cards []string) bool {
	suitFreqs := make(map[rune]int)

	for _, card := range cards {
		suit := rune(card[1])
		suitFreqs[suit]++
		if suitFreqs[suit] >= 5 {
			return true
		}
	}

	return false
}

func FindMaxRank(ranks []rune) rune {
	currentMaxRank := 'N'

	for _, rank := range ranks {
		if currentMaxRank != 'A' && (FindOrderForRank(string(rank)) > FindOrderForRank(string(currentMaxRank)) || rank == 'A') {
			currentMaxRank = rank
		}
	}

	return currentMaxRank
}

// Handlers

func GetRankBlackjackValueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	itemLabel := params["label"]
	blackjackValue := CalculateBlackjackValueForCard(itemLabel)
	json.NewEncoder(w).Encode(blackjackValue)
}

func GetBlackjackValueForCardsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cardsBody CardsBody
	err := json.NewDecoder(r.Body).Decode(&cardsBody)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	cards := cardsBody.Cards

	blackjackValue := CalculateBlackjackValueForCards(cards)

	json.NewEncoder(w).Encode(blackjackValue)
}

func GetBaccaratValueForCardsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cardsBody CardsBody
	err := json.NewDecoder(r.Body).Decode(&cardsBody)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	cards := cardsBody.Cards

	baccaratValue := CalculateBaccaratValueForCards(cards)

	json.NewEncoder(w).Encode(baccaratValue)
}

func GetBaccaratNaturalHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cardsBody CardsBody
	err := json.NewDecoder(r.Body).Decode(&cardsBody)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	cards := cardsBody.Cards

	if CalculateIsNatural(cards) {
		json.NewEncoder(w).Encode(true)
	} else {
		json.NewEncoder(w).Encode(false)
	}
}

func GetRankBaccaratValueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	itemLabel := params["label"]

	for _, item := range ranks {
		if item.Label == itemLabel {
			json.NewEncoder(w).Encode(item.BaccaratValue)
			return
		}
	}

	http.NotFound(w, r)
}

func GetBlackjackSoftHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cardsBody CardsBody
	err := json.NewDecoder(r.Body).Decode(&cardsBody)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	cards := cardsBody.Cards

	isSoft := CalculateIsSoft(cards)
	json.NewEncoder(w).Encode(isSoft)
}

func GetBlackjackBustHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cardsBody CardsBody
	err := json.NewDecoder(r.Body).Decode(&cardsBody)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	cards := cardsBody.Cards

	cardsValue := CalculateBlackjackValueForCards(cards)
	json.NewEncoder(w).Encode(cardsValue > 21)
}

func GetRankOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	itemLabel := params["label"]

	orderForRank := FindOrderForRank(itemLabel)

	json.NewEncoder(w).Encode(orderForRank)
}

func GetOrderRankHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	itemOrderInt, err := strconv.Atoi(params["order"])
	itemOrder := uint8(itemOrderInt)

	if err != nil {
		http.Error(w, "Invalid order", http.StatusBadRequest)
		return
	}
	for _, item := range ranks {
		if item.Order == itemOrder {
			json.NewEncoder(w).Encode(item.Label)
			return
		}
	}

	http.NotFound(w, r)
}

func DrawCardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	shoeSizeStr := vars["shoeSize"]

	shoeSize, err := strconv.Atoi(shoeSizeStr)
	if err != nil {
		http.Error(w, "Invalid shoe size", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	shoe, found := SizeToShoeMap[shoeSize]
	if !found {
		http.NotFound(w, r)
		return
	}

	drawnCard := shoe.DrawCard()
	json.NewEncoder(w).Encode(drawnCard.String())
	return
}

func CardsLeftHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	shoeSizeStr := vars["shoeSize"]

	shoeSize, err := strconv.Atoi(shoeSizeStr)
	if err != nil {
		http.Error(w, "Invalid shoe size", http.StatusBadRequest)
		return
	}

	shoe, found := SizeToShoeMap[shoeSize]
	if !found {
		http.NotFound(w, r)
		return
	}

	cardsLeft := shoe.CardsLeft()
	json.NewEncoder(w).Encode(cardsLeft)
	return
}

func CardResourceNameHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryParams := r.URL.Query()

	rankStr := queryParams.Get("rank")
	suitStr := queryParams.Get("suit")

	var foundRankName string
	var foundSuitName string

	for _, item := range ranks {
		if item.Label == rankStr {
			foundRankName = item.Name
		}
	}

	for _, item := range suits {
		if item.Label == suitStr {
			foundSuitName = item.Name
		}
	}

	json.NewEncoder(w).Encode(strings.ToLower(foundRankName) + "_" + strings.ToLower(foundSuitName))
	return
}

func GetBlackjackForDealerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cardsBody CardsBody
	err := json.NewDecoder(r.Body).Decode(&cardsBody)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	cards := cardsBody.Cards

	isBlackjack := CalculateIsBlackjack(cards)
	json.NewEncoder(w).Encode(isBlackjack)
}

func GetBlackjackRanksForValuesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var valuesBody ValuesBody
	err := json.NewDecoder(r.Body).Decode(&valuesBody)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	values := valuesBody.Values

	result := make([]string, len(values))

	for i, value := range values {
		for _, item := range ranks {
			if item.BlackjackValue == value {
				result[i] = item.Label + "S"
			}
		}
	}

	json.NewEncoder(w).Encode(result)
}

func GetBlackjackDescriptionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cardsBody CardsBody
	err := json.NewDecoder(r.Body).Decode(&cardsBody)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	cards := cardsBody.Cards

	handValue := CalculateBlackjackValueForCards(cards)

	var result string

	if handValue == 21 && len(cards) == 2 {
		result = " (Natural 21)"
	} else {
		isSoft := CalculateIsSoft(cards)
		if isSoft {
			result = " (Soft " + strconv.Itoa(handValue) + ")"
		} else {
			result = " (Hard " + strconv.Itoa(handValue) + ")"
		}

	}
	json.NewEncoder(w).Encode(result)
}

func GetRankComparisonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	rank1 := params["rank1"]
	rank2 := params["rank2"]

	var result int = 0

	if rank1 == "A" && rank2 != "A" {
		result = 1
	} else if rank1 != "A" && rank2 == "A" {
		result = -1
	} else {
		var rank1Order, rank2Order uint8
		for _, item := range ranks {
			if item.Label == rank1 {
				rank1Order = item.Order
			}
			if item.Label == rank2 {
				rank2Order = item.Order
			}
		}

		if rank1Order > rank2Order {
			result = 1
		} else if rank1Order < rank2Order {
			result = -1
		}
	}

	json.NewEncoder(w).Encode(result)
}

func CalculateStrategyDecision(cards []string, upCard string) string {
	handValue := CalculateBlackjackValueForCards(cards)
	dealerUpCardValue := CalculateBlackjackValueForCard(upCard)
	var result string

	if len(cards) == 2 && cards[0][0] == cards[1][0] && handValue != 10 && handValue != 20 {
		switch cards[0][0] {
		case 'A', '8':
			result = "SPLIT"
		case '2', '3', '7':
			if dealerUpCardValue < 8 {
				result = "SPLIT"
			} else {
				result = "HIT"
			}
		case '4':
			if dealerUpCardValue == 5 || dealerUpCardValue == 6 {
				result = "SPLIT"
			} else {
				result = "HIT"
			}
		case '6':
			if dealerUpCardValue < 7 {
				result = "SPLIT"
			} else {
				result = "HIT"
			}
		default:
			switch dealerUpCardValue {
			case 7, 10, 11:
				result = "STAND"
			default:
				result = "SPLIT"
			}
		}
	} else if CalculateIsSoft(cards) {
		switch handValue {
		case 19, 20:
			result = "STAND"
		case 18:
			switch dealerUpCardValue {
			case 2, 7, 8:
				result = "STAND"
			default:
				if dealerUpCardValue > 2 && dealerUpCardValue < 7 {
					if len(cards) == 2 {
						result = "DOUBLE"
					} else {
						result = "HIT"
					}
				} else {
					result = "HIT"
				}
			}
		case 17:
			if dealerUpCardValue > 2 && dealerUpCardValue < 7 {
				if len(cards) == 2 {
					result = "DOUBLE"
				} else {
					result = "HIT"
				}
			} else {
				result = "HIT"
			}
		case 15, 16:
			if dealerUpCardValue > 3 && dealerUpCardValue < 7 {
				if len(cards) == 2 {
					result = "DOUBLE"
				} else {
					result = "HIT"
				}
			} else {
				result = "HIT"
			}
		case 13, 14:
			if dealerUpCardValue > 4 && dealerUpCardValue < 7 {
				if len(cards) == 2 {
					result = "DOUBLE"
				} else {
					result = "HIT"
				}
			} else {
				result = "HIT"
			}
		default:
			result = "HIT"
		}
	} else {
		if dealerUpCardValue < 7 {
			switch handValue {
			case 9:
				if dealerUpCardValue == 2 {
					result = "HIT"
				} else if len(cards) == 2 {
					result = "DOUBLE"
				} else {
					result = "HIT"
				}
			case 10, 11:
				if len(cards) == 2 {
					result = "DOUBLE"
				} else {
					result = "HIT"
				}
			case 12:
				switch dealerUpCardValue {
				case 2, 3:
					result = "HIT"
				default:
					result = "STAND"
				}
			default:
				if handValue < 9 {
					result = "HIT"
				} else {
					result = "STAND"
				}
			}
		} else if handValue < 17 {
			if (handValue == 10 || handValue == 11) && dealerUpCardValue < handValue && len(cards) == 2 {
				result = "DOUBLE"
			} else {
				result = "HIT"
			}
		} else {
			result = "STAND"
		}
	}
	return result
}

func GetBlackjackStrategyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var strategyBody StrategyBody
	err := json.NewDecoder(r.Body).Decode(&strategyBody)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	cards := strategyBody.Cards
	upCard := strategyBody.UpCard

	result := CalculateStrategyDecision(cards, upCard)

	json.NewEncoder(w).Encode(result)
}

func SetCardsInShoeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var cardsBody CardsBody

	vars := mux.Vars(r)
	shoeSizeStr := vars["shoeSize"]

	shoeSize, shoeSizeErr := strconv.Atoi(shoeSizeStr)
	if shoeSizeErr != nil {
		http.Error(w, "Invalid shoe size", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&cardsBody)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	cards := cardsBody.Cards

	SetCardsInShoe(shoeSize, cards)
	json.NewEncoder(w).Encode(true)
}

func ResetShoeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	shoeSizeStr := vars["shoeSize"]

	shoeSize, shoeSizeErr := strconv.Atoi(shoeSizeStr)
	if shoeSizeErr != nil {
		http.Error(w, "Invalid shoe size", http.StatusBadRequest)
		return
	}

	ResetShoe(shoeSize)
	json.NewEncoder(w).Encode(true)
}

func GetPokerHandStrengthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var singleStringBody SingleStringBody
	err := json.NewDecoder(r.Body).Decode(&singleStringBody)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	singleString := singleStringBody.SingleString

	for _, item := range pokerHandTypes {
		if item.Name == singleString {
			json.NewEncoder(w).Encode(item.Strength)
			return
		}
	}

	http.NotFound(w, r)
}

func GetPokerFlushRanksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var stringsAndStringBody StringsAndStringBody
	err := json.NewDecoder(r.Body).Decode(&stringsAndStringBody)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	cards := stringsAndStringBody.Strings
	flushSuit, _ := utf8.DecodeRuneInString(stringsAndStringBody.String)

	updatedFlushRanks := UpdateFlushKeyRanks(cards, flushSuit)
	jsonResponse := MapRuneBoolResponseBody{Entries: updatedFlushRanks}

	json.NewEncoder(w).Encode(jsonResponse)
}

func GetPokerFlushHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cardsBody CardsBody
	err := json.NewDecoder(r.Body).Decode(&cardsBody)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	cards := cardsBody.Cards

	isFlush := CalculateIsFlush(cards)

	json.NewEncoder(w).Encode(isFlush)
}

func GetMaxRankHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var stringsBody StringsBody
	err := json.NewDecoder(r.Body).Decode(&stringsBody)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	strings := stringsBody.Strings
	var runeSlice []rune
	for _, str := range strings {
		runeSlice = append(runeSlice, []rune(str)...)
	}
	maxRank := FindMaxRank(runeSlice)

	json.NewEncoder(w).Encode(maxRank)
}

func init() {
	SetUpRanksAndSuits()
}

func SetUpRanksAndSuits() {
	ranks = append(ranks, Rank{Name: "Ace", BlackjackValue: 11, BaccaratValue: 1, Label: "A", Order: 1})
	ranks = append(ranks, Rank{Name: "Two", BlackjackValue: 2, BaccaratValue: 2, Label: "2", Order: 2})
	ranks = append(ranks, Rank{Name: "Three", BlackjackValue: 3, BaccaratValue: 3, Label: "3", Order: 3})
	ranks = append(ranks, Rank{Name: "Four", BlackjackValue: 4, BaccaratValue: 4, Label: "4", Order: 4})
	ranks = append(ranks, Rank{Name: "Five", BlackjackValue: 5, BaccaratValue: 5, Label: "5", Order: 5})
	ranks = append(ranks, Rank{Name: "Six", BlackjackValue: 6, BaccaratValue: 6, Label: "6", Order: 6})
	ranks = append(ranks, Rank{Name: "Seven", BlackjackValue: 7, BaccaratValue: 7, Label: "7", Order: 7})
	ranks = append(ranks, Rank{Name: "Eight", BlackjackValue: 8, BaccaratValue: 8, Label: "8", Order: 8})
	ranks = append(ranks, Rank{Name: "Nine", BlackjackValue: 9, BaccaratValue: 9, Label: "9", Order: 9})
	ranks = append(ranks, Rank{Name: "Ten", BlackjackValue: 10, BaccaratValue: 0, Label: "T", Order: 10})
	ranks = append(ranks, Rank{Name: "Jack", BlackjackValue: 10, BaccaratValue: 0, Label: "J", Order: 11})
	ranks = append(ranks, Rank{Name: "Queen", BlackjackValue: 10, BaccaratValue: 0, Label: "Q", Order: 12})
	ranks = append(ranks, Rank{Name: "King", BlackjackValue: 10, BaccaratValue: 0, Label: "K", Order: 13})

	suits = append(suits, Suit{Name: "Hearts", Label: "H"})
	suits = append(suits, Suit{Name: "Spades", Label: "S"})
	suits = append(suits, Suit{Name: "Diamonds", Label: "D"})
	suits = append(suits, Suit{Name: "Clubs", Label: "C"})

	pokerHandTypes = append(pokerHandTypes, PokerHandType{Name: "High Card", Strength: 1})
	pokerHandTypes = append(pokerHandTypes, PokerHandType{Name: "Pair", Strength: 2})
	pokerHandTypes = append(pokerHandTypes, PokerHandType{Name: "Two Pair", Strength: 3})
	pokerHandTypes = append(pokerHandTypes, PokerHandType{Name: "Three Of A Kind", Strength: 4})
	pokerHandTypes = append(pokerHandTypes, PokerHandType{Name: "Straight", Strength: 5})
	pokerHandTypes = append(pokerHandTypes, PokerHandType{Name: "Flush", Strength: 6})
	pokerHandTypes = append(pokerHandTypes, PokerHandType{Name: "Full House", Strength: 7})
	pokerHandTypes = append(pokerHandTypes, PokerHandType{Name: "Four Of A Kind", Strength: 8})
	pokerHandTypes = append(pokerHandTypes, PokerHandType{Name: "Straight Flush", Strength: 9})
	pokerHandTypes = append(pokerHandTypes, PokerHandType{Name: "Royal Flush", Strength: 10})

}

func main() {
	SetUpRanksAndSuits()

	SizeToShoeMap[1] = NewShoe(1)
	SizeToShoeMap[2] = NewShoe(2)
	SizeToShoeMap[3] = NewShoe(3)
	SizeToShoeMap[4] = NewShoe(4)
	SizeToShoeMap[5] = NewShoe(5)
	SizeToShoeMap[6] = NewShoe(6)
	SizeToShoeMap[7] = NewShoe(7)
	SizeToShoeMap[8] = NewShoe(8)
	SizeToShoeMap[9] = NewShoe(9)
	SizeToShoeMap[10] = NewShoe(10)

	router := mux.NewRouter()
	router.HandleFunc("/orders/{order}/rank", GetOrderRankHandler).Methods("GET")
	router.HandleFunc("/ranks/{label}/order", GetRankOrderHandler).Methods("GET")
	router.HandleFunc("/ranks/{rank1}/{rank2}", GetRankComparisonHandler).Methods("GET")
	router.HandleFunc("/ranks/max", GetMaxRankHandler).Methods("POST")
	router.HandleFunc("/cards/resourceName", CardResourceNameHandler).Methods("GET")
	router.HandleFunc("/shoes/{shoeSize}/setCards", SetCardsInShoeHandler).Methods("POST")
	router.HandleFunc("/shoes/{shoeSize}/reset", ResetShoeHandler).Methods("POST")
	router.HandleFunc("/shoes/{shoeSize}/draw", DrawCardHandler).Methods("GET")
	router.HandleFunc("/shoes/{shoeSize}/cardsLeft", CardsLeftHandler).Methods("GET")
	router.HandleFunc("/blackjack", GetBlackjackForDealerHandler).Methods("POST")
	router.HandleFunc("/blackjack/strategy", GetBlackjackStrategyHandler).Methods("POST")
	router.HandleFunc("/blackjack/soft", GetBlackjackSoftHandler).Methods("POST")
	router.HandleFunc("/blackjack/bust", GetBlackjackBustHandler).Methods("POST")
	router.HandleFunc("/blackjack/ranks/{label}", GetRankBlackjackValueHandler).Methods("GET")
	router.HandleFunc("/blackjack/values", GetBlackjackValueForCardsHandler).Methods("POST")
	router.HandleFunc("/blackjack/values/ranks", GetBlackjackRanksForValuesHandler).Methods("POST")
	router.HandleFunc("/blackjack/values/description", GetBlackjackDescriptionHandler).Methods("POST")
	router.HandleFunc("/baccarat/natural", GetBaccaratNaturalHandler).Methods("POST")
	router.HandleFunc("/baccarat/value", GetBaccaratValueForCardsHandler).Methods("POST")
	router.HandleFunc("/baccarat/ranks/{label}", GetRankBaccaratValueHandler).Methods("GET")
	router.HandleFunc("/poker/strength", GetPokerHandStrengthHandler).Methods("POST")
	router.HandleFunc("/poker/flush/ranks", GetPokerFlushRanksHandler).Methods("POST")
	router.HandleFunc("/poker/flush", GetPokerFlushHandler).Methods("POST")

	port := 5001
	fmt.Printf("Server is running on :%d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}
