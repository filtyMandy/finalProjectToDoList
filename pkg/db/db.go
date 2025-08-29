package db

import (
	"database/sql"
	"log"
	_ "modernc.org/sqlite"
	"os"
)

var DB *sql.DB //Connecting DB

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// SQL table
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT,
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);
`

// Opening DB, creating if not exist
func Init(dbfile string) error {
	_, err := os.Stat(dbfile)
	var install bool
	if err != nil {
		install = true // if DB not exist
	}
	// Opening DB
	DB, err = sql.Open("sqlite", dbfile)
	if err != nil {
		return err
	}
	if install {
		_, err = DB.Exec(schema)
		if err != nil {
			log.Printf("Creating table failed: %v", err)
			return err
		}
	}
	return nil
}

func AddTask(task Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
