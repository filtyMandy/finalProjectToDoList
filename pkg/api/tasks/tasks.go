package tasks

import (
	"encoding/json"
	"finalProjectToDoList/pkg/api/task"
	"finalProjectToDoList/pkg/db"
	"net/http"
	"time"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func getTaskHandlerCasual(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50)
	if err != nil {
		task.WriteJSONError(w, err.Error())
		return
	}
	resp := TasksResp{Tasks: tasks}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	limit := 50

	var tasks []*db.Task
	var err error

	if search != "" {
		date, errDate := time.Parse("02.01.2006", search)
		if errDate == nil {
			tasks, err = db.FindTaskByDate(date.Format("20060102"), limit)
		} else {
			like := "%" + search + "%"
			tasks, err = db.FindTasksBySubstring(like, limit)
		}
	} else {
		tasks, err = db.Tasks(limit)
	}
	if err != nil {
		task.WriteJSONError(w, err.Error())
		return
	}
	resp := TasksResp{Tasks: tasks}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func Init() {
	http.HandleFunc("/api/tasks", getTaskHandler)
}
