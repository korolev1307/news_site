package views

import (
	//"html/template"
	"github.com/korolev1307/news_site/db"
	"github.com/korolev1307/news_site/sessions"
	"github.com/korolev1307/news_site/types"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

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
			http.Redirect(w, r, "/login/", http.StatusBadRequest)
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
		short_content := r.Form.Get("short_content")
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
		path := "files/" + strconv.Itoa(int(LastNewsid)) + "/images/"

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
			errorfile := os.MkdirAll(path, os.ModePerm)

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

		errordb := db.AddNewsDB(title, content, short_content, path, login, filesbool, filesbool, publishing_at_main_page, publishing_at_lit_page, publishing_at_EC)
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
		log.Fatalf("failed to open image: %v, %s", err, address)
	}
	src = imaging.Fit(src, 800, 600, imaging.Lanczos)

	err = imaging.Save(src, address)
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	} else {
		log.Println("Image resized successful")
	}
}

//HomePage is used to handle the "/" URL
//TODO add http404 error
func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var context types.Context
		context.LoggedIn = sessions.IsLoggedIn(r)
		context.CurrentName, context.CurrentPatronymic = db.GetUserNameAndPatronymic(sessions.GetCurrentUserLogin(r))
		context.NewsArray, _ = db.GetAllNews()

		login := sessions.GetCurrentUserLogin(r)
		user_id, _ := db.GetUserID(login)
		context.User, _ = db.GetUserById(user_id)

		homeTemplate.Execute(w, context)
		//homeTemplate.ExecuteTemplate(w, "home.html", r)
		message = ""
	} else {
		message = "Method not allowed"
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func ShowNewsPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	var context types.Context
	id, err := strconv.Atoi(r.URL.Path[len("/news/"):])
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	context.LoggedIn = sessions.IsLoggedIn(r)
	context.CurrentName, context.CurrentPatronymic = db.GetUserNameAndPatronymic(sessions.GetCurrentUserLogin(r))
	context.News, err = db.GetNewsById(id)

	login := sessions.GetCurrentUserLogin(r)
	user_id, _ := db.GetUserID(login)
	context.User, _ = db.GetUserById(user_id)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	shownewsTemplate.Execute(w, context)
}

func EditNews(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//GET displays the upload form.
	case "GET":
		var context types.Context
		context.LoggedIn = sessions.IsLoggedIn(r)
		context.CurrentName, context.CurrentPatronymic = db.GetUserNameAndPatronymic(sessions.GetCurrentUserLogin(r))
		login := sessions.GetCurrentUserLogin(r)
		id, _ := db.GetUserID(login)
		context.User, _ = db.GetUserById(id)

		news_id, err := strconv.Atoi(r.URL.Path[len("/edit/"):])
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}
		context.News, err = db.GetNewsById(news_id)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}
		log.Println(context.User)
		log.Println(context.News.User_id)
		if context.LoggedIn && (context.User.Id == context.News.User_id || context.User.Administrator || context.User.Moderator) {
			editnewsTemplate.Execute(w, context)
		} else {
			http.Redirect(w, r, "/", http.StatusBadRequest)
		}

	case "POST":
		log.Println("Inside POST LUL")
		err := r.ParseMultipartForm(1000000)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//get a ref to the parsed multipart form
		m := r.MultipartForm
		r.ParseForm()
		News_id, _ := strconv.Atoi(r.Form.Get("news_id"))
		title := r.Form.Get("title")
		content := r.Form.Get("content")
		short_content := r.Form.Get("short_content")
		publishing_at_main_page, _ := strconv.Atoi(r.Form.Get("publishing_at_main_page"))
		publishing_at_lit_page, _ := strconv.Atoi(r.Form.Get("publishing_at_lit_page"))
		publishing_at_EC, _ := strconv.Atoi(r.Form.Get("publishing_at_EC"))
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
		path := "files/" + strconv.Itoa(int(News_id)) + "/images/"
		if len(files) > 0 {
			for i, _ := range files {
				//for each fileheader, get a handle to the actual file
				file, err := files[i].Open()
				defer file.Close()
				if err != nil {
					log.Print("Error in first step")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				// //open current file path dir
				// errorfile := os.Chdir(path)

				// if errorfile != nil {
				// 	log.Println("Error open directory")
				// 	log.Println(errorfile)
				// 	return
				// }

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
		}

		errordb := db.UpdateNewsDB(News_id, title, content, short_content, path, login, filesbool, filesbool, publishing_at_main_page, publishing_at_lit_page, publishing_at_EC)
		if errordb != nil {
			http.Error(w, "Unable to update news", http.StatusInternalServerError)
		} else {
			http.Redirect(w, r, "/news/"+strconv.Itoa(int(News_id)), 302)
		}

	default:
		http.Redirect(w, r, "/", http.StatusInternalServerError)
	}
}

func DeleteImage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		filename := r.FormValue("filepath")
		log.Println(filename)
		err := db.DeleteImageByPath(filename)
		if err != nil {
			log.Println(err)
		}
	}
}

func DeleteNews(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	var context types.Context
	news_id, err := strconv.Atoi(r.URL.Path[len("/delete-news/"):])
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	context.News, _ = db.GetNewsById(news_id)
	login := sessions.GetCurrentUserLogin(r)
	id, _ := db.GetUserID(login)
	context.User, _ = db.GetUserById(id)
	log.Println(context.User.Id)
	log.Println(context.News.User_id)
	if (context.User.Id == context.News.User_id) || context.User.Administrator || context.User.Moderator {
		error := db.DeleteNews(news_id)
		if error != nil {
			log.Println(error)
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/", 302)
	} else {
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		log.Println("No access")
	}
}

func ModerateNews(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	var context types.Context
	news_id, err := strconv.Atoi(r.URL.Path[len("/approve/"):])
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	context.News, _ = db.GetNewsById(news_id)
	login := sessions.GetCurrentUserLogin(r)
	id, _ := db.GetUserID(login)
	context.User, _ = db.GetUserById(id)
	if context.User.Administrator {
		error := db.AdministrateNews(news_id)
		if error != nil {
			log.Println(error)
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/", 302)
	} else if context.User.Moderator {
		error := db.ModerateNews(news_id, context.User.Id)
		if error != nil {
			log.Println(error)
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/", 302)
	} else {
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		log.Println("No access")
	}
}
