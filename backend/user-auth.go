package main

import (
	"fmt"
	"github.com/gorilla/mux"
	//"github.com/sirupsen/logrus"
	//"html/template"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var db *sql.DB
var err error
var username string
var password string
var user string

type Handler struct {
	url    string
	method string
}

func Orig(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "../frontend/index.html")
}

func Login(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "../frontend/successForm.html")

	username = req.FormValue("username")
	password = req.FormValue("password")

	var dbUser string
	var dbPass string

	err = db.QueryRow("SELECT username, password FROM user WHERE username=?", username).Scan(&dbUser, &dbPass)
	if err != nil {
		fmt.Println("username at least")
		http.Redirect(res, req, "/login", 301)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(password))
	if err != nil {
		fmt.Println("password too")
		http.Redirect(res, req, "/login", 301)
		return
	}
}

func Register(res http.ResponseWriter, req *http.Request) {
	//http.ServeFile(res, req, "../frontend/register.html")
	if req.Method != "POST" {
		fmt.Println("HM")
		http.ServeFile(res, req, "../frontend/register.html")
		return
	}

	fmt.Println("HM")
	username = req.FormValue("username")
	password = req.FormValue("password")
	err = db.QueryRow("SELECT username FROM user WHERE username=?", username).Scan(&user)
	if err == sql.ErrNoRows {
		cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			panic(err.Error())
		}
		_, err = db.Exec("INSERT INTO user(username, password) VALUES(?, ?)", username, cryptedPassword)
		//defer req.ParseForm()
		//defer fmt.Println(req.Form["username"])
		//defer fmt.Println(req.Form["password"])
		if err != nil {
			panic(err.Error())
		}
		return
	}
	if err != nil {
		panic(err.Error())
	}
	http.Redirect(res, req, "/", 301)

}

func main() {
	db, err = sql.Open("mysql", "root:lusMonkey2412!@/user_auth")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	router := mux.NewRouter()
	router.HandleFunc("/", Orig)
	router.HandleFunc("/login", Login)
	router.HandleFunc("/register", Register)
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}

}
