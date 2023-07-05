package main

import (
	"fmt"
	"net/http"

	"github.com/robinmuhia/GolangProjects/rssAggregator/internal/auth"
	"github.com/robinmuhia/GolangProjects/rssAggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiConfig *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		apiKey,err := auth.GetAPIkey(r.Header)
		if err != nil {
			respondWithError(w,403,fmt.Sprintf("Auth error:%v",err))
			return
		}
		user,err := apiConfig.DB.GetUserByApiKey(r.Context(),apiKey)
	
		if err != nil {
			respondWithError(w,400,fmt.Sprintf("Couldn't get user:%v",err))
			return
		}
		handler(w,r,user)
	}
}