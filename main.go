package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

type BookDB map[string]Book
type AuthorDB map[string]int

// Keeps record Books
// key - ISBN, value -> Book Object

var bookList BookDB

// Records Number of books of author
// key -> author , value -> Number of books written by that author
var authorList AuthorDB


type Author struct {
	Name string `json:"name"`
	DateOfBirth string `json:"date_of_birth"`
	BirthPlace string `json:"birth_place"`
}

type Book struct {
	BookName string `json:"book_name"`
	Author `json:"author"`
	ISBN string `json:"isbn"`
	Genre string `json:"genre"`
	Publisher string `json:"publisher"`
}

func GenerateDummyData() {
	book1 := Book{
		BookName: "A Thousand Splendid Suns",
		Author : Author{
			Name: "Khaled Hosseini",
			DateOfBirth : "March 4, 1965",
			BirthPlace: "Afganistan",
		},
		ISBN: "978-1-59448-950-1",
		Genre: "Fiction",
		Publisher: "Riverhead Books",
	}
	bookList["978-1-59448-950-1"] = book1

	book2 := Book{
		BookName: "The Alchemist",
		Author : Author{
			Name: "Paulo Coelho",
			DateOfBirth : "August 24, 1947",
			BirthPlace: "Brazil",
		},
		ISBN: "0-06-250217-4",
		Genre: "Fiction",
		Publisher: "HarperTorch ",
	}
	bookList["0-06-250217-4"] = book2

	book3 := Book{
		BookName: "The Godfather",
		Author : Author{
			Name: "Mario Puzo",
			DateOfBirth : "October 15, 1920",
			BirthPlace: "United States",
		},
		ISBN: "13:9780399103421",
		Genre: "Crime Novel",
		Publisher: "G. P. Putnam's Sons",
	}
	bookList["13:9780399103421"] = book3

	book4 := Book{
		BookName: "The Kite Runner",
		Author : Author{
			Name: "Khaled Hosseini",
			DateOfBirth : "March 4, 1965",
			BirthPlace: "Afganistan",
		},
		ISBN: "1-57322-245-3",
		Genre: "Drama",
		Publisher: "Riverhead Books",
	}
	bookList["1-57322-245-3"] = book4
}

func initializeData()  {
	bookList = make(BookDB)
	authorList = make(AuthorDB)
	GenerateDummyData()

	for key,value := range bookList {
		fmt.Println("Key",key, " -> " , "Value : " , value)
	}

}

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	err := json.NewEncoder(w).Encode("Hello From Get All Books")
	if(err != nil) {
		http.Error(w,err.Error(),400)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	bookName := chi.URLParam(r,"bookName")

	 arr := [3]string {
		bookName,
		bookName,
		"Hello From Getbook",
	}

	err := json.NewEncoder(w).Encode(arr)

	if(err != nil) {
		http.Error(w,err.Error(),400)
	}
	w.WriteHeader(http.StatusOK)

}

func main() {

	initializeData()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Get("/",func(w http.ResponseWriter , r *http.Request) {
		w.Write([]byte("root."))
	})


	r.Route("/books",func(r chi.Router) {
		r.Get("/",GetAllBooks)
		r.Get("/{bookName}",GetBook)
	})


	http.ListenAndServe(":8081",r)
}