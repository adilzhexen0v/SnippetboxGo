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

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl", data)

}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	snippet, err1 := app.snippets.Get(w, id)
	if err1 != nil {
		if errors.Is(err1, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err1)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl", data)

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	snippet1 := models.Snippet{
		Title:   "O snail",
		Content: "O snail\\nClimb Mount Fuji,\\nBut slowly, slowly!\\n\\n– Kobayashi Issa",
	}

	id, err := app.snippets.Insert(w, snippet1)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
