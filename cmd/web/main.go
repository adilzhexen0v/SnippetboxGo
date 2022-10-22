package main

import (
	"com.aitu.snippetbox/internal/models"
	"context"
	"flag"
	"github.com/jackc/pgx/v4/pgxpool"
	"html/template"
	"log"
	"net/http"
	"os"
)

var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	dbURL := "postgres://postgres:123@localhost:5432/snippetbox"
	db, err1 := openDB(dbURL)
	if err1 != nil {
		errorLog.Fatal(err1)
	}
	defer db.Close()

	templateCache, err2 := newTemplateCache()
	if err2 != nil {
		errorLog.Fatal(err2)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	pool, err1 := pgxpool.Connect(context.Background(), dsn)
	infoLog.Printf("Connected!")
	return pool, err1
}
