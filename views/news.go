package views

import (
	//"html/template"
	"io"
	"net/http"
	"os"
	"log"
	"strconv"
	"github.com/korolev1307/news_site/db"
	"github.com/korolev1307/news_site/sessions"
	"github.com/korolev1307/news_site/types"

	"github.com/disintegration/imaging"
)

func AddNews(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//GET displays the upload form.
	case "GET":
		var context types.Context
        context.LoggedIn = sessions.IsLoggedIn(r)
        context.CurrentName, context.CurrentPatronymic = db.GetUserNameAndPatronymic(sessions.GetCurrentUserLogin(r)) 
        if context.LoggedIn {
        	addnewsTemplate.Execute(w, context)
        } else {
			http.Redirect(w, r, "/login/", 302)
		}
		
	//POST takes the uploaded file(s) and saves it to disk.
	case "POST":
		//parse the multipart form in the request
		err := r.ParseMultipartForm(1000000)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//get a ref to the parsed multipart form
		m := r.MultipartForm
		r.ParseForm()
		title := r.Form.Get("title")
		content := r.Form.Get("content")
		publishing_at_main_page, _ := strconv.Atoi(r.Form.Get("publishing_at_main_page"))
		publishing_at_lit_page, _ := strconv.Atoi(r.Form.Get("publishing_at_lit_page"))
		publishing_at_EC, _ := strconv.Atoi(r.Form.Get("publishing_at_EC"))
		LastNewsid, _ := db.GetLastNewsId()
		LastNewsid++
		login := sessions.GetCurrentUserLogin(r)
		//log.Print(title + " " + content + " " + publishing_at_main_page + publishing_at_lit_page +publishing_at_EC)
		//get the *fileheaders
		files := m.File["myfiles"]
		var filesbool int
		if len(files) > 0 {
			filesbool = 1
		} else {
			filesbool = 0
		}
		path := "files/"  + strconv.Itoa(int(LastNewsid)) + "/images/"
		

		for i, _ := range files {
			//for each fileheader, get a handle to the actual file
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				log.Print("Error in first step")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//create destination file making sure the path is writeable.
			errorfile := os.MkdirAll(path,os.ModePerm)

			if errorfile != nil {
 				log.Println("Error creating directory")
 				log.Println(errorfile)
 				return
			}
			dst, err := os.Create(path + files[i].Filename)
			defer dst.Close()
			if err != nil {
				log.Print("Error in second step")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//copy the uploaded file to the destination file
			if _, err := io.Copy(dst, file); err != nil {
				log.Print("Error in third step")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ResizeImage(path + files[i].Filename)
		}

		errordb := db.AddNewsDB(title, content, path, login, filesbool, filesbool, publishing_at_main_page, publishing_at_lit_page, publishing_at_EC)
		if errordb != nil {
			http.Error(w, "Unable to add news", http.StatusInternalServerError)
		} else {
			http.Redirect(w, r, "/", 302)
		}
	default:
		http.Redirect(w, r, "/addnews/", http.StatusInternalServerError)
	}
}


func ResizeImage(address string) {
	src, err := imaging.Open(address)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	src = imaging.Resize(src, 800, 0, imaging.Lanczos)

	err = imaging.Save(src, address)
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
}

//HomePage is used to handle the "/" URL 
//TODO add http404 error
func HomePage(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        var context types.Context
        context.LoggedIn = sessions.IsLoggedIn(r)
        context.CurrentName, context.CurrentPatronymic = db.GetUserNameAndPatronymic(sessions.GetCurrentUserLogin(r)) 
        homeTemplate.Execute(w, context)
        //homeTemplate.ExecuteTemplate(w, "home.html", r)
        message = ""
    } else {
        message = "Method not allowed"
        http.Redirect(w, r, "/", http.StatusFound)
    }
}

