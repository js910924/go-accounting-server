package app

import (
	"fmt"
	"log"
	"net/http"
	"server/middleware"
	"server/models"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (s *Server) showAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := fmt.Sprintf(`SELECT * FROM User;`)
		log.Println("[Query]", query)

		row, err := s.DB.Query(query)
		if err != nil {
			log.Println(err)
		}

		var user models.User
		var allUser []models.User = []models.User{}
		for row.Next() {
			row.Scan(&user.UId, &user.Name, &user.Account, &user.Password, &user.CreateTime)
			allUser = append(allUser, user)
		}

		middleware.RenderTemplate(s.Templates, w, "allUsers", allUser)
	}
}

func (s *Server) createUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		account := r.FormValue("account")
		password := r.FormValue("password")

		//---sha256 + salt------------------------
		// salt1 := "@#$%"
		// salt2 := "^&*()"
		// h := sha256.New()
		// h.Write([]byte(salt1 + name + password + salt2))
		// bs := h.Sum(nil)
		// password = fmt.Sprintf("%x", bs)
		//---sha256 + salt------------------------
		hashPwd := middleware.HashAndSalt([]byte(password))

		query := fmt.Sprintf(`INSERT INTO User (Name, Account, Password) VALUES ("%s", "%s", "%s")`, name, account, hashPwd)
		log.Println("[Query]", query)

		_, err := s.DB.Query(query)
		if err != nil {
			log.Println(err)
		}

		http.Redirect(w, r, "/Login", http.StatusFound)
	}
}

func (s *Server) index() http.HandlerFunc {
	data := "Home"
	tmplName := "index"
	return func(w http.ResponseWriter, r *http.Request) {
		middleware.RenderTemplate(s.Templates, w, tmplName, data)
	}
}

func (s *Server) signUp() http.HandlerFunc {
	data := "Sign Up Page"
	tmplName := "register"
	return func(w http.ResponseWriter, r *http.Request) {
		middleware.RenderTemplate(s.Templates, w, tmplName, data)
	}
}

func (s *Server) login() http.HandlerFunc {
	data := "Login"
	tmplName := "login"
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("user")
		if err != nil {
			middleware.RenderTemplate(s.Templates, w, tmplName, data)
			return
		}

		http.Redirect(w, r, "/Users/"+c.Value, http.StatusFound)
	}
}

func (s *Server) logOut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("user")
		if err != nil {
			log.Fatal(err)
		}

		// Clear Cookie
		if c.String() != "" {
			c.Value = ""
			c.Expires = time.Unix(0, 0)
			c.HttpOnly = true
			http.SetCookie(w, c)
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (s *Server) showUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Haven't login yet
		c, err := r.Cookie("user")
		if err != nil {
			http.Redirect(w, r, "/Login", http.StatusFound)
			return
		}

		// Not that user's account
		userID := mux.Vars(r)["id"]
		if c.Value != userID {
			http.Redirect(w, r, "/Users/"+c.Value, http.StatusFound)
			return
		}

		// Search database
		query := fmt.Sprintf("SELECT * FROM User WHERE UId=%s;", userID)
		row, err := s.DB.Query(query)
		if err != nil {
			log.Fatal(err)
		}

		user := models.User{}
		for row.Next() {
			if err := row.Scan(&user.UId, &user.Name, &user.Account, &user.Password, &user.CreateTime); err != nil {
				log.Fatal(err)
			}
		}

		middleware.RenderTemplate(s.Templates, w, "user", user)
	}
}

func (s *Server) checkLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		account := r.FormValue("account")
		password := r.FormValue("password")

		// query := fmt.Sprintf(`SELECT * FROM User WHERE Account="%s" AND Password="%s";`, account, hashPwd)
		query := fmt.Sprintf(`SELECT * FROM User WHERE Account="%s";`, account)
		log.Println("[Query]", query)

		row, err := s.DB.Query(query)
		if err != nil {
			log.Fatal(err)
		}

		var user models.User
		for row.Next() {
			if err := row.Scan(&user.UId, &user.Name, &user.Account, &user.Password, &user.CreateTime); err != nil {
				log.Fatal(err)
			}
		}

		if middleware.ComparePasswords(user.Password, []byte(password)) {
			log.Println("[Success] Login")
			id := strconv.Itoa(user.UId)
			c := &http.Cookie{
				Name:  "user",
				Value: id,
			}

			http.SetCookie(w, c)
			http.Redirect(w, r, "/Users/"+id, http.StatusFound)
			return
		}

		log.Printf("[Warning] No Such User!!!")
		log.Printf("[User] Account: %s, Password: %s\n", account, password)
		http.Redirect(w, r, "/Login", http.StatusFound)
	}
}

