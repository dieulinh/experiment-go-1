package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	connStr := "postgres://go_app_client:mysecretpassword123@localhost:5432/restapi_go_development?sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to postgres")

}

// package db

// import (
// 	"database/sql"
// 	"log"

// 	_ "github.com/lib/pq"
// )

// var DB *sql.DB

// func Init() {
// 	var err error

// 	DB, err = sql.Open("postgres", "postgres://user:password@localhost:5432/mydb?sslmode=disable")
// 	if err != nil {
// 		log.Fatal("DB connection error:", err)
// 	}

// 	if err = DB.Ping(); err != nil {
// 		log.Fatal("DB not reachable:", err)
// 	}

// 	log.Println("✅ Connected to database")
// }
