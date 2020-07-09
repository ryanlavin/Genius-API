package main

import (
	"fmt"
	"github.com/gorilla/mux"
	//"github.com/sirupsen/logrus"
	//"html/template"
	"io/ioutil"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	//"net/url"
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
	if req.Method == "POST" {
		Register(res, req)
		return
	}
	http.ServeFile(res, req, "../frontend/index.html")
}

func Login(res http.ResponseWriter, req *http.Request) {

	username = req.FormValue("username")
	password = req.FormValue("password")

	var dbUser string
	var dbPass string

	err = db.QueryRow("SELECT username, password FROM user WHERE username=?", username).Scan(&dbUser, &dbPass)
	if err == sql.ErrNoRows {
		http.Redirect(res, req, "/", 301)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(password))
	if err != nil {
		http.Redirect(res, req, "/", 301)
		return
	}
	http.Redirect(res, req, "/homepage", 301)
}

func Register(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "../frontend/register.html")
		return
	}

	username = req.FormValue("username")
	password = req.FormValue("password")
	err = db.QueryRow("SELECT username FROM user WHERE username=?", username).Scan(&user)
	if err == sql.ErrNoRows {
		cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			panic(err.Error())
		}
		_, err = db.Exec("INSERT INTO user(username, password) VALUES(?, ?)", username, cryptedPassword)
		if err != nil {
			panic(err.Error())
		}
		http.Redirect(res, req, "/", 301)
		return
	}

	if err != nil {
		panic(err.Error())
	}
	http.Redirect(res, req, "/", 301)

}

func Homepage(res http.ResponseWriter, req *http.Request) {
	if err != nil {
		panic(err.Error())
	}
	http.ServeFile(res, req, "../frontend/homepage.html")
}

func PullLyrics(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "../frontend/pullLyrics.html")
		return
	}
	artist := req.FormValue("artist")
	title := req.FormValue("title")

/*
	if err != nil {
		panic(err.Error())
	}
	if err != nil {
		panic(err.Error())
	}
	*/
	//q := req.URL.Query()
	s := "https://api.lyrics.ovh/v1/"
	s = s + artist + "/"
	s = s + title + "/"
	
	req, err = http.NewRequest("GET", s, nil)
	if err != nil {
		panic(err.Error())
	}

	resp, err := http.Get(s)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	http.ServeFile(res, req, "../frontend/pullLyrics.html")
}

func SavedLyrics(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "../frontend/savedLyrics.html")
}

func main() {
	db, err = sql.Open("mysql", "root:password@/user_auth")
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
	router.HandleFunc("/homepage", Homepage)
	router.HandleFunc("/pulllyrics", PullLyrics)
	router.HandleFunc("/savedlyrics", SavedLyrics)
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
