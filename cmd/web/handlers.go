package main

import (
	"com.aitu.snippetbox/internal/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	snippets, err1 := app.snippets.Latest(w)
	if err1 != nil {
		app.serverError(w, err1)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)
	}
	/*
		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			"./ui/html/pages/home.tmpl",
		}
		ts, err := template.ParseFiles(files...)

		if err != nil {
			app.serverError(w, err)
			return
		}

		err = ts.ExecuteTemplate(w, "base", nil)
		if err != nil {
			app.serverError(w, err)
		}

	*/
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(w, id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	fmt.Fprintf(w, "%+v", snippet)

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	snippet1 := models.Snippet{
		Title:   "Test title",
		Content: "Test content",
	}

	id, err := app.snippets.Insert(w, snippet1)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
