package api

import (
	"net/http"
	"time"
)

func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	nowStr := r.FormValue("now")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DF, nowStr)
		if err != nil {
			http.Error(w, "bad NOW parameter", http.StatusBadRequest)
			return
		}
	}

	result, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(result))
}
