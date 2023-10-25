package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)



var secret_key = os.Getenv("SESSION_SECRET")
var client_id = os.Getenv("CLIENT_ID")
var secret_id = os.Getenv("SECRET_CLIENT")

func main()  {
    key := secret_key
    maxAge := 86400 * 7
    isProd := false
    store := sessions.NewCookieStore([]byte(key))
    store.MaxAge(maxAge)
    store.Options.Path = "/"
    store.Options.HttpOnly = true
    store.Options.Secure = isProd
    gothic.Store = store
    
    goth.UseProviders(
        google.New(client_id , secret_id, "http://localhost:3000/auth/google/callback"),
    )
    p := pat.New()
    p.Get("/auth/{provider}/callback", func(w http.ResponseWriter, r *http.Request) {
    user, err := gothic.CompleteUserAuth(w, r)
    if err != nil {
        fmt.Fprintln(w, err)
        return
    }
    t, _ := template.ParseFiles("templates/success.html")
    t.Execute(w, user)
    })
    p.Get("/auth/{provider}", func(w http.ResponseWriter, r *http.Request) {
        gothic.BeginAuthHandler(w, r)
    })
    p.Get("/", func(w http.ResponseWriter, r *http.Request) {
        t, _ := template.ParseFiles("templates/index.html")
        t.Execute(w, false)
    })
    log.Println("Listening on port 3000")
    log.Fatal(http.ListenAndServe(":3000", p))
    
}
