package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/robinmuhia/GolangProjects/gormBooks/pkg/models"
	"github.com/robinmuhia/GolangProjects/gormBooks/pkg/utils"
)

var NewBook models.Book

func GetAllBooks(w http.ResponseWriter, r *http.Request){
	newBooks := models.GetAllBooks()
	res,_ := json.Marshal(newBooks)
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetBookById(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	bookId := vars["bookId"]
	ID, err := strconv.ParseInt(bookId,0,0)
	if err != nil{
		fmt.Println("Error while parsing")
	}
	bookDetails, _ := models.GetBookById(ID)
	res,_ := json.Marshal(bookDetails)
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func CreateBook(w http.ResponseWriter, r *http.Request){
	CreatedBook := &models.Book{}
	utils.ParseBody(r,CreatedBook)
	b := CreatedBook.CreateBook()
	res,_ := json.Marshal(b)
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func DeleteBook(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	bookId := vars["bookId"]
	ID, err := strconv.ParseInt(bookId,0,0)
	if err != nil{
		fmt.Println("Error while parsing")
	}
	bookDeleted := models.DeleteBook(ID)
	res,_ := json.Marshal(bookDeleted)
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusNoContent)
	w.Write(res)
}

func UpdateBook(w http.ResponseWriter, r *http.Request){
	var updateBook = &models.Book{}
	utils.ParseBody(r,updateBook)
	vars := mux.Vars(r)
	bookId := vars["bookId"]
	ID, err := strconv.ParseInt(bookId,0,0)
	if err != nil{
		fmt.Println("Error while parsing")
	}
	bookExists, db := models.GetBookById(ID)
	if updateBook.Name != ""{
		bookExists.Name =updateBook.Name
	}
	if updateBook.Publication != ""{
		bookExists.Publication =updateBook.Publication
	}
	if updateBook.Author != ""{
		bookExists.Author = updateBook.Author
	}
	db.Save(&bookExists)
	res,_ := json.Marshal(bookExists)
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
