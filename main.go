package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"flag"
)

func prepareDB()*sql.DB{
	dbName := "couples"
	dbHost := "localhost"
	dbUser := "root"
	dbPass := ""
	log.Printf("[Database] Connecting to MySQL database %s on %s", dbName, dbHost)
	db, err := sql.Open("mysql", dbUser+":"+dbPass+"@tcp("+dbHost+")/"+dbName+"?parseTime=true&loc=Local&charset=utf8mb4")
	if err != nil {
		log.Fatal("[Error] Database connection failed: ", err)
	}
	log.Println("[Database] Connected successfully")
	return db
}

func main() {
	var filename string
	flag.StringVar(&filename,"name","","name")
	flag.Parse()
	db:=prepareDB()
	defer db.Close()
	parse_yinglian(filename,db)
}
