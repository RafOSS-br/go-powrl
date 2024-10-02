package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"powrl/internal/pow"

	"github.com/gorilla/mux"
)

func main() {
	store, err := pow.NewSQLiteChallengeStore("pow.db")
	if err != nil {
		log.Fatal(err)
	}
	powSystem := pow.NewPoWSystem(store)

	r := mux.NewRouter()
	r.HandleFunc("/generate_challenge", generateChallengeHandler(powSystem)).Methods("GET")
	r.HandleFunc("/verify_solution", verifySolutionHandler(powSystem)).Methods("POST")

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", r)
}

func generateChallengeHandler(p *pow.PoWSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		algo := r.URL.Query().Get("algo")
		difficultyStr := r.URL.Query().Get("difficulty")
		difficulty, err := strconv.Atoi(difficultyStr)
		if err != nil || difficulty <= 0 {
			difficulty = 4 // Default difficulty
		}

		challenge, err := p.GenerateChallenge(algo, difficulty)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Response to client
		response := map[string]string{
			"challenge_id": challenge.ID,
			"algorithm":    challenge.Algorithm,
			"difficulty":   strconv.Itoa(challenge.Difficulty),
			"data":         challenge.Data, // Protocol definition here
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func verifySolutionHandler(p *pow.PoWSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ChallengeID string `json:"challenge_id"`
			Solution    string `json:"solution"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		isValid, err := p.VerifySolution(req.ChallengeID, req.Solution)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response := map[string]bool{
			"valid": isValid,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
