package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Moyaz79/Blog-Aggregator/internal/database"
	"github.com/Moyaz79/Blog-Aggregator/internal/database/auth"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWtihError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID: 	uuid.New(),
		CreatedAt: 	time.Now().UTC(),
		UpdatedAt: 	time.Now().UTC(),
		Name: params.Name,
	})
	if err != nil {
		respondWtihError(w, 400, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	respondWtihJSON(w, 201, databaseUserToUser(user))
}


func (apiCfg *apiConfig) getUserHandler(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWtihError(w, 403, fmt.Sprintf("Auth error: %v", err))
		return
	}

	user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil {
		respondWtihError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
		return
	}

	respondWtihJSON(w, 200, databaseUserToUser(user))
}
