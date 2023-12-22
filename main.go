package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Rank struct {
	Name           string `json:"name"`
	BlackjackValue uint8  `json:"blackjackValue"`
	BaccaratValue  uint8  `json:"baccaratValue"`
	Label          string `json:"label"`
	Order          uint8  `json:"order"`
}

var ranks []Rank

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

func GetRankLabelHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	itemLabel := params["label"]

	for _, item := range ranks {
		if item.Label == itemLabel {
			json.NewEncoder(w).Encode(item.Label)
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

func main() {
	router := mux.NewRouter()

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

	router.HandleFunc("/ranks/{label}/blackjackValue", GetRankBlackjackValueHandler).Methods("GET")
	router.HandleFunc("/ranks/{label}/baccaratValue", GetRankBaccaratValueHandler).Methods("GET")
	router.HandleFunc("/ranks/{label}/label", GetRankLabelHandler).Methods("GET")
	router.HandleFunc("/ranks/{label}/order", GetRankOrderHandler).Methods("GET")

	port := 5001
	fmt.Printf("Server is running on :%d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}
