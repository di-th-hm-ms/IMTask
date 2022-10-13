package service

import (
	"fmt"
	"os"
	"time"
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type DBSettingService struct {}

func (d DBSettingService) OpenConnection(ch chan *sql.DB) *sql.DB {
	driver := os.Getenv("DRIVER")
	// this user is created by shell script when a mysql-container(db-dev)
	// is created and started.
	path := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=true",
		os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("DB_HOST"), os.Getenv("MYSQL_DATABASE"))
	return d.ConnectDB(driver, path, 100, ch)
}
func (d DBSettingService) ConnectDB(driver string, path string, count int, ch chan *sql.DB) *sql.DB {
	db, err := sql.Open(driver, path)
	if err != nil {
		log.Fatal("open error: ", err)
	} else {
		if res, err := db.Query("SHOW DATABASES;"); err != nil{
			log.Fatal("show error: ", err)
		} else {
			var database string
			for res.Next() {
				res.Scan(&database)
			}
		}
		if _, err := db.Query("USE " + os.Getenv("MYSQL_DATABASE") + ";"); err != nil {
			log.Fatal("use error: ", err)
		}
	}

	// Error occurs if ping() is exe before starting a container and launch the server
	// -> retry ping()
	if err = db.Ping(); err != nil {
		// log.Fatal("ping error:", err)
		time.Sleep(time.Second * 2)
		count--
		fmt.Printf("retry count:%v\n", count)
		return d.ConnectDB(driver, path, count, ch)
	}

	// if the connection is going well
	ch <- db
	return db
}
func (DBSettingService) CloseConnection(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Fatal("Failed to close DB connection:", err)
	}
}
