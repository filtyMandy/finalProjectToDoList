package done

import (
	api "finalProjectToDoList/pkg/api/nextdate"
	"finalProjectToDoList/pkg/db"
	"finalProjectToDoList/pkg/util.go"
	"log"
	"net/http"
	"strings"
	"time"
)

func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		util_go.WriteJSONError(w, http.StatusBadRequest, "ID not specified")
		return
	}

	// getting task
	t, err := db.GetTask(id)
	if err != nil {
		log.Printf("GetTask failed: %v", err)
		util_go.WriteJSONError(w, http.StatusBadRequest, "error GetTask")
		return
	}

	if strings.TrimSpace(t.Repeat) == "" {
		// one time - delete
		err = db.DeleteTask(id)
		if err != nil {
			log.Printf("DeleteTask failed: %v", err)
			util_go.WriteJSONError(w, http.StatusBadRequest, "error GetTask")
			return
		}
		util_go.WriteJSON(w, map[string]interface{}{})
		return
	}

	// counting next date
	next, err := api.NextDate(time.Now(), t.Date, t.Repeat)
	if err != nil {
		log.Printf("NextDate failed: %v", err)
		util_go.WriteJSONError(w, http.StatusBadRequest, "error NextDate")
		return
	}
	err = db.UpdateDate(id, next)
	if err != nil {
		log.Printf("UpdateDate failed: %v", err)
		util_go.WriteJSONError(w, http.StatusBadRequest, "error UpdateDate")
		return
	}
	util_go.WriteJSON(w, map[string]interface{}{})
}
