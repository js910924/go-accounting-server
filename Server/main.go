package main

import (
	"log"
	"net/http"
	"server/app"

	_ "github.com/go-sql-driver/mysql"
)

var server app.Server

func main() {
	server.Init()

	// log.SetPrefix("")
	log.SetFlags(log.Lshortfile)

	log.Println("Starting Server...")
	log.Fatal(http.ListenAndServe(":80", server.Router))
}

// func Login(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodGet:
// 		renderTemplate(w, "Login", "Login Page")

// 	case http.MethodPost:
// 		account := r.FormValue("account")
// 		password := r.FormValue("password")
// 		fmt.Fprintf(w, "Account: %s, Password: %s", account, password)
// 	}
// }
