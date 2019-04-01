package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os/user"

	"github.com/FogCreek/mini"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func params() string {
	u, err := user.Current()
	fatal(err)
	cfg, err := mini.LoadConfiguration(u.HomeDir + "/.phonebookrc")
	fatal(err)

	info := fmt.Sprintf("host=%s port=%s dbname=%s "+
		"sslmode=%s user=%s password=%s ",
		cfg.String("host", "127.0.0.1"),
		cfg.String("port", "5432"),
		cfg.String("dbname", u.Username),
		cfg.String("sslmode", "disable"),
		cfg.String("user", u.Username),
		cfg.String("pass", ""),
	)
	return info
}

var db *sql.DB

func main() {

	var err error
	db, err = sql.Open("postgres", params())
	fatal(err)
	defer db.Close()
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS " +
		`phonebook("id" SERIAL PRIMARY KEY,` +
		`"ProjectId" int, "IssueId" varchar(20), "PongsCounter" varchar(20))`)
	fatal(err)

	router := httprouter.New()
	router.GET("/api/v1/records", getRecords)
	router.GET("/api/v1/records/:id", getRecord)
	router.POST("/api/v1/records", addRecord)
	router.PUT("/api/v1/records/:id", updateRecord)
	router.DELETE("/api/v1/records/:id", deleteRecord)
	http.ListenAndServe(":8080", router)
}

//	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s " + "password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
//	sqlStatement := `INSERT INTO users (age, email, firstname, lastname)
//VALUES ($1, $2, $3, $4)`
//	db, err := sql.Open("postgres", psqlInfo)
//	if err != nil {
//		panic(err)
//	}
//	defer db.Close()
//
//	_, err = db.Exec(sqlStatement, 12, "tester 23133wow", "i like to parse bithch", "kartofelniy")
//	if err != nil {
//		panic(err)
//	}
//
//	err = db.Ping()
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println("Zaebok")
