package types

import (
	"database/sql"
)

type User struct {
	Id                   int
	Name                 string
	Surname              string
	Patronymic           string
	Login                string
	Snils                string
	Allowed_registration bool
	Administrator        bool
	Moderator            bool
}

type Context struct {
	LoggedIn          bool
	Users             []User
	CurrentName       string
	CurrentPatronymic string
	NewsArray         []News
	News              News
	User              User
}

type News struct {
	Id                        int
	Title                     string
	User_id                   int
	Author                    string
	Content                   string
	Short_content             string
	Created_date              string
	Moderated_at              sql.NullString
	Folder_name               string
	Images                    bool
	Files                     bool
	Approved_by_administrator bool
	Approved_by_moderator     bool
	Publishing_at_main_page   bool
	Publishing_at_lit_page    bool
	Publishing_at_EC          bool
	Moderated_by_id           sql.NullInt64
	Moderator_name            sql.NullString
	Filename                  string
	All_Files                 []string
}
