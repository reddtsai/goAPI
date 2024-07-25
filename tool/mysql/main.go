package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	gormmysqldriver "gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/reddtsai/goAPI/pkg/blockaction/storage"
)

func main() {
	db, err := sql.Open("mysql", "root:!QAZ2wsx@tcp(127.0.0.1:3306)/blockaction?charset=utf8mb4&parseTime=True")
	if err != nil {
		log.Fatalln(err)
		return
	}
	cfg := gormmysqldriver.Config{
		Conn: db,
	}
	conn, err := gorm.Open(gormmysqldriver.New(cfg), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
		return
	}
	err = conn.AutoMigrate(&storage.UserTable{})
	if err != nil {
		log.Fatalln(err)
		return
	}

}
