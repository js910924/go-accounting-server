package app

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
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
	// import templates
	allTemplates := []string{
		"index", "login", "register",
		"user", "allData", "income",
		"outlay", "allUsers", "editLog",
		"partials/header", "partials/footer",
	}
	for i := range allTemplates {
		allTemplates[i] = "./template/" + allTemplates[i] + ".html"
	}

	s.Templates = template.Must(template.ParseFiles(allTemplates...))

	s.connectDB("mysql", "root", "0000", "account")

	s.Router = mux.NewRouter()
	s.setRoutes()
}

func (s *Server) setRoutes() {
	s.Router.PathPrefix("/style/css/").Handler(http.StripPrefix("/style/css/", http.FileServer(http.Dir("template/style/css/"))))
	s.Router.HandleFunc("/", s.index()).Methods("GET")
	s.Router.HandleFunc("/Register", s.signUp()).Methods("GET")
	s.Router.HandleFunc("/Login", s.login()).Methods("GET")
	s.Router.HandleFunc("/Login", s.checkLogin()).Methods("POST")
	s.Router.HandleFunc("/Logout", s.logOut()).Methods("GET")
	s.Router.HandleFunc("/PageNotFound", s.pageNotFound()).Methods("GET")

	user := s.Router.PathPrefix("/Users").Subrouter().StrictSlash(true)
	user.HandleFunc("", s.showAllUsers()).Methods("GET")
	user.HandleFunc("", s.createUser()).Methods("POST")
	user.HandleFunc("/{id}", s.showUser()).Methods("GET")
	user.HandleFunc("/{id}/AllData", s.showAllData()).Methods("GET")
	user.HandleFunc("/{id}/Outlay", s.showOutlay()).Methods("GET")
	user.HandleFunc("/{id}/Outlay", s.createAction()).Methods("POST")
	user.HandleFunc("/{id}/Income", s.showIncome()).Methods("GET")
	user.HandleFunc("/{id}/Income", s.createAction()).Methods("POST")
	user.HandleFunc("/{id}/AllData/{actionType}/{logId}", s.editLog()).Methods("GET")
	user.HandleFunc("/{id}/AllData/{actionType}/{logId}", s.editLog()).Methods("POST")
}

func (s *Server) connectDB(driverName string, userName string, password string, dbName string) {
	var err error
	dataSourceName := fmt.Sprintf("%s:%s@tcp(172.20.0.2)/%s", userName, password, dbName)
	s.DB, err = sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatal("[Fail] Open DB fail")
		panic(err)
	}

	if err = s.DB.Ping(); err != nil {
		log.Println("[Fail]", err.Error())

		switch true {
		case strings.Contains(err.Error(), "connection refused"):
			log.Printf("[Handling] Wait mysql server started... On %s\n", runtime.GOOS)
			// log.Printf("[Handling] Trying to start mysql server... On %s\n", runtime.GOOS)
			// var cmd *exec.Cmd
			// switch runtime.GOOS {
			// case "darwin":
			// 	cmd = exec.Command("mysql.server", "start")

			// case "linux":
			// 	cmd = exec.Command("service", "mysql", "start")
			// }

			// if err := cmd.Run(); err != nil {
			// 	log.Println("[Fail] Restart fail, wait 3 sec..., Erro:", err)
			// }

			<-time.After(3 * time.Second) // Wait 3 seconds
			s.connectDB(driverName, userName, password, dbName)
			return

		case strings.Contains(err.Error(), "Access denied"):
			log.Println("[Handling] Please try another password...")
			fmt.Scan(&password)
			log.Println("[New Password]", password)
			s.connectDB(driverName, userName, password, dbName)
			return

		default:
			log.Fatalln("[Ping] Unknown Error")
		}
	}

	log.Println("[Success] MySQL server already start!")
}
