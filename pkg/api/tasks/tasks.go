package tasks

import (
	"encoding/json"
	"finalProjectToDoList/pkg/db"
	"finalProjectToDoList/pkg/util.go"
	"net/http"
	"time"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	limit := 50

	var tasks []*db.Task
	var err error

	if search != "" {
		date, errDate := time.Parse("02.01.2006", search)
		if errDate == nil {
			tasks, err = db.FindTaskByDate(date.Format("20060102"), limit)
		} else {
			tasks, err = db.FindTasksBySubstring(search, limit)
		}
	} else {
		tasks, err = db.Tasks(limit)
	}
	if err != nil {
		util_go.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	resp := TasksResp{Tasks: tasks}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
