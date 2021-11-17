package Handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
	"time"
)

func initDummyUser() {
	UserDB = make(map[string]string)
	UserDB["User1"] = "Password1"
	UserDB["User2"] = "Password2"
}

func initAuth() {

	TokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
	initDummyUser()
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")

	var cred Credential
	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		fmt.Println("Can't Decode")
		http.Error(w,err.Error(),400)
		return
	}

	if _,ok := UserDB[cred.Username] ; ok == false {
		fmt.Println("User does not exist")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if pass,_ := UserDB[cred.Username] ; pass != cred.Password {
		fmt.Println("Password does not match")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expireTime := time.Now().Add(10 * time.Minute)
	_,tokenString,err := TokenAuth.Encode(map[string]interface{}{
		"aud": "Abdullah Al Shaad",
		"exp": expireTime.Unix(),
	})

	fmt.Println(tokenString)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}



	http.SetCookie(w,&http.Cookie{
		Name : "jwt",
		Value : tokenString,
		Expires: expireTime,
	})
	w.WriteHeader(http.StatusOK)
}
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	http.SetCookie(w,&http.Cookie{
		Name : "jwt",
		Expires: time.Now(),
	})
	w.WriteHeader(http.StatusOK)
}

