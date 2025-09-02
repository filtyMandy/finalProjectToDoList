package task

import (
	"encoding/json"
	"finalProjectToDoList/pkg/db"
	"log"
	"net/http"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Printf("[DeleteTask] id not specified")
		WriteJSONError(w, "ID not specified")
		return
	}
	err := db.DeleteTask(id)
	if err != nil {
		log.Printf("[DeleteTask] failed to delete id %s: %v", id, err)
		WriteJSONError(w, "error DeleteTask")
		return
	}
	WriteJSON(w, map[string]interface{}{})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		WriteJSONError(w, "ID not specified")
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		WriteJSONError(w, "error GetTask")
		return
	}
	WriteJSON(w, task)
}

func putTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		WriteJSONError(w, "Invalid JSON")
		return
	}

	err = db.ValidateTask(&t)
	if err != nil {
		WriteJSONError(w, "error ValidateTask")
		return
	}

	err = db.UpdateTask(&t)
	if err != nil {
		WriteJSONError(w, "error UpdateTask")
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
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func Init() {
	http.HandleFunc("/api/task", taskHandler)
}
