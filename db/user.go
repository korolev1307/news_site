package db

import (
	"log"
	"github.com/korolev1307/news_site/types"
)

func CreateUser(name, surname, patronymic, login, password string) error {
	err := QueryFunc("insert into users(name, surname, patronumic, login, password) values(?,?,?,?,?)", name, surname, patronymic, login, password)
	return err
}

func ValidUser(username, password string) bool {
	var passwordFromDB string
	userSQL := "select password from users where login=?"
	log.Print("validating user ", username)
	rows := database.query(userSQL, username)

	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&passwordFromDB)
		if err != nil {
			return false
		}
	}
	//If the password matches, return true
	if password == passwordFromDB {
		return true
	}
	//by default return false
	return false
}

func GetUserID(username string) (int, error) {
	var userID int
	userSQL := "select id from users where login=?"
	rows := database.query(userSQL, username)

	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&userID)
		if err != nil {
			return -1, err
		}
	}
	return userID, nil
}

func GetUserNameAndPatronymic(username string) (string, string) {
	var name, patronymic string
	userSQL := "select name, patronumic from users where login=?"
	rows := database.query(userSQL, username)

	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&name, &patronymic)
		if err != nil {
			return "", ""
		}
	}
	return name, patronymic
}

func GetAllUsers() (types.Context, error) {
	var user types.User
	var users []types.User
	var context types.Context
	userSQL := "select id, name, surname, patronumic, login, administrator, moderator from users"
	rows := database.query(userSQL)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Name, &user.Surname, &user.Patronymic, &user.Login, &user.Administrator, &user.Moderator)
		if err != nil {
			log.Println(err)
		}
		users = append(users, user)
	}
	context = types.Context{Users: users}
	return context, nil
}

func UpdateUserRole(id, administrator, moderator int) error {
	err := QueryFunc("update users set administrator=?, moderator=? where id=?", administrator, moderator, id)
	return err
}