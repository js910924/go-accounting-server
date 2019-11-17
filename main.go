package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"server/middleware"
	"server/models"
	"strconv"
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
	allTemplates := []string{"index", "login", "register", "user", "allData", "income", "outlay"}
	for i := range allTemplates {
		allTemplates[i] = "./static/public/" + allTemplates[i] + ".html"
	}

	a.Templates = template.Must(template.ParseFiles(allTemplates...))
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
	user := a.Router.PathPrefix("/User").Subrouter()
	user.HandleFunc("", ShowUser).Methods("GET")
	user.HandleFunc("", CreateUser).Methods("POST")
	user.HandleFunc("/AllData", ShowAllData).Methods("GET")
	user.HandleFunc("/Outlay", ShowOutlay).Methods("GET")
	user.HandleFunc("/Outlay", CreateOutlay).Methods("POST")
	user.HandleFunc("/Income", ShowIncome).Methods("GET")

	a.Router.HandleFunc("/Register", Register).Methods("GET")
	a.Router.HandleFunc("/Login", login).Methods("GET")
	a.Router.HandleFunc("/Login", checkLogin).Methods("POST")

	log.Println("Starting Server...")
	log.Fatal(http.ListenAndServe(":80", a.Router))
}

func ShowUser(w http.ResponseWriter, r *http.Request) {
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
		Name     string
		Account  string
		Password string
	}{
		Name:     values[0],
		Account:  values[1],
		Password: values[2],
	}

	middleware.RenderTemplate(a.Templates, w, "user", item)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	account := r.FormValue("account")
	password := r.FormValue("password")

	query := fmt.Sprintf(`INSERT INTO User (Name, Account, Password) VALUES ("%s", "%s", "%s")`, name, account, password)
	log.Println("[Query]", query)

	_, err := a.DB.Query(query)
	if err != nil {
		log.Println(err)
	}

	c := &http.Cookie{
		Name:  "user",
		Value: name + " " + account + " " + password,
	}

	http.SetCookie(w, c)
	http.Redirect(w, r, "/User", http.StatusFound)
}

func index(w http.ResponseWriter, r *http.Request) {
	middleware.RenderTemplate(a.Templates, w, "index", "Home")
}

func Register(w http.ResponseWriter, r *http.Request) {
	middleware.RenderTemplate(a.Templates, w, "register", "Register Page")
}

func login(w http.ResponseWriter, r *http.Request) {
	middleware.RenderTemplate(a.Templates, w, "login", "Login")
}

func checkLogin(w http.ResponseWriter, r *http.Request) {
	account := r.FormValue("account")
	password := r.FormValue("password")

	query := fmt.Sprintf(`SELECT * FROM User WHERE Account="%s" AND Password="%s";`, account, password)
	log.Println("[Query]", query)

	row, err := a.DB.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	var user models.User
	for row.Next() {
		if err := row.Scan(&user.UId, &user.Name, &user.Account, &user.Password, &user.CreateTime); err != nil {
			log.Fatal(err)
		}
	}

	if user.UId != 0 {
		log.Println("[Success] Login")
		c := &http.Cookie{
			Name:  "user",
			Value: user.Name + " " + user.Account + " " + user.Password,
		}

		http.SetCookie(w, c)
		http.Redirect(w, r, "/User", http.StatusFound)
		return
	}

	log.Printf("[Warning] No Such User!!!")
	log.Printf("[User] UId: %d, Name: %s, Account: %s, Password: %s, CreateTime: %s\n", user.UId, user.Name, user.Account, user.Password, string(user.CreateTime))
	http.Redirect(w, r, "/Login", http.StatusFound)
}

func ShowIncome(w http.ResponseWriter, r *http.Request) {
	query := "SELECT * FROM incomeType ORDER BY type"
	row, err := a.DB.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	allIncome := []models.Income{}
	var income models.Income
	for row.Next() {
		row.Scan(&income.Type, &income.TypeName)
		allIncome = append(allIncome, income)
	}

	middleware.RenderTemplate(a.Templates, w, "income", allIncome)
}

func ShowOutlay(w http.ResponseWriter, r *http.Request) {
	query := "SELECT * FROM outlayType ORDER BY type"
	row, err := a.DB.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	allOutlay := []models.Outlay{}
	var outlay models.Outlay
	for row.Next() {
		row.Scan(&outlay.Type, &outlay.TypeName)
		allOutlay = append(allOutlay, outlay)
	}
	log.Println("[Success] All Outlays:", allOutlay)
	middleware.RenderTemplate(a.Templates, w, "outlay", allOutlay)
}

func ShowAllData(w http.ResponseWriter, r *http.Request) {
	middleware.RenderTemplate(a.Templates, w, "allData", "All Data Page")
}

func CreateOutlay(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Form)
	detailType, err := strconv.Atoi(r.FormValue("detailType"))
	if err != nil {
		log.Fatalln("[Fail]", err)
	}

	money, err := strconv.Atoi(r.FormValue("money"))
	if err != nil {
		log.Fatalln("[Fail]", err)
	}

	var data models.Data = models.Data{
		UserId:      1,
		ActionType:  1,
		DetailType:  detailType,
		Money:       money,
		Description: r.FormValue("description"),
	}

	query := fmt.Sprintf("INSERT INTO pool (UserId, ActionType, DetailType, Money, Description) values ('%d', '%d', '%d', '%d', '%s');", data.UserId, data.ActionType, data.DetailType, data.Money, data.Description)
	log.Println("[Query]", query)
	row, err := a.DB.Query(query)
	if err != nil {
		log.Fatalln("[Fail]", err)
	}

	var result string
	for row.Next() {
		row.Scan(&result)
	}

	log.Println("[Success]", result, "Data:", data)
	http.Redirect(w, r, "/User", http.StatusFound)
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
