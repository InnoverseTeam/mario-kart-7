package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Player struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Score int    `json:"score"`
}
type Server struct {
	players map[string]*Player
	mu      sync.Mutex
}
func newServer() *Server {
	return &Server{
		players: make(map[string]*Player),
	}
}
func (s *Server) addPlayer(w http.ResponseWriter, r *http.Request) {
	var player Player
	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	s.players[player.ID] = &player
	s.mu.Unlock()

	w.WriteHeader(http.StatusCreated)
}
func (s *Server) getPlayer(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	s.mu.Lock()
	player, ok := s.players[id]
	s.mu.Unlock()

	if !ok {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(player)
}
func (s *Server) updateScore(w http.ResponseWriter, r *http.Request) {
	var update struct {
		ID    string `json:"id"`
		Score int    `json:"score"`
	}

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	player, ok := s.players[update.ID]
	if ok {
		player.Score = update.Score
	}
	s.mu.Unlock()

	if !ok {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (s *Server) listPlayers(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	players := make([]*Player, 0, len(s.players))
	for _, player := range s.players {
		players = append(players, player)
	}
	s.mu.Unlock()

	json.NewEncoder(w).Encode(players)
}
func (s *Server) Start() {
	http.HandleFunc("/addPlayer", s.addPlayer)
	http.HandleFunc("/getPlayer", s.getPlayer)
	http.HandleFunc("/updateScore", s.updateScore)
	http.HandleFunc("/listPlayers", s.listPlayers)

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
