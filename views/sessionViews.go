package views

import (
	"github.com/korolev1307/news_site/db"
	"github.com/korolev1307/news_site/sessions"
	"github.com/korolev1307/news_site/types"
	"log"
	"net/http"
)

func RequiresLogin(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !sessions.IsLoggedIn(r) {
			http.Redirect(w, r, "/login/", 302)
			return
		}
		handler(w, r)
	}
}

//LogoutFunc Implements the logout functionality. WIll delete the session information from the cookie store
func LogoutPage(w http.ResponseWriter, r *http.Request) {
	session, err := sessions.Store.Get(r, "session")
	if err == nil { //If there is no error, then remove session
		if session.Values["loggedin"] != "false" {
			session.Values["loggedin"] = "false"
			session.Save(r, w)
			log.Print("Successfully logged out")
		}
	}
	http.Redirect(w, r, "/login", 302) //redirect to login irrespective of error or not
}

//LoginFunc implements the login functionality, will add a cookie to the cookie store for managing authentication
func LoginPage(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "session")

	switch r.Method {
	case "GET":
		var context types.Context
		context.LoggedIn = sessions.IsLoggedIn(r)
		context.CurrentName, context.CurrentPatronymic = db.GetUserNameAndPatronymic(sessions.GetCurrentUserLogin(r))
		if context.LoggedIn {
			http.Redirect(w, r, "/", 302)
		} else {
			loginTemplate.Execute(w, context)
		}
	case "POST":
		log.Print("Inside POST")
		r.ParseForm()
		login := r.Form.Get("login")
		password := r.Form.Get("password")

		if (login != "" && password != "") && db.ValidUser(login, password) {
			session.Values["loggedin"] = "true"
			session.Values["login"] = login
			session.Save(r, w)
			log.Print("user ", login, " is authenticated")
			http.Redirect(w, r, "/", 302)
			return
		}
		log.Print("Invalid user " + login)
		loginTemplate.Execute(w, nil)
	default:
		http.Redirect(w, r, "/login/", http.StatusUnauthorized)
	}
}

func SignUpPage(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		var context types.Context
		context.LoggedIn = sessions.IsLoggedIn(r)
		if context.LoggedIn {
			http.Redirect(w, r, "/", 302)
		} else {
			signupTemplate.Execute(w, context)
		}

	case "POST":
		r.ParseForm()
		name := r.Form.Get("name")
		surname := r.Form.Get("surname")
		patronumic := r.Form.Get("patronumic")
		login := r.Form.Get("login")
		password := r.Form.Get("password")
		snils := r.Form.Get("snils")
		log.Println(name, surname, patronumic, login, password, snils)

		err := db.CreateUser(name, surname, patronumic, login, password, snils)
		if err != nil {
			http.Error(w, "Unable to sign user up", http.StatusInternalServerError)
		} else {
			http.Redirect(w, r, "/", 302)
		}
	}
}
