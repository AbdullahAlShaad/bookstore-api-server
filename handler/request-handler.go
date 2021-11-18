package handler

import (
	"net/http"
)

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"strings"
)

func setJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	WriteJSONResponse(w,http.StatusOK,bookList)
}

func GetBooksNameSimplified(w http.ResponseWriter, r *http.Request) {

	var booksName []string
	for _,bookInfo := range bookList {
		booksName = append(booksName,bookInfo.BookName)
	}
	WriteJSONResponse(w,http.StatusOK,booksName)
}

func GetBookByName(w http.ResponseWriter, r *http.Request) {

	bookName := chi.URLParam(r,"bookName")

	for _,bookInfo := range bookList {
		bookNameWOSpace := strings.ReplaceAll(bookInfo.BookName," ","")
		if bookInfo.BookName == bookName || bookNameWOSpace == bookName {

			WriteJSONResponse(w,http.StatusOK,bookInfo)
			return
		}
	}

	WriteJSONResponse(w,http.StatusBadRequest,"Book not found. Check if spelling is correct")
}
func GetBookByISBN(w http.ResponseWriter, r *http.Request) {

	ISBN := chi.URLParam(r,"ISBN")

	book, ok := bookList[ISBN]
	if !ok {
		WriteJSONResponse(w,http.StatusNotFound,"Requested book doesn't exist")
		return
	}
	WriteJSONResponse(w,http.StatusOK,book)
}

func AddBook(w http.ResponseWriter, r *http.Request) {

	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)

	if err != nil {
		WriteJSONResponse(w,http.StatusBadRequest,"Can't Decode Given Data")
		return
	}

	if len(book.BookName) == 0 || len(book.ISBN) == 0 || len(book.AuthorInfo.Name) == 0 {
		WriteJSONResponse(w,http.StatusBadRequest,"Name, ISBN and Author Name are required field to add a book")
		return
	}

	if _,ok := bookList[book.ISBN]; ok == true {
		WriteJSONResponse(w,http.StatusConflict,"The book already exist")
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
	WriteJSONResponse(w,http.StatusOK,"Successfully Added Book")
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {

	ISBN := chi.URLParam(r,"ISBN")

	if _,ok := bookList[ISBN] ; !ok {
		WriteJSONResponse(w,http.StatusConflict,"Wrong ISBN")
		return
	}

	updatedBookInfo := bookList[ISBN]
	err := json.NewDecoder(r.Body).Decode(&updatedBookInfo)

	if err != nil {
		WriteJSONResponse(w,http.StatusBadRequest,"Can't Decode Given Data")
		return
	}

	if bookList[ISBN].AuthorInfo.Name != updatedBookInfo.AuthorInfo.Name || bookList[ISBN].BookName != updatedBookInfo.BookName || bookList[ISBN].ISBN != updatedBookInfo.ISBN || ISBN != updatedBookInfo.ISBN {
		WriteJSONResponse(w,http.StatusBadRequest,"Cant edit book name , author name or isbn. If there is mistake in these fields please delete the entry and try re adding")
		return
	}
	bookList[ISBN] = updatedBookInfo

	WriteJSONResponse(w,http.StatusOK,"Successfully Updated Book")
}

func removeFromSlice(s []string, index int) []string {
	//return append(s[:index], s[index+1:]...)
	s[index] = s[len(s)-1]
	return s[:len(s)-1]
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {

	ISBN := chi.URLParam(r,"ISBN")

	if _,ok := bookList[ISBN] ; ok == false {
		WriteJSONResponse(w,http.StatusBadRequest,"Wrong ISBN")
		return
	}

	bookInfo := bookList[ISBN]

	if cnt := authorBookCount[bookInfo.AuthorInfo.Name] ; cnt == 1 {
		delete(authorBookCount,bookInfo.AuthorInfo.Name)
		delete(authorList,bookInfo.AuthorInfo.Name)
	} else {
		authorBookCount[bookInfo.AuthorInfo.Name]--
		authorBooks := authorList[bookInfo.AuthorInfo.Name].Books

		for i,bookName := range authorBooks {
			if bookName == bookInfo.BookName {
				authorList[bookInfo.AuthorInfo.Name].Books = removeFromSlice(authorBooks,i)
				break
			}
		}
	}
	delete(bookList,ISBN)

	WriteJSONResponse(w,http.StatusOK,"Successfully Deleted Book")
}

func GetAllAuthors(w http.ResponseWriter, r *http.Request){

	var authorNames []string
	for authorName := range authorBookCount {
		authorNames = append(authorNames,authorName)
	}
	WriteJSONResponse(w, http.StatusOK, authorNames)
}

func GetAuthorInfo(w http.ResponseWriter, r *http.Request) {
	authorName := chi.URLParam(r, "AuthorName")

	if _, ok := authorList[authorName]; !ok {
		WriteJSONResponse(w,http.StatusNotFound,"Author Does Not Exist")
		return
	}
	WriteJSONResponse(w,http.StatusOK,authorList[authorName])
}

func WriteJSONResponse(w http.ResponseWriter, code int, data interface{}) {

	setJSONHeader(w)
	w.WriteHeader(code)
	switch x := data.(type) {
	case string:
		w.Write([]byte(x))
	case []byte:
		w.Write(x)
	default:
		err := json.NewEncoder(w).Encode(x)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}