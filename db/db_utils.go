package db

import(
	"log"
	//"fmt"
	"database/sql"
  _ "github.com/mattn/go-sqlite3" 
)
var database Database
var taskStatus map[string]int
var err error

type Database struct {
	db *sql.DB
}

//Queryfunc encapsulates running multiple queries which don't do much things
func QueryFunc(sql string, args ...interface{}) error {
	log.Print("creating insert query")
	SQL := database.prepare(sql)
	tx := database.begin()
	_, err = tx.Stmt(SQL).Exec(args...)
	if err != nil {
		log.Println("Insert: ", err)
		tx.Rollback()
	} else {
		err = tx.Commit()
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("Commit successful")
	}
	return err
}

//Begins a transaction
func (db Database) begin() (tx *sql.Tx) {
	tx, err := db.db.Begin()
	if err != nil {
		log.Println(err)
		return nil
	}
	return tx
}

func (db Database) prepare(q string) (stmt *sql.Stmt) {
	stmt, err := db.db.Prepare(q)
	if err != nil {
		log.Println(err)
		return nil
	}
	return stmt
}

func (db Database) query(q string, args ...interface{}) (rows *sql.Rows) {
	rows, err := db.db.Query(q, args...)
	if err != nil {
		log.Println(err)
		return nil
	}
	return rows
}

func init() {
	database.db, err = sql.Open("sqlite3", "./news.db")
	taskStatus = map[string]int{"COMPLETE": 1, "PENDING": 2, "DELETED": 3}
	if err != nil {
		log.Fatal(err)
	}
}

func Close() {
	database.db.Close()
}


func SearchName(num int) (id int, name string) {
 //  var (
	// 	id int
	// 	name string
	// )
	db, err := sql.Open("sqlite3", "news.db")
    if err != nil {
      log.Fatal(err)
    }
    defer db.Close()
	rows, err := db.Query("select id, name from users where id = ?", num)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
  return id, name
}