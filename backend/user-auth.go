package main

import (
	"context"
	"database/sql"
	//"encoding/json"
	//"fmt"
	_ "net/http"	
	//"encoding/json"
	//"github.com/sirupsen/logrus"
	_ "html/template"
	_ "io/ioutil"
	_ "github.com/go-sql-driver/mysql"
	//"log"
	_ "golang.org/x/crypto/bcrypt"
	//"net/url"
)
type DBsql struct {
	db *sql.DB
	username string
	password string
	ctx *context.Context	
}
//var ctx context.Context
//var db *sql.DB
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

func (s server) OpenDBcxn() *sql.DB {
	userLogin := s.Config.DBUser + ":" + s.Config.DBPass + "@/" + s.Config.AuthTable
	db, err := sql.Open(s.Config.DBHost, userLogin)
	checkErr(err)
	return db
	//s.database.db = db
	//err = s.database.db.Ping()
	//checkErr(err)
	//return db
}

// Error checking 
func checkErr(err error) {
	if err != nil {
		panic(err.Error)
	}
}

// Finds the # of rows in mysql database to ensure correct ordering of songs
func (s server) findOrder() int {
	var count int
	err = s.database.db.QueryRow("SELECT COUNT(*) FROM song").Scan(&count)	
	checkErr(err)
	return count
}

