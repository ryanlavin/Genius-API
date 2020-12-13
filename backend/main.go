package main

var err error

func main() {
	var s server
	s.ConfigFile = "config.json"
	s.Config = s.LoadConfig()
	s.database = &DBsql{db: s.OpenDBcxn()}
	err = s.StartServer()
	checkErr(err)
}
