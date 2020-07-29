package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type server struct {
	database   *DBsql
	Config     *Conf
	ConfigFile string
}

func (s server) StartServer() error {
	//s.ConfigFile = "config.json"
	//s.Config = s.LoadConfig()
	//fmt.Println(s.Config.IndexURL, s.Config.LoginURL, s.Config.RegisterURL, s.Config.HomepageURL, s.Config.PullLyricsURL)
	router := mux.NewRouter()
	router.HandleFunc(s.Config.IndexURL, s.Orig)
	router.HandleFunc(s.Config.LoginURL, s.Login)
	router.HandleFunc(s.Config.RegisterURL, s.Register)
	router.HandleFunc(s.Config.HomepageURL, s.Homepage)
	router.HandleFunc(s.Config.PullLyricsURL, s.PullLyrics)
	err = http.ListenAndServe(":8080", router)
	checkErr(err)
	return nil
}

// Original html file dealt
func (s server) Orig(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		s.Register(res, req)
		return
	}
	http.ServeFile(res, req, s.Config.IndexFile)
}

// Login Form
func (s server) Login(res http.ResponseWriter, req *http.Request) {

	username = req.FormValue("email")
	password = req.FormValue("password")

	var dbUser string
	var dbPass string
	err = s.database.db.QueryRow("SELECT username, password FROM user WHERE username=?", username).Scan(&dbUser, &dbPass)
	//err = db.QueryRow("SELECT username, password FROM user WHERE username=?", username).Scan(&dbUser, &dbPass)
	if err == sql.ErrNoRows {
		http.Redirect(res, req, s.Config.IndexURL, 301)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(password))
	if err != nil {
		http.Redirect(res, req, s.Config.IndexURL, 301)
		return
	}
	http.Redirect(res, req, s.Config.HomepageURL, 301)
}

// Form dealt to register account in mysql db, submitting redirects to index form
// to login
func (s server) Register(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, s.Config.RegisterFile)
		return
	}

	username = req.FormValue("email")
	password = req.FormValue("password")
	err = s.database.db.QueryRow("SELECT username FROM user WHERE username=?", username).Scan(&user)
	if err == sql.ErrNoRows {
		cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		checkErr(err)
		_, err = s.database.db.Exec("INSERT INTO user(username, password) VALUES(?, ?)", username, cryptedPassword)
		checkErr(err)
		http.Redirect(res, req, s.Config.IndexURL, 301)
		return
	}
	checkErr(err)
	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	checkErr(err)
	_, err = s.database.db.Exec("INSERT INTO user(username, password) VALUES(?, ?)", username, cryptedPassword)
	checkErr(err)
	http.Redirect(res, req, s.Config.IndexURL, 301)
	return
}

// for writing to json file to be read from by javascript to write to frontend html
type toFrontend struct {
	Row1 string `json: "Row1"`
	Row2 string `json: "Row2"`
	Row3 string `json: "Row3"`
}

// First page dealt after successfully signing in
func (s server) Homepage(res http.ResponseWriter, req *http.Request) {
	var row1, row2, row3 string
	err = s.database.db.QueryRow("SELECT songLink from song WHERE songOrder=0").Scan(&row1)
	checkErr(err)
	err = s.database.db.QueryRow("SELECT songLink from song WHERE songOrder=1").Scan(&row2)
	checkErr(err)
	err = s.database.db.QueryRow("SELECT songLink from song WHERE songOrder=2").Scan(&row3)
	checkErr(err)
	fmt.Println("rows:", row1, row2, row3)
	dataStr := toFrontend{
		Row1: row1,
		Row2: row2,
		Row3: row3,
	}
	data, _ := json.MarshalIndent(dataStr, "", " ")
	checkErr(err)
	_ = ioutil.WriteFile("artist.json", data, os.FileMode(0666))
	//FirstSong := req.FormValue("Song1")
	//SecondSong := req.FormValue("Song2")
	//ThirdSong := req.FormValue("Song3")
	//fmt.Println("current web names:", FirstSong, SecondSong, ThirdSong)
	//Field := "Field"
	//tmpl, err := template.ParseFiles(s.Config.HomepageFile)
	//checkErr(err)
	//f := new(Field)
	//f := new([3]string{row1, row2, row3})
	//f.SongOne = row1
	//f.SongTwo = row2
	//f.SongThree = row3
	//err = tmpl.Execute(res, &f)
	//checkErr(err)
	//array := [3]string{row1, row2, row3}
	//jsonData, err := json.Marshal(array)
	//checkErr(err)
	//json.NewEncoder(res).Encode(jsonData)
	http.ServeFile(res, req, s.Config.HomepageFile)
}

// Page for pulling lyrics from lyrics.ovh API
// Calls getLyrics() function to stream data into db & download lyrics as .txt file
func (s server) PullLyrics(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		fmt.Println("Not POST")
		http.ServeFile(res, req, s.Config.PullLyricsFile)
		return
	}
	artist := req.FormValue("artist")
	title := req.FormValue("title")
	url := "https://api.lyrics.ovh/v1/"
	url = url + artist + "/"
	url = url + title + "/"
	fileName := artist + " - " + title + ".txt"
	h := dbHandler{artist, title, url, fileName}
	req = s.getLyrics(h)
	http.ServeFile(res, req, s.Config.PullLyricsFile)
}

// Implements song interface, this function executes upon a POST request being handled by PullLyrics
// Checks if user input of the artist and song name is in the database and proceeds to pull the lyrics from the API and stream them both into the database and into a txt file for easy retrieval
func (s server) getLyrics(h dbHandler) (req *http.Request) {
	req, err := http.NewRequest("GET", h.url, nil)
	checkErr(err)

	resp, err := http.Get(h.url)
	checkErr(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	songLink := strings.TrimRight(h.file, ".txt")
	fmt.Println("songLink:", songLink)
	err = s.database.db.QueryRow("SELECT fileName FROM song WHERE fileName=?", h.file).Scan(&h.file)
	if err == sql.ErrNoRows {
		num := s.findOrder()
		_, err = s.database.db.Exec("INSERT INTO song(fileName, songLink, songOrder) VALUES(?, ?, ?)", h.file, songLink, num)
		checkErr(err)

		file, err := os.Create(h.file)
		checkErr(err)
		defer file.Close()

		fmt.Fprintf(file, string(body), h.file)
		return req
	}
	checkErr(err)
	num := s.findOrder()
	_, err = s.database.db.Exec("INSERT INTO song(fileName, songLink, songOrder) VALUES(?, ?, ?)", h.file, h.url, num)
	checkErr(err)
	return req
}

// Work in progress
func (s server) SavedLyrics(res http.ResponseWriter, req *http.Request) {
	//http.ServeFile(res, req, "../frontend/savedLyrics.html")

	// put stirngs into a list of json objects and then return the list
}
