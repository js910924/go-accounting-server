package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"server/middleware"
	"server/models"
	"strings"

	"database/sql"
	// "go-account-server/handler"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type app struct {
	Templates *template.Template
	DB        *sql.DB
	Router    *mux.Router
}

func (a *app) Init() {
	var err error
	a.Templates = template.Must(template.ParseFiles(
		"./static/public/index.html",
		"./static/public/login.html",
		"./static/public/register.html",
		"./static/public/user.html",
	))
	a.DB = ConnectDB()
	if err != nil {
		log.Fatal("Open DB fail")
	}

	a.Router = mux.NewRouter()
}

func ConnectDB() *sql.DB {
	db, err := sql.Open("mysql", "root:0924@/account")
	if err != nil {
		log.Fatal("Open DB fail")
	}

	return db
}

var a app

func main() {
	a.Init()

	a.Router.HandleFunc("/", index).Methods("GET")
	// user := a.Router.PathPrefix("/User").Subrouter()
	a.Router.HandleFunc("/User", UserHandler).Methods("GET")
	a.Router.HandleFunc("/User", AddUserHandler).Methods("POST")
	a.Router.HandleFunc("/Register", Register).Methods("GET")
	a.Router.HandleFunc("/Login", login).Methods("GET")
	a.Router.HandleFunc("/Login", checkLogin).Methods("POST")

	log.Println("Starting Server...")
	log.Fatal(http.ListenAndServe(":80", a.Router))
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	// renderTemplate(a.Templates, w, "user", nil)
	c, err := r.Cookie("user")
	if err != nil {
		c = &http.Cookie{
			Name:  "user",
			Value: "",
		}

		fmt.Fprint(w, "No Cookie")
		return
	}

	values := strings.Split(c.Value, " ")
	item := struct {
		Account  string
		Password string
	}{
		Account:  values[0],
		Password: values[1],
	}

	middleware.RenderTemplate(a.Templates, w, "user", item)
}

func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	account := r.FormValue("account")
	password := r.FormValue("password")

	query := fmt.Sprintf(`insert into User (Name, Account, Password) values ("%s", "%s", "%s")`, name, account, password)
	log.Println("[Query]", query)

	_, err := a.DB.Query(query)
	if err != nil {
		log.Println(err)
	}

	c := &http.Cookie{
		Name:  "user",
		Value: account + " " + password,
	}

	http.SetCookie(w, c)
	http.Redirect(w, r, "/User", http.StatusFound)
}

func index(w http.ResponseWriter, r *http.Request) {
	middleware.RenderTemplate(a.Templates, w, "index", "Home")
}

func Register(w http.ResponseWriter, r *http.Request) {
	middleware.RenderTemplate(a.Templates, w, "register", "Registe Page")
}

func login(w http.ResponseWriter, r *http.Request) {
	middleware.RenderTemplate(a.Templates, w, "login", "Login")
}

func checkLogin(w http.ResponseWriter, r *http.Request) {
	account := r.FormValue("account")
	password := r.FormValue("password")

	query := fmt.Sprintf(`select * from User where Account="%s" And Password="%s";`, account, password)
	log.Println("[Query]", query)

	rs, err := a.DB.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	var user models.User
	for rs.Next() {
		if err := rs.Scan(&user.UId, &user.Name, &user.Account, &user.Password, &user.CreateTime); err != nil {
			log.Fatal(err)
		}
	}

	if user.UId != 0 {
		log.Printf("[Warning] No Such User!!!")
		c := &http.Cookie{
			Name:  "user",
			Value: account + " " + password,
		}

		http.SetCookie(w, c)
		http.Redirect(w, r, "/User", http.StatusFound)
		return
	}

	log.Printf("[User] UId: %d, Name: %s, Account: %s, Password: %s, CreateTime: %s\n", user.UId, user.Name, user.Account, user.Password, string(user.CreateTime))
	http.Redirect(w, r, "/Login", http.StatusFound)
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
