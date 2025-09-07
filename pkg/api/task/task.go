package task

import (
	"encoding/json"
	"finalProjectToDoList/pkg/db"
	"finalProjectToDoList/pkg/util.go"
	"log"
	"net/http"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Printf("[DeleteTask] id not specified")
		util_go.WriteJSONError(w, http.StatusBadRequest, "ID not specified")
		return
	}
	err := db.DeleteTask(id)
	if err != nil {
		log.Printf("[DeleteTask] failed to delete id %s: %v", id, err)
		util_go.WriteJSONError(w, http.StatusBadRequest, "error DeleteTask")
		return
	}
	util_go.WriteJSON(w, map[string]interface{}{})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		util_go.WriteJSONError(w, http.StatusBadRequest, "ID not specified")
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		util_go.WriteJSONError(w, http.StatusBadRequest, "error GetTask")
		return
	}
	util_go.WriteJSON(w, task)
}

func putTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		util_go.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	err = db.ValidateTask(&t)
	if err != nil {
		util_go.WriteJSONError(w, http.StatusBadRequest, "error ValidateTask")
		return
	}

	err = db.UpdateTask(&t)
	if err != nil {
		util_go.WriteJSONError(w, http.StatusBadRequest, "error UpdateTask")
		return
	}
	util_go.WriteJSON(w, map[string]interface{}{})
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
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
