package main

import (
	"forum/internal/app"
	"forum/internal/handlers"
	"forum/internal/service"
	"forum/internal/sqlite"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Add a templateCache field to the application struct.

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := sqlite.NewRepo("./data/database.db")
	if err != nil {
		log.Fatal(err)
	}
	templateCache, err := app.NewTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}
	service := service.NewService(db)
	app := app.New(infoLog, errorLog, templateCache)
	handlers := handlers.New(service, app)
	infoLog.Println("Server is running on :http://localhost:4000")
	srv := &http.Server{
		Addr:     ":4000",
		ErrorLog: app.ErrorLog,
		Handler:  handlers.Routes(),
	}
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
