package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/robinmuhia/GolangProjects/gormBooks/pkg/routes"
)

func main(){
	r := mux.NewRouter()
	routes.RegisterBookStoreRouters(r)
	http.Handle("/",r)
	fmt.Println("Server started")
	log.Fatal(http.ListenAndServe("localhost:8080",r))
}
