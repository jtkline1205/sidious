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

	"github.com/gorilla/mux"
)

// Structs

type StrategyBody struct {
	Cards  []string `json:"cards"`
	UpCard string   `json:"upCard"`
}

type CardsBody struct {
	Cards []string `json:"cards"`
}

type ValuesBody struct {
	Values []uint8 `json:"values"`
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
var SizeToShoeMap = make(map[int]*Shoe)
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
		panic("shoe is empty")
	}

	// Choose a random deck from the shoe
	rand.Seed(time.Now().UnixNano())
	deckIndex := rand.Intn(len(s.Decks))

	// Draw a card from the selected deck
	return s.Decks[deckIndex].DrawCard()
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

	if (CalculateBaccaratValueForCards(cards) == 9 || CalculateBaccaratValueForCards(cards) == 8) && len(cards) == 2 {
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

	for _, item := range ranks {
		if item.Label == itemLabel {
			json.NewEncoder(w).Encode(item.Order)
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

func init() {
	println("init running")
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
}

func main() {
	println("main running")
	SetUpRanksAndSuits()

	SizeToShoeMap[1] = NewShoe(1)
	SizeToShoeMap[2] = NewShoe(2)
	SizeToShoeMap[4] = NewShoe(4)
	SizeToShoeMap[8] = NewShoe(8)
	SizeToShoeMap[9] = NewShoe(9)
	SizeToShoeMap[10] = NewShoe(10)

	router := mux.NewRouter()
	router.HandleFunc("/ranks/{label}/blackjackValue", GetRankBlackjackValueHandler).Methods("GET")
	router.HandleFunc("/ranks/{label}/baccaratValue", GetRankBaccaratValueHandler).Methods("GET")
	router.HandleFunc("/ranks/{label}/order", GetRankOrderHandler).Methods("GET")
	router.HandleFunc("/shoes/{shoeSize}/draw", DrawCardHandler).Methods("GET")
	router.HandleFunc("/shoes/{shoeSize}/cardsLeft", CardsLeftHandler).Methods("GET")
	router.HandleFunc("/cards/resourceName", CardResourceNameHandler).Methods("GET")
	router.HandleFunc("/blackjack/value", GetBlackjackValueForCardsHandler).Methods("POST")
	router.HandleFunc("/baccarat/value", GetBaccaratValueForCardsHandler).Methods("POST")
	router.HandleFunc("/baccarat/natural", GetBaccaratNaturalHandler).Methods("POST")
	router.HandleFunc("/blackjack/soft", GetBlackjackSoftHandler).Methods("POST")
	router.HandleFunc("/blackjack/bust", GetBlackjackBustHandler).Methods("POST")
	router.HandleFunc("/blackjack", GetBlackjackForDealerHandler).Methods("POST")
	router.HandleFunc("/blackjack/values/ranks", GetBlackjackRanksForValuesHandler).Methods("POST")
	router.HandleFunc("/ranks/{rank1}/{rank2}", GetRankComparisonHandler).Methods("GET")
	router.HandleFunc("/blackjack/value/description", GetBlackjackDescriptionHandler).Methods("POST")
	router.HandleFunc("/blackjack/strategy", GetBlackjackStrategyHandler).Methods("POST")

	port := 5001
	fmt.Printf("Server is running on :%d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}
