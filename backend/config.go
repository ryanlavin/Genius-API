package main

import (
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"os"
)

type Conf struct {
	DBHost         string `json:"DBHost"`
	DBName         string `json:"DBName"`
	DBUser         string `json:"DBUser"`
	DBPass         string `json:"DBPass"`
	AuthTable      string `json:"AuthTable"`
	SongTable      string `json:"SongTable"`
	IndexURL       string `json:"IndexURL"`
	IndexFile      string `json:"IndexFile"`
	RegisterURL    string `json:"RegisterURL"`
	RegisterFile   string `json:"RegisterFile"`
	LoginURL       string `json:"LoginURL"`
	HomepageURL    string `json:"HomepageURL"`
	HomepageFile   string `json:"HomepageFile"`
	PullLyricsURL  string `json:"PullLyricsURL"`
	PullLyricsFile string `json:"PullLyricsFile"`
	JsonURL string `json:"jsonURL"`
	JsonFile string `json:"jsonFile"`
}

func (s server) LoadConfig() *Conf {
	jsonFile, err := os.Open(s.ConfigFile)
	checkErr(err)
	defer jsonFile.Close()
	confData, err := ioutil.ReadAll(jsonFile)
	checkErr(err)
	json.Unmarshal(confData, &s.Config)
	return s.Config
	//s.Config = getConfig()
}

func getConfig() *Conf {
	config := &Conf{
		DBHost:         "mysql",
		DBName:         "user_auth",
		DBUser:         "root",
		DBPass:         "lusMonkey2412!",
		AuthTable:      "user_auth",
		SongTable:      "song",
		IndexURL:       "/",
		IndexFile:      "../frontend/index.html",
		RegisterURL:    "/register",
		RegisterFile:   "../frontend/register.html",
		LoginURL:       "/login",
		HomepageURL:    "/homepage",
		HomepageFile:   "../frontend/homepage.html",
		PullLyricsURL:  "/pulllyrics",
		PullLyricsFile: "../frontend/pullLyrics.html",
	}
	return config
}
