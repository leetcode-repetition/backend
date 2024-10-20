package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"sync"
)

var challenge_solutions sync.Map

func generate_challenge() (string, string) {
	token := make([]byte, 16)
	rand.Read(token)
	challenge_token := base64.StdEncoding.EncodeToString(token)

	a, _ := rand.Int(rand.Reader, big.NewInt(100))
	b, _ := rand.Int(rand.Reader, big.NewInt(100))
	c, _ := rand.Int(rand.Reader, big.NewInt(100))

	challenge := fmt.Sprintf("return (%d * %d * %d);", a, b, c)
	solution := fmt.Sprintf("%d", a.Int64()*b.Int64()*c.Int64())
	challenge_solutions.Store(challenge_token, solution)

	return challenge_token, challenge
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("EXTENSION_ORIGIN"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Challenge-Token, X-Challenge-Response")
		w.Header().Set("Access-Control-Expose-Headers", "X-Challenge, X-Challenge-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		token := r.Header.Get("X-Challenge-Token")
		solution_attempt := r.Header.Get("X-Challenge-Response")

		if token != "" && solution_attempt != "" {
			if solution, ok := challenge_solutions.Load(token); ok {
				if solution_attempt == solution.(string) {
					challenge_solutions.Delete(token)
					next.ServeHTTP(w, r)
					fmt.Println("Successfully completed challenge")
					return
				}
			}
			http.Error(w, "Invalid challenge solution: UNAUTHORIZED", http.StatusUnauthorized)
			return
		}

		token, challenge := generate_challenge()
		w.Header().Set("X-Challenge-Token", token)
		w.Header().Set("X-Challenge", challenge)

		http.Error(w, "Challenge required: UNAUTHORIZED", http.StatusUnauthorized)
	}
}
