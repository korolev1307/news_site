package sessions

import (
	"github.com/gorilla/sessions"
	"net/http"
)

//Store the cookie store which is going to store session data in the cookie
var Store = sessions.NewCookieStore([]byte("secret-password"))
var session *sessions.Session

//IsLoggedIn will check if the user has an active session and return True
func IsLoggedIn(r *http.Request) bool {
	session, _ := Store.Get(r, "session")

	if session.Values["loggedin"] == "true" {
		return true
	}
	return false
}

func GetCurrentUserLogin(r *http.Request) string {
	session, err := Store.Get(r, "session")
	if session.Values["loggedin"] == "true" && err == nil {
		return session.Values["login"].(string)
	}
	return ""
}
