package views

import (
	"github.com/korolev1307/news_site/db"
	"github.com/korolev1307/news_site/sessions"
	"log"
	"net/http"
	"strconv"
)

func UserListPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		context, _ := db.GetAllUsers()
		context.LoggedIn = sessions.IsLoggedIn(r)
		login := sessions.GetCurrentUserLogin(r)
		user_id, _ := db.GetUserID(login)
		context.User, _ = db.GetUserById(user_id)

		if !(context.LoggedIn && context.User.Administrator) {
			http.Redirect(w, r, "/", http.StatusFound)
		}
		context.CurrentName, context.CurrentPatronymic = db.GetUserNameAndPatronymic(sessions.GetCurrentUserLogin(r))
		log.Println(context)
		userlistTemplate.Execute(w, context)
		message = ""
	case "POST":
		log.Print("Inside POST")
		r.ParseForm()
		administrator, _ := strconv.Atoi(r.Form.Get("administrator"))
		moderator, _ := strconv.Atoi(r.Form.Get("moderator"))
		allowed_registration, _ := strconv.Atoi(r.Form.Get("allowed_registration"))
		id, _ := strconv.Atoi(r.Form.Get("id"))
		log.Println(id, administrator, moderator)

		err := db.UpdateUserRole(id, administrator, moderator, allowed_registration)
		if err != nil {
			http.Error(w, "Unable to update user roles", http.StatusInternalServerError)
		} else {
			log.Println("User successfully updated")
			http.Redirect(w, r, "/userlist/", 302)
		}
	default:
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
