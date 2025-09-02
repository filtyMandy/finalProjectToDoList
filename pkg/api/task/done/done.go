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
		task.WriteJSONError(w, "ID not specified")
		return
	}

	// getting task
	t, err := db.GetTask(id)
	if err != nil {
		task.WriteJSONError(w, "error GetTask")
		return
	}

	if strings.TrimSpace(t.Repeat) == "" {
		// one time - delete
		err = db.DeleteTask(id)
		if err != nil {
			task.WriteJSONError(w, "error GetTask")
			return
		}
		task.WriteJSON(w, map[string]interface{}{})
		return
	}

	// counting next date
	next, err := api.NextDate(time.Now(), t.Date, t.Repeat)
	if err != nil {
		task.WriteJSONError(w, "error NextDate")
		return
	}
	err = db.UpdateDate(id, next)
	if err != nil {
		task.WriteJSONError(w, "error UpdateDate")
		return
	}
	task.WriteJSON(w, map[string]interface{}{})
}

func Init() {
	http.HandleFunc("/api/task/done", doneTaskHandler)
}
