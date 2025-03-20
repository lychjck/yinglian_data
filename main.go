package main

import (
	"database/sql"
	"io/fs"
	"log"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
)

func prepareDB()*sql.DB{
	dbName := "yinglian"
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


	db:=prepareDB()
	defer db.Close()
	filepath.WalkDir("./txt", func(path string, d fs.DirEntry, err error) error {
		if err!= nil {
			return err
		}
		if !d.IsDir() {
			parse_yinglian(path,db)
			// fmt.Println(path)
		}
		return nil
	})
	
}
