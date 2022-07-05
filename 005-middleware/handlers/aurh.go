package handlers

import (
	"fmt"
	"net/http"
)

type UsersHandler struct {
	handler http.Handler
}

type User struct {
	userid string
}

func NewUsersHandler(h http.Handler) *UsersHandler {
	return &UsersHandler{
		handler: h,
	}
}

func GetAuthenticatedUser(r *http.Request) (*User, error) {
	//validate the session token in the request,
	//fetch the session state from the session store,
	//and return the authenticated user
	//or an error if the user is not authenticated
	return &User{
		userid: "me",
	}, nil
}

func UsersMeHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetAuthenticatedUser(r)
	if err != nil {
		http.Error(w, "please sign-in", http.StatusUnauthorized)
		return
	}
	fmt.Println(user.userid)
	//GET = respond with current user's profile
	//PATCH = update current user's profile
}
