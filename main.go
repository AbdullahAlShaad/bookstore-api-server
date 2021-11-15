package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"strings"
)

type BookDB map[string]Book
type AuthorDB map[string] *Author

// Keeps record Books
// key - ISBN, value -> Book Object
var bookList BookDB

// Records Number of books of author
// key -> author , value -> Number of books written by that author
var authorList AuthorDB
var authorBookCount map[string]int


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

/*
{
    "book_name" : "Harry Potter",
    "author_info" : {
        "name" : "JK Rowling",
        "date_of_birth" : "31 July 1965",
        "birth_place" : "England"
    },
    "ISBN" : "0-7475-3269-9",
    "Genre" : "Fantasy",
    "Publisher" : "Bloomsbury"
}

{
    "book_name" : "The Sicilian",
    "author_info" : {
        "name" : "Mario Puzo",
        "date_of_birth" : "October 15, 1920",
        "birth_place" : "United States"
    },
    "ISBN" : "0-671-43564-7",
    "Genre" : "Thriller",
    "Publisher" : "	Random House"
}

 */

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

func addAuthorToList(authorName string , authorObject *Author) {
	authorNameWOSpace := strings.ReplaceAll(authorName," ","")

	//fmt.Println("Adding new obj" , authorName,authorNameWOSpace,authorObject.Books)

	authorList[authorName] = authorObject
	authorList[authorNameWOSpace] = authorObject
}

func addBookToAuthor(authorName string, bookName string) {
	authorList[authorName].Books = append(authorList[authorName].Books, bookName)
}

func GenerateAuthorInfo() {

	for _,bookInfo := range bookList {

		authorName := bookInfo.AuthorInfo.Name

		if _,ok := authorBookCount[authorName] ; ok == false {
			authorBookCount[authorName] = 1
			authorObj := &Author{AuthorInfo : bookInfo.AuthorInfo,
				Books : []string{bookInfo.BookName},}

			addAuthorToList(authorName,authorObj)

		} else {
			authorBookCount[authorName]++
			addBookToAuthor(authorName,bookInfo.BookName)
		}
	}
}

func initializeData()  {
	bookList = make(BookDB)
	authorList = make(AuthorDB)
	authorBookCount = make(map[string]int)
	GenerateDummyData()
	GenerateAuthorInfo()

}

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Call get all book")
	w.Header().Set("Content-Type","application/json")
	err := json.NewEncoder(w).Encode(bookList)
	if err != nil {
		http.Error(w,err.Error(),400)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetBooksNameSimplified(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	var booksName []string
	for _,bookInfo := range bookList {
		booksName = append(booksName,bookInfo.BookName)
	}
	err := json.NewEncoder(w).Encode(booksName)

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
			if err != nil {
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

func AddBook(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Called Add book method")

	w.Header().Set("Content-Type","application/json")
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)

	if err != nil {
		fmt.Println("Cant decode")
		http.Error(w,err.Error(),400)
		return
	}

	if _,ok := bookList[book.ISBN]; ok == true {
		return
	}

	bookList[book.ISBN] = book

	authorName := book.AuthorInfo.Name

	if _,ok := authorBookCount[authorName] ; ok == true {
		authorBookCount[authorName]++
		addBookToAuthor(authorName,book.BookName)
	} else {
		author := &Author{AuthorInfo :book.AuthorInfo ,Books: []string {book.BookName},}
		addAuthorToList(authorName,author)
	}
	w.WriteHeader(http.StatusOK)
}

func GetAllAuthors(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")

	var authorNames []string
	for authorName := range authorBookCount {
		authorNames = append(authorNames,authorName)
	}

	err := json.NewEncoder(w).Encode(authorNames)

	if err != nil {
		http.Error(w,err.Error(),400)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetAuthorInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")

	authorName := chi.URLParam(r,"AuthorName")

	if _,ok := authorList[authorName] ; ok == false {
		w.WriteHeader(404)
		return
	}

	err := json.NewEncoder(w).Encode( authorList[authorName])
	if err != nil {
		http.Error(w,err.Error(),400)
		return
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
		r.Get("/simplified",GetBooksNameSimplified)
		r.Get("/ISBN/{ISBN}",GetBookByISBN)
		r.Post("/AddBook",AddBook)
	})

	r.Route("/authors",func(r chi.Router) {
		r.Get("/",GetAllAuthors)
		r.Get("/{AuthorName}",GetAuthorInfo)
	})




	http.ListenAndServe(":8081",r)
}