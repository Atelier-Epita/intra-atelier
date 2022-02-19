package main

import (
	"intra/api"
	"intra/cmd"
	"intra/db"
)

var user_table = `
CREATE TABLE IF NOT EXISTS user (
    id int  not null auto_increment primary key,
    login VARCHAR(255) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    promotion smallint NOT NULL
);`

var account_table = `CREATE TABLE IF NOT EXISTS account (
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
	cmd.Init()
	db.Connect()
	defer db.Close()

	var server = api.CreateRouter()
	server.Run()

	db.DB.MustExec(user_table)
	db.DB.MustExec(account_table)
}

