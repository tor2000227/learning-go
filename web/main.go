package main

import (
	"cmd1/pkg/models/mysql"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
    infoLog  *log.Logger
    errorLog *log.Logger
    snippets *mysql.SnippetModel
    templateCache map[string]*template.Template
}

func main() {
    addr := flag.String("addr", ":4000", "HTTP Network Address")
    dsn := flag.String("dsn", "", "MySQL DSN (required)")
    flag.Parse()
    
    infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
    errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

    db, err := openDB(*dsn)
    if err != nil {
        errorLog.Fatal(err)
    }
    defer db.Close()

    templateCache, err := newTemplateCache("./ui/html")
    if err != nil {
        errorLog.Fatal(err)
        
    }

    app := &application{
        errorLog: errorLog,
        infoLog:  infoLog,
        snippets: &mysql.SnippetModel{DB: db},
        templateCache: templateCache,
    }

    srv := &http.Server{
        Addr:     *addr,
        ErrorLog: errorLog,
        Handler:  app.routes(),
    }

    infoLog.Printf("Starting server on %s", *addr)
    err = srv.ListenAndServe()
    errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }
    if err = db.Ping(); err != nil {
        return nil, err
    }
    return db, nil
}
func newTemplateCache(dir string) (map[string]*template.Template, error) {
    cache := map[string]*template.Template{}
    pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
    if err != nil {
        return nil, err
    }
    for _, page := range pages{
        name := filepath.Base(page)
        ts, err := template.ParseFiles(filepath.Join(dir, "base.layout.tmpl"))
        if err != nil {
            return nil, err
        }
        ts, err = ts.ParseFiles(page)
        if err != nil {
            return nil, err
        }
        partials, err := filepath.Glob(filepath.Join(dir, "*.partial.tmpl"))
        if err != nil {
            return nil, err
        }
        for _, partial := range partials {
            ts, err = ts.ParseFiles(partial)
            if err != nil {
                return nil, err
            }
        }
        cache[name] = ts
    }
    return cache, nil
}