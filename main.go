package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	s := r.PathPrefix("/User").Subrouter()
	// r.HandleFunc("/", home)
	r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("/static"))))

	s.HandleFunc("/", users).Methods("GET")
	s.HandleFunc("/{id}", user).Methods("GET")
	s.HandleFunc("/new", AddUser).Methods("POST")

	fmt.Println("Starting Server...")
	log.Fatal(http.ListenAndServe(":80", r))
}

func home(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("index.html")
	if err != nil {
		panic(err)
	}

	temp.Execute(w, "hello world!")

	fmt.Fprint(w, "Home")
}

func users(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Users")
}

func user(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "User: %v\n", vars["id"])
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "New User")
}
