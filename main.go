package main

import (
	"fmt"
	"strings"
)

func main() {
	insert_user("joe.mama", 2024)
}

func insert_user(login string, promotion uint16) error {
	substrings := strings.Split(login, ".")
	if len(substrings) != 2 {
		fmt.Println("insert_user: Login incorrect")
		return nil
	}
	_, err := db.Query("INSERT INTO user (login, first_name, last_name, promotion) VALUES (?, ?, ?, ?)",
		login,
		substrings[0],
		substrings[1],
		promotion,
	)

	if err != nil {
		return fmt.Errorf("cannot insert user %v, promotion %v", login, promotion)
	}

	return nil
}
