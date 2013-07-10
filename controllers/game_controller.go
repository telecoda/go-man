package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/telecoda/go-man/models"
	"log"
	"net/http"
)

func GameList(w http.ResponseWriter, r *http.Request) {
	addResponseHeaders(w)

	stateFilter := r.URL.Query().Get("state")
	fmt.Println("Query parameters:", stateFilter)

	// get all games
	boards, err := models.ReadAllGameBoards(stateFilter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	returnBoardsSummaryAsJson(w, boards)
	//fmt.Fprint(w, Response{"success": true, "message": "Here are the current games", "method": r.Method})
}

func GameCreate(w http.ResponseWriter, r *http.Request) {

	log.Println("GameCreate started")
	addResponseHeaders(w)

	var board = models.NewGameBoard()

	board.SaveGameBoard()

	log.Println("GameCreate finshed")
	returnBoardAsJson(w, board)
}

func GameById(w http.ResponseWriter, r *http.Request) {

	addResponseHeaders(w)

	vars := mux.Vars(r)
	gameId := vars["gameId"]

	board, err := models.LoadGameBoard(gameId)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	returnBoardAsJson(w, board)

}

func returnBoardsSummaryAsJson(w http.ResponseWriter, boardsSummary *[]models.GameBoardSummary) {

	json.NewEncoder(w).Encode(&boardsSummary)

}

func returnBoardAsJson(w http.ResponseWriter, board *models.GameBoard) {

	json.NewEncoder(w).Encode(&board)

}

func returnPlayerAsJson(w http.ResponseWriter, player *models.Player) {

	json.NewEncoder(w).Encode(&player)

}

// received MainPlayer as JSON request
func AddPlayer(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Add player started")
	addResponseHeaders(w)

	jsonBody, err := getRequestBody(r)
	if err != nil {
		http.Error(w, "Failed to get request body", http.StatusInternalServerError)
		return
	}

	// unmarshall Player request
	player, err := unmarshallPlayer(jsonBody)

	if err != nil {
		http.Error(w, "Failed to unmarshall player", http.StatusInternalServerError)
		return
	}

	// fetch current board
	vars := mux.Vars(r)
	gameId := vars["gameId"]

	fmt.Println("Getting game board", gameId)
	board, err := models.LoadGameBoard(gameId)

	if board == nil || err != nil {
		http.NotFound(w, r)
		return
	}

	player, err = board.AddPlayer(player)

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Save game board", gameId)
	err = board.SaveGameBoard()

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	returnPlayerAsJson(w, player)

}

// received MainPlayer as JSON request
func UpdatePlayer(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Update player started")
	addResponseHeaders(w)

	jsonBody, err := getRequestBody(r)
	if err != nil {
		http.Error(w, "Failed to get request body", http.StatusInternalServerError)
		return
	}

	// unmarshall Player request
	mainPlayer, err := unmarshallPlayer(jsonBody)

	if err != nil {
		http.Error(w, "Failed to unmarshall player", http.StatusInternalServerError)
		return
	}

	// fetch current board
	vars := mux.Vars(r)
	gameId := vars["gameId"]
	//playerId := vars["playerId"]

	fmt.Println("Getting game board", gameId)
	board, err := models.LoadGameBoard(gameId)

	if board == nil || err != nil {
		http.NotFound(w, r)
		return
	}

	err = board.MovePlayer(mainPlayer)

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Save game board", gameId)
	err = board.SaveGameBoard()

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	returnBoardAsJson(w, board)

}

func unmarshallPlayer(jsonBody []byte) (*models.Player, error) {

	var mainPlayer models.Player

	err := json.Unmarshal(jsonBody, &mainPlayer)

	return &mainPlayer, err

}
