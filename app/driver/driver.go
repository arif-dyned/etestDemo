package driver

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

//Type DB
type DB struct {
	SQL *sql.DB
}

var (
	//DBConn ...
	dbConn = &DB{}

	dbUsername, _ = os.LookupEnv("DB_USERNAME")
	dbPassword, _ = os.LookupEnv("DB_PASSWORD")
	dbHost, _     = os.LookupEnv("DB_HOST")
	dbName, _     = os.LookupEnv("DB_NAME")
	dbPort, _     = os.LookupEnv("DB_PORT")
)

// Connect to MYSQL
func ConnectToMysql() (*sql.DB, error) {

	dbSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8",
		"root",
		"dynedjayasejahtera",
		"68.183.162.169",
		"3306",
		"etest_api",
	)

	d, err := sql.Open("mysql", dbSource)

	if err != nil {
		log.Fatal(err.Error())
		os.Exit(-1)
	}

	err = d.Ping()
	if err != nil {
		log.Fatalf("Can't establish connection to MySQL server: %s - %v", dbHost, err)
	}

	return d, err
}

// ConnectMgo ....
func ConnectMgo(host, port, uname, pass string) error {

	return nil

}
