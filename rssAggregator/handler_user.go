package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/robinmuhia/GolangProjects/rssAggregator/internal/database"
)

func (apiConfig *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request){
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w,400,fmt.Sprintf("Error parsing JSON:%v",err))
		return
	}

	user,err := apiConfig.DB.CreateUser(r.Context(),database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
	})
	if err != nil{
		respondWithError(w,400,fmt.Sprintf("Couldn't create user:%v",err))
		return
	}
	respondWithJson(w,200,databaseUserToUser(user))
}