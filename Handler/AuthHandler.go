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

func Login(w http.ResponseWriter, r *http.Request) {
	setJSONHeader(w)

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

	encodeErr := json.NewEncoder(w).Encode("Succesfully Logged IN...")
	if encodeErr != nil {
		http.Error(w,encodeErr.Error(),404)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func Logout(w http.ResponseWriter, r *http.Request) {

	setJSONHeader(w)
	http.SetCookie(w,&http.Cookie{
		Name : "jwt",
		Expires: time.Now(),
	})

	encodeErr := json.NewEncoder(w).Encode("Succesfully Logged Out")
	if encodeErr != nil {
		http.Error(w,encodeErr.Error(),http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}


func Register(w http.ResponseWriter, r *http.Request) {

	setJSONHeader(w)
	var cred Credential
	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		fmt.Println("Can't Decode")
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}

	if len(cred.Username) == 0 || len(cred.Password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _,ok := UserDB[cred.Username] ; ok == true {
		fmt.Println("Can't Register. Username taken")
		w.WriteHeader(http.StatusConflict)
		return
	}

	UserDB[cred.Username] = cred.Password

	encodeErr := json.NewEncoder(w).Encode("Registration Successful")
	if encodeErr != nil {
		http.Error(w,encodeErr.Error(),404)
		return
	}
	w.WriteHeader(http.StatusOK)

}