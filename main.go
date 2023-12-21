package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Rank struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	BlackjackValue uint8  `json:"blackjackValue"`
}

var ranks []Rank

// var rank = "Spades"

// func GetItemsHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(rank)
// }

func GetRankBlackjackValueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	// itemID := params["id"]
	itemName := params["name"]

	for _, item := range ranks {
		if item.Name == itemName {
			json.NewEncoder(w).Encode(item.BlackjackValue)
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

	ranks = append(ranks, Rank{ID: "1", Name: "Ace", BlackjackValue: 11})
	ranks = append(ranks, Rank{ID: "2", Name: "Two", BlackjackValue: 2})
	ranks = append(ranks, Rank{ID: "3", Name: "Three", BlackjackValue: 3})
	ranks = append(ranks, Rank{ID: "4", Name: "Four", BlackjackValue: 4})
	ranks = append(ranks, Rank{ID: "5", Name: "Five", BlackjackValue: 5})
	ranks = append(ranks, Rank{ID: "6", Name: "Six", BlackjackValue: 6})
	ranks = append(ranks, Rank{ID: "7", Name: "Seven", BlackjackValue: 7})
	ranks = append(ranks, Rank{ID: "8", Name: "Eight", BlackjackValue: 8})
	ranks = append(ranks, Rank{ID: "9", Name: "Nine", BlackjackValue: 9})
	ranks = append(ranks, Rank{ID: "10", Name: "Ten", BlackjackValue: 10})
	ranks = append(ranks, Rank{ID: "11", Name: "Jack", BlackjackValue: 10})
	ranks = append(ranks, Rank{ID: "12", Name: "Queen", BlackjackValue: 10})
	ranks = append(ranks, Rank{ID: "13", Name: "King", BlackjackValue: 10})

	// router.HandleFunc("/rank", GetItemsHandler).Methods("GET")
	router.HandleFunc("/ranks/{name}/blackjackValue", GetRankBlackjackValueHandler).Methods("GET")

	port := 5001
	fmt.Printf("Server is running on :%d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}
