package Handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
	"time"
)

func initDummyUser() {
	UserDB = map[string]string {
		"User1": "Password1",
		"User2": "Password2",
		"user": "pass",
	}
}

func initAuth() {

	TokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
	initDummyUser()
}

func Login(w http.ResponseWriter, r *http.Request) {

	var cred Credential
	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		WriteJSONResponse(w,http.StatusBadRequest,"Can't Decode Given Data")
		return
	}

	if _,ok := UserDB[cred.Username] ; ok == false {
		WriteJSONResponse(w,http.StatusBadRequest,"Given User does not exist")
		return
	}

	if pass,_ := UserDB[cred.Username] ; pass != cred.Password {
		WriteJSONResponse(w,http.StatusUnauthorized,"Wrong Password")
		return
	}

	expireTime := time.Now().Add(10 * time.Minute)
	_,tokenString,err := TokenAuth.Encode(map[string]interface{}{
		"aud": "Abdullah Al Shaad",
		"exp": expireTime.Unix(),
	})

	fmt.Println(tokenString)

	if err != nil {
		WriteJSONResponse(w,http.StatusInternalServerError,"Can not Generate Auth Token")
		return
	}

	http.SetCookie(w,&http.Cookie{
		Name : "jwt",
		Value : tokenString,
		Expires: expireTime,
	})

	WriteJSONResponse(w,http.StatusOK,"Successfully Logged In... Welcome")

}
func Logout(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w,&http.Cookie{
		Name:    "jwt",
		Value:   "",
		Expires: time.Now(),
	})
	WriteJSONResponse(w,http.StatusOK,"Successfully Logged Out... Good Bye")
}


func Register(w http.ResponseWriter, r *http.Request) {

	var cred Credential
	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		WriteJSONResponse(w,http.StatusBadRequest,"Can't Decode Given Data")
		return
	}

	if len(cred.Username) == 0 || len(cred.Password) == 0 {
		WriteJSONResponse(w,http.StatusBadRequest,"Username and Password can not be empty")
		return
	}

	if _,ok := UserDB[cred.Username] ; ok == true {
		WriteJSONResponse(w,http.StatusConflict,"Can't Register. Username Taken")
		return
	}

	UserDB[cred.Username] = cred.Password
	WriteJSONResponse(w,http.StatusOK,"Registration Successful. Login to continue")
}