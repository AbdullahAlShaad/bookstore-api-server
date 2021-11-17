package Handler

import (
	"github.com/go-chi/jwtauth/v5"
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

var TokenAuth *jwtauth.JWTAuth

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

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

func addAuthorToList(authorName string , authorObject *Author) {
	authorNameWOSpace := strings.ReplaceAll(authorName," ","")

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


// Key : UserID , Value : Password
var UserDB map[string]string



func InitializeData()  {
	bookList = make(BookDB)
	authorList = make(AuthorDB)
	authorBookCount = make(map[string]int)
	GenerateDummyData()
	GenerateAuthorInfo()
	initAuth()

}
