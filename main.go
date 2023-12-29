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

func (card Card) String() string {
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

// Handlers

func GetRankBlackjackValueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	itemLabel := params["label"]

	for _, item := range ranks {
		if item.Label == itemLabel {
			json.NewEncoder(w).Encode(item.BlackjackValue)
			return
		}
	}

	http.NotFound(w, r)
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

	port := 5001
	fmt.Printf("Server is running on :%d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}
