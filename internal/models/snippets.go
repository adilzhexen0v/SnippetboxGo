package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *pgxpool.Pool
}

func (sm *SnippetModel) Insert(w http.ResponseWriter, snippet Snippet) (int, error) {
	conn, err := sm.DB.Acquire(context.Background())
	if err != nil {
		fmt.Printf("Unable to acquire a database connection: %v\n", err)
		w.WriteHeader(500)
		return 0, nil
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		"INSERT INTO snippets (title, content, created, expires) VALUES ($1, $2, current_timestamp, current_timestamp + interval '365 days') RETURNING id",
		snippet.Title, snippet.Content)
	var id int
	err = row.Scan(&id)
	if err != nil {
		fmt.Printf("Unable to INSERT: %v\n", err)
		w.WriteHeader(500)
		return 0, nil
	}

	return id, err
}

func (sm *SnippetModel) Get(w http.ResponseWriter, snippetId int) (*Snippet, error) {
	conn, err := sm.DB.Acquire(context.Background())
	if err != nil {
		fmt.Printf("Unable to acquire a database connection: %v\n", err)
		w.WriteHeader(500)
		return nil, nil
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		"SELECT id, title, content, created, expires FROM snippets WHERE expires > current_timestamp AND id = $1",
		snippetId)
	s := &Snippet{}
	err = row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, err
}

func (sm *SnippetModel) Latest(w http.ResponseWriter) ([]*Snippet, error) {
	conn, err := sm.DB.Acquire(context.Background())
	if err != nil {
		fmt.Printf("Unable to acquire a database connection: %v\n", err)
		w.WriteHeader(500)
		return nil, nil
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(),
		"SELECT id, title, content, created, expires FROM snippets WHERE expires > current_timestamp ORDER BY id DESC LIMIT 10",
	)
	snippets := []*Snippet{}
	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			fmt.Println("ErrNoRecord from snippets.go")
			return nil, ErrNoRecord
		} else {
			fmt.Println("Another error from snippets.go")

			return nil, err
		}
	}
	return snippets, err
}
