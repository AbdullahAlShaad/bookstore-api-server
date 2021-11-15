package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"strings"
)

type BookDB map[string]Book
type AuthorDB map[string]Author

// Keeps record Books
// key - ISBN, value -> Book Object
var bookList BookDB

// Records Number of books of author
// key -> author , value -> Number of books written by that author
var authorList AuthorDB


type AuthorInfo struct {
	Name string `json:"name"`
	DateOfBirth string `json:"date_of_birth"`
	BirthPlace string `json:"birth_place"`
}

type Author struct {
	AuthorInfo `json:"author"`
	Books []string `json:"books"` //List of books the author wrote
}

type Book struct {
	BookName string `json:"book_name"`
	AuthorInfo `json:"author_info"`
	ISBN string `json:"isbn"`
	Genre string `json:"genre"`
	Publisher string `json:"publisher"`
}

func GenerateDummyData() {
	book1 := Book{
		BookName: "A Thousand Splendid Suns",
		AuthorInfo : AuthorInfo{
			Name: "Khaled Hosseini",
			DateOfBirth : "March 4, 1965",
			BirthPlace: "Afganistan",
		},
		ISBN: "1",
		Genre: "Fiction",
		Publisher: "Riverhead Books",
	}
	bookList["1"] = book1

	book2 := Book{
		BookName: "The Alchemist",
		AuthorInfo : AuthorInfo{
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
		AuthorInfo : AuthorInfo{
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
		AuthorInfo : AuthorInfo{
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

}

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	err := json.NewEncoder(w).Encode(bookList)
	if err != nil {
		http.Error(w,err.Error(),400)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func GetBookByName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	bookName := chi.URLParam(r,"bookName")

	for _,bookInfo := range bookList {
		bookNameWOSpace := strings.ReplaceAll(bookInfo.BookName," ","")
		if bookInfo.BookName == bookName || bookNameWOSpace == bookName {
			err := json.NewEncoder(w).Encode(bookInfo)
			if(err != nil) {
				http.Error(w,err.Error(),400)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	err := json.NewEncoder(w).Encode("Book not found. Check if spelling is correct")

	if err != nil {
		http.Error(w,err.Error(),400)
	}
	w.WriteHeader(http.StatusOK)

}
func GetBookByISBN(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")

	ISBN := chi.URLParam(r,"ISBN")

	if _,ok := bookList[ISBN] ; ok == false {
		w.WriteHeader(404)
		return
	}

	err := json.NewEncoder(w).Encode(bookList[ISBN])
	if err != nil {
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
		r.Get("/Name/{bookName}", GetBookByName)
		r.Get("/ISBN/{ISBN}",GetBookByISBN)
	})


	http.ListenAndServe(":8081",r)
}