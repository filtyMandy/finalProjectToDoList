package task

import (
	"encoding/json"
	api "finalProjectToDoList/pkg/api/nextdate"
	"finalProjectToDoList/pkg/db"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const DF = "20060102"

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		WriteJSON(w, map[string]string{"error": "Error JSON: " + err.Error()})
		return
	}
	t.Title = strings.TrimSpace(t.Title)
	if t.Title == "" {
		WriteJSON(w, map[string]string{"error": "Title is empty"})
		return
	}

	err = checkDate(&t)
	if err != nil {
		WriteJSON(w, map[string]string{"error": "Date format is wrong"})
		return
	}

	id, err := db.AddTask(t)
	if err != nil {
		WriteJSON(w, map[string]string{"error": err.Error()})
		return
	}

	WriteJSON(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}

func checkDate(task *db.Task) error {
	now := time.Now()
	dateString := strings.TrimSpace(task.Date)
	if dateString == "" {
		task.Date = now.Format("20060102")
		dateString = task.Date
	}
	t, err := time.Parse("20060102", dateString)
	if err != nil {
		return fmt.Errorf("date format error: %s", dateString)
	}
	if strings.TrimSpace(task.Repeat) == "" {
		if t.Before(dateOnly(now)) {
			task.Date = now.Format("20060102") // автоматом подставляем сегодня
		}
	} else {
		next, err := api.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("The repetition rule is incorrect: %v", err)
		}
		if t.Before(dateOnly(now)) {
			task.Date = next
		}
	}
	return nil
}

func dateOnly(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
