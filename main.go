package main

import (
	"github.com/korolev1307/news_site/views"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {
	views.PopulateTemplates()
	http.HandleFunc("/", views.HomePage)
	http.HandleFunc("/login/", views.LoginPage)
	http.HandleFunc("/logout/", views.RequiresLogin(views.LogoutPage))
	http.HandleFunc("/signup/", views.SignUpPage)
	http.HandleFunc("/addnews/", views.AddNews)
	http.HandleFunc("/userlist/", views.UserListPage)
	http.HandleFunc("/news/", views.ShowNewsPage)
	http.HandleFunc("/edit/", views.EditNews)

	http.Handle("/static/", http.FileServer(http.Dir("public")))
	http.Handle("/files/", http.FileServer(http.Dir("")))

	http.HandleFunc("/edit/delete-image", views.DeleteImage)
	http.HandleFunc("/delete-news/", views.DeleteNews)
	http.HandleFunc("/approve/", views.ModerateNews)

	log.Println("Server is running...")
	http.ListenAndServe(":8888", nil)

}
