package views

import (
	"io/ioutil"
	"os"
	"strconv"
	// "github.com/korolev1307/news_site/db"
	// "github.com/korolev1307/news_site/sessions"
	// "github.com/korolev1307/news_site/types"
	"html/template"
	"log"
	"strings"
)

var (
	homeTemplate     *template.Template
	signupTemplate   *template.Template
	addnewsTemplate  *template.Template
	userlistTemplate *template.Template
	// deletedTemplate   *template.Template
	// completedTemplate *template.Template
	loginTemplate *template.Template
	// editTemplate      *template.Template
	// searchTemplate    *template.Template
	templates *template.Template
	message   string
	//message will store the message to be shown as notification
	err error
)

//PopulateTemplates is used to parse all templates present in
//the templates folder
func PopulateTemplates() {
	var allFiles []string
	templatesDir := "./templates/"
	files, err := ioutil.ReadDir(templatesDir)
	if err != nil {
		log.Println(err)
		os.Exit(1) // No point in running app if templates aren't read
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".html") {
			allFiles = append(allFiles, templatesDir+filename)
		}
	}

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	templates, err = template.ParseFiles(allFiles...)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	homeTemplate = templates.Lookup("home.html")
	loginTemplate = templates.Lookup("login.html")
	signupTemplate = templates.Lookup("signup.html")
	addnewsTemplate = templates.Lookup("addnews.html")
	userlistTemplate = templates.Lookup("userlist.html")
	// deletedTemplate = templates.Lookup("deleted.html")
	// editTemplate = templates.Lookup("edit.html")
	// searchTemplate = templates.Lookup("search.html")
	// completedTemplate = templates.Lookup("completed.html")

}

func ParseAllImagesByID(id int) []string {
	var allFiles []string
	var filesdir = "./files/" + strconv.Itoa(int(id)) + "/"
	files, err := ioutil.ReadDir(filesdir)
	if err != nil {
		log.Println(err)
	}

	for _, file := range files {
		filename := file.Name()
		allFiles = append(allFiles, filename)
	}
	return allFiles
}