func (s *Server) showIncome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := mux.Vars(r)["id"]
		c, err := r.Cookie("user")
		if err != nil {
			http.Redirect(w, r, "/Login", http.StatusFound)
			return
		}

		if c.Value != userID {
			http.Redirect(w, r, "/Users/"+c.Value, http.StatusFound)
			return
		}

		query := "SELECT * FROM Action WHERE ActionType=2 ORDER BY DetailType"
		row, err := s.DB.Query(query)
		if err != nil {
			log.Fatal(err)
		}

		allIncome := []models.Action{}
		var action models.Action
		for row.Next() {
			row.Scan(&action.ActionType, &action.DetailType, &action.DetailName)
			allIncome = append(allIncome, action)
		}

		packet := models.Packet{
			UserID: userID,
			Data:   allIncome,
		}

		middleware.RenderTemplate(s.Templates, w, "income", packet)
	}
}

func (s *Server) showOutlay() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := mux.Vars(r)["id"]
		c, err := r.Cookie("user")
		if err != nil {
			http.Redirect(w, r, "/Login", http.StatusFound)
			return
		}

		if c.Value != userID {
			http.Redirect(w, r, "/Users/"+c.Value, http.StatusFound)
			return
		}

		query := "SELECT * FROM Action WHERE ActionType=1 ORDER BY DetailType"
		row, err := s.DB.Query(query)
		if err != nil {
			log.Fatal(err)
		}

		allOutlay := []models.Action{}
		var action models.Action
		for row.Next() {
			row.Scan(&action.ActionType, &action.DetailType, &action.DetailName)
			allOutlay = append(allOutlay, action)
		}

		packet := models.Packet{
			UserID: userID,
			Data:   allOutlay,
		}

		log.Println("[Success] All Outlays:", allOutlay)
		middleware.RenderTemplate(s.Templates, w, "outlay", packet)
	}
}

func (s *Server) showAllData() http.HandlerFunc {
	tmplName := "allData"
	return func(w http.ResponseWriter, r *http.Request) {
		userID := mux.Vars(r)["id"]
		c, err := r.Cookie("user")
		if err != nil {
			http.Redirect(w, r, "/Login", http.StatusFound)
			return
		}

		if c.Value == userID {
			query := fmt.Sprintf("SELECT * FROM Log WHERE UserId=%s", userID)
			row, err := s.DB.Query(query)
			if err != nil {
				log.Fatal(err)
			}

			var data models.Data
			var allData []models.Data = []models.Data{}
			for row.Next() {
				row.Scan(&data.UserId, &data.ActionType, &data.DetailType, &data.Money, &data.Description, &data.CreateTime)
				allData = append(allData, data)
			}

			packet := models.Packet{
				UserID: userID,
				Data:   allData,
			}

			middleware.RenderTemplate(s.Templates, w, tmplName, packet)
		} else {
			http.Redirect(w, r, "/Users/"+c.Value+"/AllData", http.StatusFound)
		}
	}
}

func (s *Server) createAction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := mux.Vars(r)["id"]
		actionType, err := strconv.Atoi(r.FormValue("actionType"))
		if err != nil {
			log.Fatal(err)
		}
		id, err := strconv.Atoi(userID)
		if err != nil {
			log.Fatal(err)
		}

		detailType, err := strconv.Atoi(r.FormValue("detailType"))
		if err != nil {
			log.Fatalln("[Fail]", err)
		}

		money, err := strconv.Atoi(r.FormValue("money"))
		if err != nil {
			log.Fatalln("[Fail]", err)
		}

		var data models.Data = models.Data{
			UserId:      id,
			ActionType:  actionType,
			DetailType:  detailType,
			Money:       money,
			Description: r.FormValue("description"),
		}

		query := fmt.Sprintf("INSERT INTO Log (UserId, ActionType, DetailType, Money, Description) values ('%d', '%d', '%d', '%d', '%s');", data.UserId, data.ActionType, data.DetailType, data.Money, data.Description)
		log.Println("[Query]", query)
		row, err := s.DB.Query(query)
		if err != nil {
			log.Fatalln("[Fail]", err)
		}

		var result string
		for row.Next() {
			row.Scan(&result)
		}

		log.Println("[Success]", result, "Data:", data)
		http.Redirect(w, r, "/Users/"+userID, http.StatusFound)
	}
}

type Log struct {
	ActionName  string
	DetailName  string
	Money       int
	Description string
	CreateTime  []int8
}
