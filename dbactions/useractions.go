package dbactions

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var (
	DbHost     = "localhost"
	DbPort     = 5432
	DbUsername = "myuser"
	DbPassword = "mypassword"
	DbName     = "mydb"
)

type UserEntity struct {
	Username, Password, Name, Surname string
}

type TaskEntity struct {
	Goal, Measure, Relevance, Deadline string
}

func openConnection() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DbHost, DbPort, DbUsername, DbPassword, DbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	} else {
		return db, nil
	}
}
func userExists(username string) (int, error) {
	db, err := openConnection()

	if err != nil {
		log.Println(err)
		return -1, err
	}
	defer db.Close()

	userId := -1
	query := `
		SELECT id FROM users
		WHERE username = $1;`

	rows, err := db.Query(query, username)

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return -1, err
		}
		userId = id

	}
	return userId, nil

}

func CreateUser(u UserEntity) (int, error) {
	db, err := openConnection()

	if err != nil {
		log.Println(err)
		return -1, err
	}

	defer db.Close()

	foundId, err := userExists(u.Username)
	if err != nil {
		return -1, err
	}

	if foundId > 0 {
		errMsg := fmt.Sprintf("User %s already exists", u.Username)
		log.Println(errMsg)
		return -1, errors.New(errMsg)
	}

	query := `
		INSERT INTO users (username, password, name, surname, createdat)
		VALUES ($1, $2, $3, $4, now()) RETURNING id;`

	row := db.QueryRow(query, u.Username, u.Password, u.Name, u.Surname)

	userId := -1
	if err := row.Scan(&userId); err != nil {
		panic(err)
	}
	log.Printf("User %s created with id %d\n", u.Username, userId)

	return userId, nil
}

func ChangePassword(userId int, oldPwd, newPwd string) error {
	db, err := openConnection()

	if err != nil {
		log.Println(err)
		return err
	}

	defer db.Close()

	query := `
		SELECT password FROM users
		WHERE id = $1;
	`
	row := db.QueryRow(query, userId)

	var oldPwdFromDb string
	if err := row.Scan(&oldPwdFromDb); err != nil {
		errMsg := fmt.Sprintf("User %d not found", userId)
		log.Println(errMsg)
		return errors.New(errMsg)
	}

	if oldPwd != oldPwdFromDb {
		errMsg := fmt.Sprintf("User %d provided incorrect password", userId)
		log.Println(errMsg)
		return errors.New(errMsg)
	}

	query = `
		UPDATE users
		SET password = $1
		WHERE id = $2;
	`

	db.Exec(query, newPwd, userId)
	log.Printf("Password changed for user %d\n", userId)

	return nil
}
