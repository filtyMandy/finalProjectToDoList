package main

import (
	api "finalProjectToDoList/pkg/api/nextdate"
	"finalProjectToDoList/pkg/api/signin"
	"finalProjectToDoList/pkg/api/task"
	"finalProjectToDoList/pkg/api/task/done"
	"finalProjectToDoList/pkg/api/tasks"
	"finalProjectToDoList/pkg/db"
	"log"
	"net/http"
	"os"
)

func main() {
	// finding TODO_DBFILE **
	dbFile := "scheduler.db"
	if env := os.Getenv("TODO_DBFILE"); env != "" {
		dbFile = env
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatalf("Error DB connection: %v", err)
	}

	webDir := "web"
	port := "7540"

	// finding TODO_PORT *
	if env := os.Getenv("TODO_PORT"); env != "" {
		port = env
	}

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	http.HandleFunc("/api/signin", signin.SigninHandler)
	http.HandleFunc("/api/nextdate", signin.Auth(api.NextDateHandler))
	http.HandleFunc("/api/task", signin.Auth(task.TaskHandler))
	http.HandleFunc("/api/tasks", signin.Auth(tasks.GetTaskHandler))
	http.HandleFunc("/api/task/done", signin.Auth(done.DoneTaskHandler))

	log.Printf("Listening on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}
