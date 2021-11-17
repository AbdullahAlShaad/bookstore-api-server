package Handler

import "net/http"

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"strings"
)
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

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")

	ISBN := chi.URLParam(r,"ISBN")

	if _,exist := bookList[ISBN] ; exist == false {
		w.WriteHeader(404)
		return
	}

	updatedBookInfo := bookList[ISBN]
	err := json.NewDecoder(r.Body).Decode(&updatedBookInfo)

	if err != nil {
		fmt.Println("Cant decode")
		http.Error(w,err.Error(),400)
		return
	}

	if bookList[ISBN].AuthorInfo.Name != updatedBookInfo.AuthorInfo.Name || bookList[ISBN].BookName != updatedBookInfo.BookName || bookList[ISBN].ISBN != updatedBookInfo.ISBN || ISBN != updatedBookInfo.ISBN {
		json.NewEncoder(w).Encode("Cant edit book name , author name or isbn. If there is mistake in these fields please delete the entry and try re adding")
		w.WriteHeader(400)
		return
	}
	bookList[ISBN] = updatedBookInfo
	w.WriteHeader(http.StatusOK)

}

func removeFromSlice(s []string, index int) []string {
	s[index] = s[len(s)-1]
	return s[:len(s)-1]
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	ISBN := chi.URLParam(r,"ISBN")

	if _,ok := bookList[ISBN] ; ok == false {
		w.WriteHeader(400)
		return
	}

	bookInfo := bookList[ISBN]

	if cnt,_ := authorBookCount[bookInfo.AuthorInfo.Name] ; cnt == 1 {
		delete(authorBookCount,bookInfo.AuthorInfo.Name)
		delete(authorList,bookInfo.AuthorInfo.Name)
	} else {
		authorBookCount[bookInfo.AuthorInfo.Name]--
		authorBooks := authorList[bookInfo.AuthorInfo.Name].Books

		for i,bookName := range(authorBooks) {
			if bookName == bookInfo.BookName {
				authorList[bookInfo.AuthorInfo.Name].Books = removeFromSlice(authorBooks,i)
				break
			}
		}

	}
	delete(bookList,ISBN)
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
	w.Header().Set("Content-Type", "application/json")

	authorName := chi.URLParam(r, "AuthorName")

	if _, ok := authorList[authorName]; ok == false {
		w.WriteHeader(404)
		return
	}

	err := json.NewEncoder(w).Encode(authorList[authorName])
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(http.StatusOK)

}
