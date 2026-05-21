package main

import (
    "fmt"
    "html/template"
    "net/http"
    "runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
    trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
    app.errorLog.Output(2, trace)
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
    http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
    app.clientError(w, http.StatusNotFound)
}
func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
    files := []string{
        fmt.Sprintf("ui/html/%s", name),
        "ui/html/base.layout.tmpl",
        "ui/html/footer.partial.tmpl",
    }
 /*   ts, ok := app.templateCache[name]
    if !ok {
        app.serverError(w, fmt.Errorf("The template %s does not exist", name))
        return
    }
    err := ts.Execute(w, td)
    if err != nil {
        app.serverError(w, err)
    }
*/
     ts, err := template.ParseFiles(files...)
    if err != nil {
        app.serverError(w, err)
        return
    }
    
    err = ts.ExecuteTemplate(w, "base", td)
    if err != nil {
        app.serverError(w, err)
    }        
}