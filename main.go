package main

import (
	"fmt"
	"intra/db"
	"strings"
)

var schema = `
CREATE TABLE IF NOT EXISTS user (
    id int  not null auto_increment primary key,
    login VARCHAR(255) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    promotion smallint NOT NULL
);

CREATE TABLE IF NOT EXISTS account (
    id int not null auto_increment primary key,
    user_id int not null, 
    foreign key (user_id) references user(id)
);`

type User struct {
	Login     string `db:"login"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Promotion string `db:"promotion"`
}

type Account struct {
	Id int `db:"id"`
}

func main() {
	fmt.Println("Entering main...")
	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatal("Error loading .env file")
	//}
	db.Connect()
	fmt.Println("Connected")
	db.DB.MustExec(schema)
	fmt.Println("schema executed")
	insert_user("joe.mama", 2024)
	fmt.Println("End")
}

func insert_user(login string, promotion uint16) error {
	substrings := strings.Split(login, ".")
	if len(substrings) != 2 {
		fmt.Println("insert_user: Login incorrect")
		return nil
	}

	tx := db.DB.MustBegin()
	tx.MustExec("INSERT INTO user (login, first_name, last_name, promotion) VALUES (?, ?, ?, ?)",
		login,
		substrings[0],
		substrings[1],
		promotion,
	)

	if tx.Commit() != nil {
		return fmt.Errorf("cannot insert user %v, promotion %v", login, promotion)
	}

	return nil
}
