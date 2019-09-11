package db

import (
	"database/sql"
	"github.com/korolev1307/news_site/types"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

func ParseAllImagesByID(id int) []string {
	var allFiles []string
	var filesdir = "./files/" + strconv.Itoa(int(id)) + "/images/"
	files, err := ioutil.ReadDir(filesdir)
	if err != nil {
		log.Println(err)
	}
	for _, file := range files {
		filename := "files/" + strconv.Itoa(int(id)) + "/images/" + file.Name()
		allFiles = append(allFiles, filename)
	}
	return allFiles
}

func DeleteImageByPath(path string) error {
	dir := "./" + path
	err := os.Remove(dir)
	return err
}

func DeleteNews(id int) error {
	query := "delete from news where id=?"
	err := QueryFunc(query, id)
	return err
}

func AdministrateNews(id int) error {
	query := "update news set approved_by_administrator=? where id=?"
	error := QueryFunc(query, 1, id)
	return error
}

func ModerateNews(id, moderator int) error {
	query := "update news set approved_by_moderator=?, moderated_by_id=?, moderated_at=? where id=?"
	t := time.Now()
	error := QueryFunc(query, 1, moderator, t, id)
	return error
}

func AddNewsDB(title, content, short_content, folder_name, login string, images, files, publishing_at_main_page, publishing_at_lit_page, publishing_at_EC int) error {
	log.Println("AddNews: started function")
	userID, err := GetUserID(login)

	if err != nil && (title != "" || content != "" || short_content != "") {
		return err
	}

	error := QueryFunc("insert into news(title, user_id, content, short_content, created_date, folder_name, images, files, publishing_at_main_page, publishing_at_lit_page, publishing_at_EC) values(?,?,?,?,datetime(),?,?,?,?,?,?)", title, userID, content, short_content, folder_name, images, files, publishing_at_main_page, publishing_at_lit_page, publishing_at_EC)
	return error
}

func UpdateNewsDB(id int, title, content, short_content, folder_name, login string, images, files, publishing_at_main_page, publishing_at_lit_page, publishing_at_EC int) error {
	log.Println("UpdateNews: started function")
	userID, err := GetUserID(login)

	if err != nil && (title != "" || content != "" || short_content != "") {
		return err
	}

	error := QueryFunc("update news set title=?, user_id=?, content=?, short_content=?, created_date=datetime(), folder_name=?, images=?, files=?, publishing_at_main_page=?, publishing_at_lit_page=?, publishing_at_EC=? where id=?", title, userID, content, short_content, folder_name, images, files, publishing_at_main_page, publishing_at_lit_page, publishing_at_EC, id)
	return error
}

func GetLastNewsId() (int, error) {
	query := "SELECT MAX(ID) AS LastID FROM news"
	var id int
	rows := database.query(query)
	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			log.Println(err)
			return 0, err
		}
	}
	return id, nil
}

func GetNewsById(id int) (types.News, error) {
	var News types.News
	var date string
	var moderator_date sql.NullString
	var user types.User
	querySQL := "select title, user_id, content, short_content, created_date, moderated_at, folder_name, images, approved_by_administrator, approved_by_moderator, publishing_at_main_page, publishing_at_lit_page, publishing_at_EC, moderated_by_id from news where id=?"
	rows := database.query(querySQL, id)
	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&News.Title, &News.User_id, &News.Content, &News.Short_content, &date, &moderator_date, &News.Folder_name, &News.Images, &News.Approved_by_administrator, &News.Approved_by_moderator, &News.Publishing_at_main_page, &News.Publishing_at_lit_page, &News.Publishing_at_EC, &News.Moderated_by_id)
		if err != nil {
			log.Println("Error in GetNewsById")
			log.Println(err)
		}

		log.Println(News.Title)
		News.Id = id
		user, _ = GetUserById(News.User_id)
		News.Author = user.Surname + " " + user.Name + " " + user.Patronymic

		t, _ := time.Parse("2006-01-02T15:04:05Z07:00", date)
		addedtime := t.Add(time.Hour * 3)
		News.Created_date = addedtime.Format("02-01-2006 15:04")

		News.All_Files = ParseAllImagesByID(News.Id)
		log.Println("Moderator:" + strconv.Itoa(int(News.Moderated_by_id.Int64)))
		if News.Moderated_by_id.Int64 != 0 {
			moderator, _ := GetUserById(int(News.Moderated_by_id.Int64))
			News.Moderator_name.String = moderator.Surname + " " + moderator.Name + " " + moderator.Patronymic

			f, _ := time.Parse("2006-01-02T15:04:05Z07:00", moderator_date.String)
			News.Moderated_at.String = f.Format("02-01-2006 15:04")
		}
	}
	return News, err
}

func GetAllNews() ([]types.News, error) {
	query := "select id, title, user_id, content, short_content, created_date, folder_name, images, approved_by_administrator, approved_by_moderator from news order by id desc"
	var NewsArray []types.News
	var News types.News
	var user types.User
	var date string
	rows := database.query(query)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&News.Id, &News.Title, &News.User_id, &News.Content, &News.Short_content, &date, &News.Folder_name, &News.Images, &News.Approved_by_administrator, &News.Approved_by_moderator)
		if err != nil {
			log.Println("Error in GetAllNews")
			log.Println(err)
		}
		t, _ := time.Parse("2006-01-02T15:04:05Z07:00", date)
		addedtime := t.Add(time.Hour * 3)
		News.Created_date = addedtime.Format("02-01-2006 15:04")
		user, _ = GetUserById(News.User_id)
		News.Author = user.Surname + " " + user.Name + " " + user.Patronymic
		News.Filename = ParseAllImagesByID(News.Id)[0]
		log.Println("Image name is:" + News.Filename)
		NewsArray = append(NewsArray, News)
	}
	return NewsArray, err
}
