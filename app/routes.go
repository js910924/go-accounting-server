package app

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"server/middleware"
	"server/models"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) createUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		account := r.FormValue("account")
		password := r.FormValue("password")

		salt1 := "@#$%"
		salt2 := "^&*()"
		h := sha256.New()
		h.Write([]byte(salt1 + name + password + salt2))
		bs := h.Sum(nil)
		password = fmt.Sprintf("%x", bs)

		query := fmt.Sprintf(`INSERT INTO User (Name, Account, Password) VALUES ("%s", "%s", "%s")`, name, account, password)
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

		query := fmt.Sprintf(`SELECT * FROM User WHERE Account="%s" AND Password="%s";`, account, password)
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

		if user.UId != 0 {
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

		query := "SELECT * FROM incomeType ORDER BY type"
		row, err := s.DB.Query(query)
		if err != nil {
			log.Fatal(err)
		}

		allIncome := []models.Income{}
		var income models.Income
		for row.Next() {
			row.Scan(&income.Type, &income.TypeName)
			allIncome = append(allIncome, income)
		}

		middleware.RenderTemplate(s.Templates, w, "income", allIncome)
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

		query := "SELECT * FROM outlayType ORDER BY type"
		row, err := s.DB.Query(query)
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
		middleware.RenderTemplate(s.Templates, w, "outlay", allOutlay)
	}
}

func (s *Server) showAllData() http.HandlerFunc {
	data := "All Data Paga"
	tmplName := "allData"
	return func(w http.ResponseWriter, r *http.Request) {
		userID := mux.Vars(r)["id"]
		c, err := r.Cookie("user")
		if err != nil {
			http.Redirect(w, r, "/Login", http.StatusFound)
			return
		}

		if c.Value == userID {
			middleware.RenderTemplate(s.Templates, w, tmplName, data)
		} else {
			http.Redirect(w, r, "/Users/"+c.Value, http.StatusFound)
		}
	}
}

func (s *Server) createOutlay() http.HandlerFunc {
	// Todo
	return func(w http.ResponseWriter, r *http.Request) {
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
		row, err := s.DB.Query(query)
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
}
