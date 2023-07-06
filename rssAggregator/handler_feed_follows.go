package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/robinmuhia/GolangProjects/rssAggregator/internal/database"
)

func (apiConfig *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User){
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w,400,fmt.Sprintf("Error parsing JSON:%v",err))
		return
	}

	feedFollow,err := apiConfig.DB.CreateFeedFollows(r.Context(),database.CreateFeedFollowsParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: params.FeedID,
	})
	if err != nil{
		respondWithError(w,400,fmt.Sprintf("Couldn't create feed follow:%v",err))
		return
	}
	respondWithJson(w,201,databaseFeedFollowToFeedFollow(feedFollow))
}


func (apiConfig *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User){
	feedFollows,err := apiConfig.DB.GetFeedFollows(r.Context(),user.ID)
	if err != nil{
		respondWithError(w,400,fmt.Sprintf("Couldn't create feed:%v",err))
		return
	}
	respondWithJson(w,200,databaseFeedFollowsToFeedFollows(feedFollows))
}

func (apiConfig *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User){
	FeedFollowIdStr := chi.URLParam(r, "feedFollowId")
	feedFollowId,err := uuid.Parse(FeedFollowIdStr)

	if err != nil{
		respondWithError(w,400,fmt.Sprintf("Couldn't parse feed follow id: %v",err))
	}
	err = apiConfig.DB.DeleteFeedFollow(r.Context(),database.DeleteFeedFollowParams{
		ID: feedFollowId,
		UserID: user.ID,
	})
	if err != nil{
		respondWithError(w,400,fmt.Sprintf("Couldn't delete feed follow:%v",err))
	}
	respondWithJson(w,204,struct{}{})
}


