package done

import (
	api "finalProjectToDoList/pkg/api/nextdate"
	"finalProjectToDoList/pkg/api/task"
	"finalProjectToDoList/pkg/db"
	"net/http"
	"strings"
	"time"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		task.WriteJSON(w, map[string]string{"error": "ID not specified"})
		return
	}

	// getting task
	t, err := db.GetTask(id)
	if err != nil {
		task.WriteJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if strings.TrimSpace(t.Repeat) == "" {
		// one time - delete
		err = db.DeleteTask(id)
		if err != nil {
			task.WriteJSON(w, map[string]string{"error": err.Error()})
			return
		}
		task.WriteJSON(w, map[string]interface{}{})
		return
	}

	// counting next date
	next, err := api.NextDate(time.Now(), t.Date, t.Repeat)
	if err != nil {
		task.WriteJSON(w, map[string]string{"error": err.Error()})
		return
	}
	err = db.UpdateDate(id, next)
	if err != nil {
		task.WriteJSON(w, map[string]string{"error": err.Error()})
		return
	}
	task.WriteJSON(w, map[string]interface{}{})
}

func Init() {
	http.HandleFunc("/api/task/done", doneTaskHandler)
}
