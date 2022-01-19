package db

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	DB *sqlx.DB
)

func Connect() {
	connect(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))
}

func connect(host, user, pass, database string) {
	var err error
	//var arg = user + ":" + pass + "@(" + host + ")/" + database + "?charset=utf8mb4,utf8&parseTime=true&loc=Europe%2FParis&time_zone=%27Europe%2FParis%27"
	var arg = "root:root@(localhost)/intradb"
	fmt.Println(arg)
	DB, err = sqlx.Connect("mysql", arg)
	if err != nil {
		fmt.Println(err)
		zap.S().Fatal("could not connect to database.")
	}
	zap.S().Info("Connected to database.")
}

func Close() {
	if DB == nil {
		return
	}
	zap.S().Info("Closing DB...")
	err := DB.Close()
	if err != nil {
		zap.S().Error()
		return
	}
	zap.S().Info("Closed DB")
}
