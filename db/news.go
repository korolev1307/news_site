package db

import (
	"log"
	//"github.com/korolev1307/news_site/types"
)

func AddNewsDB(title, content, folder_name, login string, images, files, publishing_at_main_page, publishing_at_lit_page, publishing_at_EC int) error {
	log.Println("AddNews: started function")
	userID, err := GetUserID(login)

	if err != nil && (title != "" || content != "") {
		return err
	}

	error := QueryFunc("insert into news(title, user_id, content, created_date, folder_name, images, files, publishing_at_main_page, publishing_at_lit_page, publishing_at_EC) values(?,?,?,datetime(),?,?,?,?,?,?)", title, userID, content, folder_name, images, files, publishing_at_main_page, publishing_at_lit_page, publishing_at_EC)
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
			//send email to respective people
		}
	}
	return id, nil
}

// func GetAllNews() (types.Context, error) {
// 	query := "select id, title, user_id, content, created_date, folder_name, images, approved_by_administrator, approved_by_moderator from news order by id asc"

// }
