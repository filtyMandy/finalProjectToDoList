package db

import (
	"database/sql"
	"fmt"
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

// show tasks by search

func FindTasksBySubstring(substring string, limit int) ([]*Task, error) {
	rows, err := DB.Query(`SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`,
		substring, substring, limit)
	if err != nil {
		return nil, fmt.Errorf("FindTasksBySubstring: %s", err)
	}
	defer rows.Close()
	var tasks []*Task
	for rows.Next() {
		t := &Task{}
		var id int64
		err := rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, fmt.Errorf("FindTasksBySubstring(Scan): %s", err)
		}
		t.ID = fmt.Sprintf("%d", id)
		tasks = append(tasks, t)
	}
	if tasks == nil {
		tasks = []*Task{}
	}
	return tasks, nil
}

func FindTaskByDate(date string, limit int) ([]*Task, error) {
	rows, err := DB.Query(`SELECT * FROM scheduler WHERE date = ? LIMIT ?`, date, limit)
	if err != nil {
		return nil, fmt.Errorf("FindTaskByDate: %s", err)
	}
	defer rows.Close()
	var tasks []*Task
	for rows.Next() {
		t := &Task{}
		var id int64
		err := rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, fmt.Errorf("FindTaskByDate(Scan): %s", err)
		}
		t.ID = fmt.Sprintf("%d", id)
		tasks = append(tasks, t)
	}
	if tasks == nil {
		tasks = []*Task{}
	}
	return tasks, nil
}

// show tasks by list

func Tasks(limit int) ([]*Task, error) {
	rows, err := DB.Query(`SELECT id, date, title, comment, repeat FROM scheduler
		ORDER BY date ASC LIMIT ?`, limit)
	if err != nil {
		return nil, fmt.Errorf("DB error: %v", err)
	}
	defer rows.Close()
	var tasks []*Task
	for rows.Next() {
		t := &Task{}
		var id int64
		err := rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, fmt.Errorf("Scan error: %v", err)
		}
		t.ID = fmt.Sprintf("%d", id)
		tasks = append(tasks, t)
	}
	//if NO tasks - return empty slice
	if tasks == nil {
		tasks = []*Task{}
	}
	return tasks, nil
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
