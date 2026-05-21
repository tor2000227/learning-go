package main

import (
    "fmt"
    "html/template"
    "net/http"
    "runtime/debug"
    "strings"
    "time"
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

func truncate(s string, n int) string {
    if len(s) <= n {
        return s
    }
    return s[:n] + "..."
}
func formatDate(t time.Time) string {
    return t.Format("Jan 02, 20006 at 15:04")

}
func formatDateShort(t time.Time) string {
    return t.Format("2006-01-02")
}

func isExpired(expires time.Time) bool {
    return expires.Before(time.Now())
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
    ts, ok := app.templateCache[name]
    if !ok {
        app.errorLog.Printf("Template %s not found in cache", name)
        app.serverError(w, fmt.Errorf("the template %s does not exist", name))
        return
    }
    clone, err := ts.Clone()
    if err != nil {
        app.serverError(w, err)
        return
    }

    clone = clone.Funcs(template.FuncMap{
        "truncate":     truncate,
        "formatDate":   formatDate,
        "formatDateShort":  formatDateShort,
        "isExpired":    isExpired,
        "add":          func(a,b int) int {return a + b },
        "subtract":     func(a,b int) int {return a - b },
        "safeHtml": func(s string) template.HTML {
            return template.HTML(s)
        },
        "nl2br": func(s string) template.HTML {
            return template.HTML(strings.ReplaceAll(s, "\n", "<br>"))
        },
        "timeSince": func(t time.Time) string {
            duration := time.Since(t)
            if duration.Hours() > 24 {
                days := int(duration.Hours() / 24)
                return fmt.Sprintf("%d days ago", days)
            }
            if duration.Hours() > 1 {
                return fmt.Sprintf("%.of hours ago", duration.Hours)
            }
            if duration.Minutes() > 1 {
                return fmt.Sprintf("%.of minutes ago", duration.Minutes())
            }
            return "just now"
        },
    })

    err = clone.ExecuteTemplate(w, "base", td)
    if err != nil {
        app.errorLog.Printf("Failed to execute template: %v", err)
        app.serverError(w, err)
    }
}

