package app

import (
	"fmt"
	"log"
	"net/http"
	"server/middleware"
	"server/models"
	"strconv"
	"strings"
)

func (s *Server) createUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		account := r.FormValue("account")
		password := r.FormValue("password")

		query := fmt.Sprintf(`INSERT INTO User (Name, Account, Password) VALUES ("%s", "%s", "%s")`, name, account, password)
		log.Println("[Query]", query)

		fmt.Println("[Query] I'm Find!!!")
		_, err := s.DB.Query(query)
		if err != nil {
			log.Println(err)
		}

		fmt.Println("[New Cookie] I'm Find!!!")
		c := &http.Cookie{
			Name:  "user",
			Value: name + " " + account + " " + password,
		}

		fmt.Println("[SetCookie] I'm Find!!!")
		http.SetCookie(w, c)
		fmt.Println("[Redirect] I'm Find!!!")
		http.Redirect(w, r, "/User", http.StatusFound)
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
		middleware.RenderTemplate(s.Templates, w, tmplName, data)
	}
}

func (s *Server) showUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		middleware.RenderTemplate(s.Templates, w, "user", item)
	}
}

// func index(w http.ResponseWriter, r *http.Request) {
// 	middleware.RenderTemplate(a.Templates, w, "index", "Home")
// }

// func Register(w http.ResponseWriter, r *http.Request) {
// 	middleware.RenderTemplate(a.Templates, w, "register", "Register Page")
// }

// func login(w http.ResponseWriter, r *http.Request) {
// 	middleware.RenderTemplate(a.Templates, w, "login", "Login")
// }

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
}

func (s *Server) showIncome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		middleware.RenderTemplate(s.Templates, w, tmplName, data)
	}
}

func (s *Server) createOutlay() http.HandlerFunc {
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
