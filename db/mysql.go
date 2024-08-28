package db

import (
	"database/sql"
	"log"

	"github.com/tubagusmf/payment-service-gb1/internal/helper"

	_ "github.com/go-sql-driver/mysql"
)

func NewMysql() *sql.DB {
	db, err := sql.Open("mysql", helper.GetConnectionString())
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}
