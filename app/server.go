package app

import (
	"database/sql"
	"html/template"
	"log"

	"github.com/gorilla/mux"
)

type Server struct {
	Templates *template.Template
	DB        *sql.DB
	Router    *mux.Router
}

func (s *Server) Init() {
	allTemplates := []string{"index", "login", "register", "user", "allData", "income", "outlay"}
	for i := range allTemplates {
		allTemplates[i] = "./static/public/html/" + allTemplates[i] + ".html"
	}

	s.Templates = template.Must(template.ParseFiles(allTemplates...))
	s.connectDB()

	//---import sql file-----------------------------------
	// file, err := ioutil.ReadFile("./db/CreateDB.sql")
	// if err != nil {
	// 	panic(err)
	// }

	// a.DB.Exec(string(file))
	// file, err = ioutil.ReadFile("./db/InsertTable.sql")
	// if err != nil {
	// 	panic(err)
	// }
	// a.DB.Exec(string(file))
	//---import sql file-----------------------------------

	s.Router = mux.NewRouter()
	s.setRoutes()
}

func (s *Server) setRoutes() {
	s.Router.HandleFunc("/", s.index()).Methods("GET")
	s.Router.HandleFunc("/Register", s.signUp()).Methods("GET")
	s.Router.HandleFunc("/Login", s.login()).Methods("GET")
	s.Router.HandleFunc("/Login", s.checkLogin()).Methods("POST")

	user := s.Router.PathPrefix("/User").Subrouter()
	user.HandleFunc("", s.showUser()).Methods("GET")
	user.HandleFunc("", s.createUser()).Methods("POST")
	user.HandleFunc("/AllData", s.showAllData()).Methods("GET")
	user.HandleFunc("/Outlay", s.showOutlay()).Methods("GET")
	user.HandleFunc("/Outlay", s.createOutlay()).Methods("POST")
	user.HandleFunc("/Income", s.showIncome()).Methods("GET")
}

func (s *Server) connectDB() {
	var err error
	s.DB, err = sql.Open("mysql", "root:0924@/account")
	if err != nil {
		log.Fatal("[Fail] Open DB fail")
	}
}
