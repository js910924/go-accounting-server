# Accounting Server

###### tags: `Golang` `Docker` `Ubuntu18.04`

---

## Build And Run
```shell
$ docker-compose up
$ curl http://localhost:5000
```	

## Todo
- [ ] Modify Cookie format
- [ ] More Detail Search
	> Ex. Show data only in Nov.
- [x] Solved MySQL timestamp to Golang time.Time problem
	> [Link](https://www.jianshu.com/p/3f7fc9093db4)
- [x] Let user select Date (2020.01.08)
- [x] Add delete and edit log feature (2019.12.30)
- [x] Rebuild Project (2019.12.27)
- [x] Use maxAge to delete cookie (2019.12.1)
- [x] Add db connect error handling while os is linux (2109.12.1)
- [x] Rebuild Dockerfile & add deploy project to docker (2019.12.1)
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

### MySQL server can't start
```shell
$service mysql start
 * Starting MySQL database server mysqld
No directory, logging in with HOME=/		[fail]
```

Solution
```shell
$service mysql stop
$usermod -d /var/lib/mysql/ mysql
$ln -s /var/lib/mysql/mysql.sock /tmp/mysql.sock
$chown -R mysql:mysql /var/lib/mysql
$service mysql start
 * Starting MySQL database server mysqld			[ OK ]
```

### Set MySQL root password at first time
```shell
$sudo cat /etc/mysql/debian.cnf
\# Automatically generated for Debian scripts. DO NOT TOUCH!
[client]
host     = localhost
user     = debian-sys-maint
password = ggeu53390yo8tpVY
socket   = /var/run/mysqld/mysqld.sock
[mysql_upgrade]
host     = localhost
user     = debian-sys-maint
password = ggeu53390yo8tpVY
socket   = /var/run/mysqld/mysqld.sock

$mysql -u debian-sys-maint -p
Enter password: # ggeu53390yo8tpVY

$mysql> use mysql;
$mysql> UPDATE user SET plugin="mysql_native_password" WHERE user="root";
$mysql> UPDATE user SET authentication_string=PASSWORD("0000") WHERE user="root";
$mysql> FLUSH PRIVILEGES;
$mysql> exit;
$mysql -u root -p
Enter password: # 0000
```

### Import SQL file
```shell
$mysql -u root -p < db/CreateDB.sql
Enter password: # 0000
$mysql -u root -p < db/InsertTable.sql
Enter password: # 0000
```
