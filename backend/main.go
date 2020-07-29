package main

var err error

func main() {
	var s server
	s.ConfigFile = "config.json"
	s.Config = s.LoadConfig()
	s.database = &DBsql{db: s.OpenDBcxn()}
	err = s.StartServer()
	checkErr(err)

	/*
		db, err = sql.Open("mysql", "root:lusMonkey2412!@/user_auth")
		checkErr(err)
		defer db.Close()
		err = db.Ping()
		checkErr(err)

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
	*/
}
