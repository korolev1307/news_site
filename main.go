package main

import (
	"log"
	//"fmt"
	"net/http"
  _ "github.com/mattn/go-sqlite3" 
  //"github.com/korolev1307/news_site/db"
  "github.com/korolev1307/news_site/views"
)

func main() {


  //log.Println(db.SearchName(1))
  views.PopulateTemplates()
  
  http.HandleFunc("/", views.HomePage)
  http.HandleFunc("/login/", views.LoginPage)
  http.HandleFunc("/logout/",views.RequiresLogin(views.LogoutPage))
  http.HandleFunc("/signup/",views.SignUpPage)
  http.HandleFunc("/addnews/",views.AddNews)
  http.HandleFunc("/userlist/",views.UserListPage)


  http.Handle("/static/", http.FileServer(http.Dir("public")))
  log.Println("Server is running...")
  http.ListenAndServe(":8888", nil)
  
}


