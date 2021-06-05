package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var host = os.Getenv("HOST")
var port = os.Getenv("PORT")
var user = os.Getenv("USER")
var password = os.Getenv("PASSWORD")
var dbname = os.Getenv("DBNAME")
var sslmode = os.Getenv("SSLMODE")

var dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

func CollectData(chat_id int, username string, message string) error {
	db, err := sql.Open("postgres", dbInfo)

	if err != nil {
		return err
	}
	defer db.Close()

	//Создаем SQL запрос
	data := `INSERT INTO questions(chat_id, username, message) VALUES($1, $2, $3);`

	//Выполняем наш SQL запрос
	if _, err = db.Exec(data, chat_id, username, message); err != nil {
		return err
	}

	return nil
}

func GetAnswer(id int) string {
	var answer string

	db, err := sql.Open("postgres", dbInfo)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	q := fmt.Sprintf("SELECT answer FROM answers WHERE id = %d", id)

	row := db.QueryRow(q)
	err = row.Scan(&answer)

	if err != nil {
		log.Fatal(err)
	}

	return answer
}
