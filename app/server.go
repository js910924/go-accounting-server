package app

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	Templates *template.Template
	DB        *sql.DB
	Router    *mux.Router
}

func (s *Server) Init() {
	log.SetFlags(log.Lshortfile)
	allTemplates := []string{"index", "login", "register", "user", "allData", "income", "outlay"}
	for i := range allTemplates {
		allTemplates[i] = "./static/public/html/" + allTemplates[i] + ".html"
	}

	s.Templates = template.Must(template.ParseFiles(allTemplates...))
	s.connectDB("mysql", "root", "0924", "account")

	//---import sql file-----------------------------------
	// file, err := ioutil.ReadFile("./db/CreateDB.sql")
	// if err != nil {
	// 	panic(err)
	// }

	// rs, err := s.DB.Exec(string(file))
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(rs)
	// file, err = ioutil.ReadFile("./db/InsertTable.sql")
	// if err != nil {
	// 	panic(err)
	// }
	// s.DB.Exec(string(file))
	//---import sql file-----------------------------------

	s.Router = mux.NewRouter()
	s.setRoutes()
}

func (s *Server) setRoutes() {
	s.Router.HandleFunc("/", s.index()).Methods("GET")
	s.Router.HandleFunc("/Register", s.signUp()).Methods("GET")
	s.Router.HandleFunc("/Login", s.login()).Methods("GET")
	s.Router.HandleFunc("/Login", s.checkLogin()).Methods("POST")

	user := s.Router.PathPrefix("/Users").Subrouter()
	// user.HandleFunc("", s.showAllUser()).Methods("GET")
	user.HandleFunc("", s.createUser()).Methods("POST")
	user.HandleFunc("/{id}", s.showUser()).Methods("GET")
	user.HandleFunc("/{id}/AllData", s.showAllData()).Methods("GET")
	user.HandleFunc("/{id}/Outlay", s.showOutlay()).Methods("GET")
	user.HandleFunc("/{id}/Outlay", s.createOutlay()).Methods("POST")
	user.HandleFunc("/{id}/Income", s.showIncome()).Methods("GET")
	// user.HandleFunc("/{id}/Income", s.showIncome()).Methods("POST")
}

func (s *Server) connectDB(driverName string, userName string, password string, dbName string) {
	var err error
	dataSourceName := fmt.Sprintf("%s:%s@/%s", userName, password, dbName)
	s.DB, err = sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatal("[Fail] Open DB fail")
		panic(err)
	}

	if err = s.DB.Ping(); err != nil {
		log.Println("[Fail]", err.Error())

		if runtime.GOOS != "darwin" {
			log.Fatal("[Fail] Connect DB ERROR... Unless you use MacOS")
			return
		}

		switch true {
		case strings.Contains(err.Error(), "connection refused"):
			log.Println("[Handling] Trying start mysql server...")
			cmd := exec.Command("mysql.server", "start")
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}

			<-time.After(3 * time.Second) // Wait 3 seconds
			s.connectDB(driverName, userName, password, dbName)

		case strings.Contains(err.Error(), "Access denied"):
			log.Println("[Handling] Please try another password...")
			fmt.Scan(&password)
			fmt.Println("[New Password]", password)
			s.connectDB(driverName, userName, password, dbName)

		default:
			log.Fatalln("[Ping] Unknown Error")
		}
	}

	log.Println("[Success] MySQL server already start!")
}
