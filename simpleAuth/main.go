package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
//	"github.com/gorilla/securecookie"
	"gopkg.in/boj/redistore.v1"
)



//var secret_key = securecookie.GenerateRandomKey(8)

var store, err = redistore.NewRediStore(10, "tcp", ":6379", "", []byte("hello"))

var users = map[string]string{"amsem": "pass", "admin": "admin"}

func LoginHandler(w http.ResponseWriter, r *http.Request)  {
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Please pass data as URL fro; encoded ", http.StatusBadRequest)
        return
    }
    username := r.PostForm.Get("username")
    password := r.PostForm.Get("password")
    if originalPass, ok := users[username]; ok {
        session, _ := store.Get(r, "session.id")
        if password == originalPass {
            session.Values["authenticated"] = true
            session.Save(r, w)
        }else {
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
            return
        }
    }else {
        http.Error(w, "User is not found", http.StatusNotFound)
        return
    }
    w.Write([]byte("Logged in SUCCESSFULLY"))
}

func LogOutHandler(w http.ResponseWriter, r *http.Request)  {
    session, _ := store.Get(r, "session.id")
    session.Options.MaxAge = -1
    session.Save(r, w)
    w.Write([]byte("logged out"))
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request)  {
    session, _ := store.Get(r, "session.id")
    if (session.Values["authenticated"] != nil) && session.Values["authenticated"] != false {
        w.Write([]byte(time.Now().String()))
    }else {
        http.Error(w, "Forbidden", http.StatusForbidden)
    }
}

func main()  {
    r := mux.NewRouter()
//    fs := http.FileServer(http.Dir("./static/"))
//    r.PathPrefix("/static").Handler(http.StripPrefix("/static/", fs))
    r.HandleFunc("/login", LoginHandler)
    r.HandleFunc("/logout", LogOutHandler)
    r.HandleFunc("/health", HealthCheckHandler)
    http.Handle("/", r)
    srv := http.Server{
        Addr: ":8000",
        WriteTimeout: 15 * time.Second,
        ReadTimeout: 15 * time.Second,
    }
    log.Fatal(srv.ListenAndServe())


}


