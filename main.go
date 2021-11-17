package main

import (
	"github.com/Shaad7/BookstoreAPIServer/Handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)


func main() {
	Handler.InitializeData()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Get("/",func(w http.ResponseWriter , r *http.Request) {
		w.Write([]byte("root."))
	})

	r.Post("/login", Handler.Login)
	r.Post("/logout", Handler.Logout)
	r.Post("/register",Handler.Register)

	r.Group(func(r chi.Router) {

		r.Route("/books",func(r chi.Router) {

				r.Get("/",Handler.GetAllBooks)
				r.Get("/name/{bookName}", Handler.GetBookByName)
				r.Get("/simple",Handler.GetBooksNameSimplified)
				r.Get("/isbn/{ISBN}", Handler.GetBookByISBN)

				r.Group(func(r chi.Router) {

					r.Use(jwtauth.Verifier(Handler.TokenAuth))
					r.Use(jwtauth.Authenticator)

					r.Post("/",Handler.AddBook)
					r.Put("/{ISBN}",Handler.UpdateBook)
					r.Delete("/{ISBN}",Handler.DeleteBook)
				})

		})

		r.Route("/authors",func(r chi.Router) {
			r.Get("/",Handler.GetAllAuthors)
			r.Get("/{AuthorName}",Handler.GetAuthorInfo)
		})
	})


	http.ListenAndServe(":8081",r)
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
    "Publisher" : "Random House"
}

*/