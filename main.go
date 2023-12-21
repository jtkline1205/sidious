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
}

var ranks []Rank

func GetRankBlackjackValueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	itemName := params["name"]

	for _, item := range ranks {
		if item.Name == itemName {
			json.NewEncoder(w).Encode(item.BlackjackValue)
			return
		}
	}

	http.NotFound(w, r)
}

func GetRankBaccaratValueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	itemName := params["name"]

	for _, item := range ranks {
		if item.Name == itemName {
			json.NewEncoder(w).Encode(item.BaccaratValue)
			return
		}
	}

	http.NotFound(w, r)
}

// func CreateItemHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	var newItem Item
// 	err := json.NewDecoder(r.Body).Decode(&newItem)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	items = append(items, newItem)
// 	json.NewEncoder(w).Encode(items)
// }

func main() {
	router := mux.NewRouter()

	ranks = append(ranks, Rank{Name: "Ace", BlackjackValue: 11, BaccaratValue: 1})
	ranks = append(ranks, Rank{Name: "Two", BlackjackValue: 2, BaccaratValue: 2})
	ranks = append(ranks, Rank{Name: "Three", BlackjackValue: 3, BaccaratValue: 3})
	ranks = append(ranks, Rank{Name: "Four", BlackjackValue: 4, BaccaratValue: 4})
	ranks = append(ranks, Rank{Name: "Five", BlackjackValue: 5, BaccaratValue: 5})
	ranks = append(ranks, Rank{Name: "Six", BlackjackValue: 6, BaccaratValue: 6})
	ranks = append(ranks, Rank{Name: "Seven", BlackjackValue: 7, BaccaratValue: 7})
	ranks = append(ranks, Rank{Name: "Eight", BlackjackValue: 8, BaccaratValue: 8})
	ranks = append(ranks, Rank{Name: "Nine", BlackjackValue: 9, BaccaratValue: 9})
	ranks = append(ranks, Rank{Name: "Ten", BlackjackValue: 10, BaccaratValue: 0})
	ranks = append(ranks, Rank{Name: "Jack", BlackjackValue: 10, BaccaratValue: 0})
	ranks = append(ranks, Rank{Name: "Queen", BlackjackValue: 10, BaccaratValue: 0})
	ranks = append(ranks, Rank{Name: "King", BlackjackValue: 10, BaccaratValue: 0})

	router.HandleFunc("/ranks/{name}/blackjackValue", GetRankBlackjackValueHandler).Methods("GET")
	router.HandleFunc("/ranks/{name}/baccaratValue", GetRankBaccaratValueHandler).Methods("GET")

	port := 5001
	fmt.Printf("Server is running on :%d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}
