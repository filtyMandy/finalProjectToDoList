package task

import (
	"encoding/json"
	"finalProjectToDoList/pkg/db"
	"net/http"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		WriteJSON(w, map[string]string{"error": "ID not specified"})
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		WriteJSON(w, map[string]string{"error": err.Error()})
		return
	}
	WriteJSON(w, task)
}

func putTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		WriteJSON(w, map[string]string{"error": "Invalid JSON"})
		return
	}

	err = db.ValidateTask(&t)
	if err != nil {
		WriteJSON(w, map[string]string{"error": err.Error()})
		return
	}

	err = db.UpdateTask(&t)
	if err != nil {
		WriteJSON(w, map[string]string{"error": err.Error()})
		return
	}
	WriteJSON(w, map[string]interface{}{})
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:

		addTaskHandler(w, r)
	case http.MethodPut:
		putTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func Init() {
	http.HandleFunc("/api/task", taskHandler)
}
