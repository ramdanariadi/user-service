package setup

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/ramdanariadi/grocery-user-service/utils"
	"os"
	"strings"
	"time"
)

func NewDbConnection() (*sql.DB, error) {

	dbUsr := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	args := os.Args
	for _, arg := range args {
		split := strings.Split(arg, "=")
		switch split[0] {
		case "DB_USERNAME":
			dbUsr = split[1]
			break
		case "DB_PASS":
			dbPass = split[1]
			break
		case "DB_NAME":
			dbName = split[1]
			break
		case "DB_HOST":
			dbHost = split[1]
		}
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUsr, dbPass, dbHost, dbName)
	db, err := sql.Open("postgres", connStr)
	utils.PanicIfError(err)

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	return db, err
}
