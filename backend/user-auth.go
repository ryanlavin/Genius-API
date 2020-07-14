package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	//"github.com/sirupsen/logrus"
	//"html/template"
	"io/ioutil"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	//"net/url"
)
var ctx context.Context
var db *sql.DB
var err error
var username string
var password string
var user string
var songName string

type dbHandler struct {
	artist string
	title string
	url    string
	file string
}

type song interface {
	getLyrics()
	//reorderLyrics()
}

// Implements song interface, this function executes upon a POST request being handled by PullLyrics
// Checks if user input of the artist and song name is in the database and proceeds to pull the lyrics from the API and stream them both into the database and into a txt file for easy retrieval
func (h dbHandler) getLyrics() (req1 *http.Request) {
	req, err := http.NewRequest("GET", h.title, nil)
	checkErr(err)	

	resp, err := http.Get(h.url)
	checkErr(err)	
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	
	err = db.QueryRow("SELECT fileName FROM song WHERE fileName=?", h.file).Scan(&h.file)	
	if err == sql.ErrNoRows {
		num := findOrder()
		_, err = db.Exec("INSERT INTO song(fileName, songLink, songOrder) VALUES(?, ?, ?)", h.file, h.url, num)
		checkErr(err)	
	
		file, err := os.Create(h.file)
		checkErr(err)	
		defer file.Close()
		
		fmt.Fprintf(file, string(body), h.file)
		return req	
	}
	checkErr(err)
	num := findOrder()
	_, err = db.Exec("INSERT INTO song(fileName, songLink, songOrder) VALUES(?, ?, ?)", h.file, h.url, num)
	checkErr(err)
	return req	
}

// Error checking 
func checkErr(err error) {
	if err != nil {
		panic(err.Error)
	}
}

// Finds the # of rows in mysql database to ensure correct ordering of songs
func findOrder() int {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM song").Scan(&count)	
	checkErr(err)
	return count
}

// Original html file dealt
func Orig(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		Register(res, req)
		return
	}
	http.ServeFile(res, req, "../frontend/index.html")
}

// Login Form
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

// Form dealt to register account in mysql db, submitting redirects to index form 
// to login
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
		checkErr(err)	
		_, err = db.Exec("INSERT INTO user(username, password) VALUES(?, ?)", username, cryptedPassword)
		checkErr(err)	
		http.Redirect(res, req, "/", 301)
		return
	}
	checkErr(err)
	http.Redirect(res, req, "/", 301)

}

// First page dealt after successfully signing in
func Homepage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "../frontend/homepage.html")
}

// Page for pulling lyrics from lyrics.ovh API
// Calls getLyrics() function to stream data into db & download lyrics as .txt file
func PullLyrics(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		fmt.Println("Not POST")
		http.ServeFile(res, req, "../frontend/pullLyrics.html")
		return
	}
	artist := req.FormValue("artist")
	title := req.FormValue("title")
	link := "https://api.lyrics.ovh/v1/"
	link = link + artist + "/"
	link = link + title + "/"

	fileName := artist + " - " + title + ".txt"
	h := dbHandler{artist, title, link, fileName}	
	req = h.getLyrics()	
	http.ServeFile(res, req, "../frontend/pullLyrics.html")
}

// Work in progress
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
