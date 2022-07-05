package handlers

import (
	"fmt"
	"net/http"
)

type EnsureAuth struct {
	handler http.Handler
}

type User struct {
	userid string
}

func (app *EnsureAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := GetAuthenticatedUser(r)
	if err != nil {
		http.Error(w, "please sign-in", http.StatusUnauthorized)
		return
	}

	fmt.Println(user.userid)
	//TODO: call the real handler, but how do we share the user?
	app.handler.ServeHTTP(w, r)
}

func NewEnsureAuth(handlerToWrap http.Handler) *EnsureAuth {
	return &EnsureAuth{handlerToWrap}
}

func (app *EnsureAuth) GetAuthenticatedUser(r *http.Request) (*User, error) {
	//validate the session token in the request,
	//fetch the session state from the session store,
	//and return the authenticated user
	//or an error if the user is not authenticated
	return &User{
		userid: "me",
	}, nil
}

func (app *EnsureAuth) UsersMeHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetAuthenticatedUser(r)
	if err != nil {
		http.Error(w, "please sign-in", http.StatusUnauthorized)
		return
	}
	fmt.Println(user.userid)
	//GET = respond with current user's profile
	//PATCH = update current user's profile
}
