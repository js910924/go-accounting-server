# Accounting Server

###### tags: `Golang` `Docker` `Ubuntu18.04`

---

## Create Dockerfile
```shell
$ sudo touch Dockerfile
$ sudo vim Dockerfile
```	

```Dockerfile
FROM golang:latest
MAINTAINER Jason jswind@myemail.com
LABEL description="This is a accounting server example" version="1.0" owner="Jason Chen"
RUN apt-get update
RUN apt-get install vim net-tools tmux
RUN mkdir $HOME/go
COPY .vimrc ./
```

## Todo
- [ ] Let user can record income
- [ ] Modify Cookie format
- [x] Encrypy Password - Done(2019.11.22)
- [x] Make API more like RESTful APi - Done
	> Ex. /Users/:id
- [x] Show All users data - Done(2019.11.23)
	> Ex. /Users
- [x] Let user can see all his records	- Done(2019.11.24)
	> Ex. /Users/:id/Pool

## Note

### Set Template Value
```go
func UserHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("user")
	if err != nil {
		c = &http.Cookie{
			Name:  "user",
			Value: "",
		}

		fmt.Fprint(w, "No Cookie")
		return
	}

    // Set Template Value
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
```

### Get Form Value, Get Query Values, Set Cookie, Redirect
```go
func checkLogin(w http.ResponseWriter, r *http.Request) {
    // Get Form Value
	account := r.FormValue("account")
	password := r.FormValue("password")

	query := fmt.Sprintf(`select UId from User where Account="%s" And Password="%s";`, account, password)
	fmt.Println(query)

	rs, err := a.DB.Query(query)
	if err != nil {
		log.Fatal(err)
	}

    // Get Query Values 
	var UId int
	for rs.Next() {
		if err := rs.Scan(&UId); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("UId: ", UId)
	if UId != 0 {
        // Set Cookie
		c := &http.Cookie{
			Name:  "user",
			Value: account + " " + password,
		}

		http.SetCookie(w, c)
		http.Redirect(w, r, "/User", http.StatusFound)
		return
	}

    // Redirect
	http.Redirect(w, r, "/Login", http.StatusFound)
}
```

### Scanner

- Scan Word
    ```go
    var input string = ""
    for {
        fmt.Scan(&input)
        fmt.Println(input)
    }
    ```

- Scan Whole Line
    ```go
    var input string = ""
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        input = scanner.Text()
        fmt.Println(input)
    }
    ```
