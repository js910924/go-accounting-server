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
RUN apt-get install vim -y
RUN apt-get install net-tools -y
RUN apt-get install tmux -y
COPY ./ ./src/go-account
```

## Todo
- [ ] Modify Cookie format
- [ ] More Detail Search
	> Ex. Show data only in Nov.

- [x] Add logout feature => logout and delete cookie (2019.11.30)
- [x] Add static file serve	(2019.11.27)
	> Ex. .css & .js file
- [x] Add html header & footer (2019.11.27)
- [x] Add Navbar (2019.11.27)
- [x] Let user can record income (2019.11.26)
- [x] Encrypy Password (2019.11.22)
- [x] Make API more like RESTful APi (2019.11.24)
	> Ex. /Users/:id
- [x] Let user can see all his records (2019.11.24)
	> Ex. /Users/:id/Logs
- [x] Show All users data (2019.11.23)
	> Ex. /Users

## Bug
- [ ] Only main router can find css path, subrouter should find too.

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
