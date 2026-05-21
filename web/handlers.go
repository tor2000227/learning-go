package main

import (
    "fmt"
    
    "net/http"
    "strconv"

    "cmd1/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        app.notFound(w)
        return
    }
    
    snippets, err := app.snippets.Latest()
    if err != nil {
        app.serverError(w, err)
        return
    }
    app.render(w,r, "home.page.tmpl", &templateData{
        Snippets: snippets,
    })
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil || id < 1 {
        app.notFound(w)
        return
    }

    snippet, err := app.snippets.Get(id)
    if err == models.ErrNoRecord {
        app.notFound(w)
        return
    } else if err != nil {
        app.serverError(w, err)
        return
    }
    app.render(w, r, "show.page.tmpl", &templateData{
        Snippet: snippet,
    })
    
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.Header().Set("Allow", http.MethodPost)
        app.clientError(w, http.StatusMethodNotAllowed)
        return
    }
    
    title := "New Snippet"
    content := "This is the content of my new snippet"
    expires := 7
    
    id, err := app.snippets.Insert(title, content, expires)
    if err != nil {
        app.serverError(w, err)
        return
    }
    
    http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}