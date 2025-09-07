package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"strings"
	"time"
	"unicode"
)

var DB *sql.DB //Connecting DB

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func DeleteTask(id string) error {
	res, err := DB.Exec("DELETE FROM scheduler WHERE id=?", id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("task does not exist")
	}
	return nil
}

func UpdateDate(id string, date string) error {
	_, err := DB.Exec("UPDATE scheduler SET date=? WHERE id=?", date, id)
	return err
}

func ValidateTask(t *Task) error {
	// is date
	if len(t.Date) != 8 || !onlyDigits(t.Date) {
		return fmt.Errorf("invalid date: %s", t.Date)
	}

	//  date exist
	_, err := time.Parse("20060102", t.Date)
	if err != nil {
		return fmt.Errorf("invalid date: %s", t.Date)
	}
	// title empty
	if strings.TrimSpace(t.Title) == "" {
		return fmt.Errorf("invalid title: %s", t.Title)
	}

	if !validateRepeat(t.Repeat) {
		return fmt.Errorf("invalid repeat: %s", t.Repeat)
	}

	return nil
}

func validateRepeat(r string) bool {
	return r == "" || strings.HasPrefix(r, "d ") || strings.HasPrefix(r, "w ")
}

func onlyDigits(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func UpdateTask(t *Task) error {
	query := `UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`
	res, err := DB.Exec(query, t.Date, t.Title, t.Comment, t.Repeat, t.ID)
	if err != nil {
		return err
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if cnt == 0 {
		return fmt.Errorf("task %s not found", t.ID)
	}
	return nil
}

func GetTask(id string) (*Task, error) {
	row := DB.QueryRow(`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id)
	t := Task{}
	var idDB int64
	err := row.Scan(&idDB, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("error getting task: %w", err)
	}
	t.ID = fmt.Sprintf("%d", idDB)
	return &t, nil
}

// show tasks by search

func FindTasksBySubstring(substring string, limit int) ([]*Task, error) {
	substringReq := "%" + substring + "%"
	rows, err := DB.Query(`SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`,
		substringReq, substringReq, limit)
	if err != nil {
		return nil, fmt.Errorf("FindTasksBySubstring: %w", err)
	}
	defer rows.Close()
	var tasks []*Task
	for rows.Next() {
		t := &Task{}
		var id int64
		err := rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, fmt.Errorf("FindTasksBySubstring(Scan): %w", err)
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
		return nil, fmt.Errorf("FindTaskByDate: %w", err)
	}
	defer rows.Close()
	var tasks []*Task
	for rows.Next() {
		t := &Task{}
		var id int64
		err := rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, fmt.Errorf("FindTaskByDate(Scan): %w", err)
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
		return nil, fmt.Errorf("DB error: %w", err)
	}
	defer rows.Close()
	var tasks []*Task
	for rows.Next() {
		t := &Task{}
		var id int64
		err := rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, fmt.Errorf("Scan error: %w", err)
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
