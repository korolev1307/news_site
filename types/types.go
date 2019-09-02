package types

type User struct {
	Id            int
	Name          string
	Surname       string
	Patronymic    string
	Login         string
	Administrator bool
	Moderator     bool
}

type Context struct {
	LoggedIn          bool
	Users             []User
	CurrentName       string
	CurrentPatronymic string
	NewsArray         []News
}

type News struct {
	Id                        int
	Title                     string
	User_id                   string
	Content                   string
	Created_date              string
	Moderated_at              string
	Folder_name               string
	Images                    bool
	Files                     bool
	Approved_by_administrator bool
	Approved_by_moderator     bool
	Publishing_at_main_page   bool
	Publishing_at_lit_page    bool
	Publishing_at_EC          bool
	Moderated_by_id           int
}
