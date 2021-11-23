# BookstoreAPIServer

A simple api server which can be used as backend of a e-bookstore. User can get list of books and authors. Details of a book can be found using GET method by passing book name or ISBN as parameter. Details about authors can be found using GET method by passing author name as parameter. After registration and login, User can Add/Update and Delete books from the server

# Installation
To install Go, run the following command:
```shell
$ go_version=1.17.1
$ cd ~/Downloads
$ sudo apt-get update
$ sudo apt-get install -y build-essential git curl wget
$ wget https://dl.google.com/go/go${go_version}.linux-amd64.tar.gz
$ sudo tar -C /usr/local -xzf go${go_version}.linux-amd64.tar.gz
$ sudo chown -R $(id -u):$(id -g) /usr/local/go
$ rm go${go_version}.linux-amd64.tar.gz
```
Add go to your $PATH variable
```shell
$ mkdir $HOME/go
$ nano ~/.bashrc
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:/usr/local/go/bin:$PATH
$ source ~/.bashrc
$ go version
```

## Running the Server
After changing the directory to project directory
`go run .`


## API Calls

|method|url|body|actions|
|---|---|---|---|
|POST|`http://localhost:8081/register`|```json
{ "username": "shaad","password":"1234"}```| Register new user|
|POST|`http://localhost:8081/login`|{```"username": "shaad","password":"1234"}```| Loggin In the user|
|POST|`http://localhost:8080/logout`||Logging Out|
|GET|`http://localhost:8081/books`||Returns Details of all books|
|GET|`http://localhost:8081/books/simple`||Return list of books name|
|GET|`http://localhost:8081/books/isbn/{ISBN}`||Returns Details of the Book of the given ISBN|
|GET|`http://localhost:8081/books/name{BookName}`||Return Details of the given book name|
|GET|`http://localhost:8081/authors`||Return list of authors|
|GET|`http://localhost:8081/authors/{AuthorName}`||Returns Details of the given author|
|DELETE|`http://localhost:8081/books/{ISBN}`||Delete the book entry matching ISBN|
|POST|`http://localhost:8081/books/`|Given Below|Adds a new book in the list|
|PUT|`http://localhost:8081/books/{ISBN}`|Given Below|Updates a book information|
-----------------

Request Body for POST and PUT methods to add and update book information. To update book information we must pass book_name, authorname and isbn correctly and only fill other fields which needs to be updated
```json
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
```



