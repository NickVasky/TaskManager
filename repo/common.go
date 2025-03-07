package repo

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var (
	DbHost     = "localhost"
	DbPort     = 25912
	DbUsername = "go"
	DbPassword = "some_secure_passworld"
	DbName     = "task_manager_db"
)

type Database struct {
	Db *sql.DB
}

type UserRepo struct {
	Db *sql.DB
}

type TaskRepo struct {
	Db *sql.DB
}

type SessionRepo struct {
	Db *sql.DB
}

func NewConnection() *Database {
	psqlConnStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DbHost, DbPort, DbUsername, DbPassword, DbName)

	db, err := sql.Open("postgres", psqlConnStr)

	if err != nil {
		return nil
	}

	return &Database{db}
}

func (db *Database) Close() {
	if err := db.Db.Close(); err != nil {
		log.Println("Error while closing DB connection: ", err)
	}
}

func (d *Database) NewUserRepo() *UserRepo {
	return &UserRepo{d.Db}
}

func (d *Database) NewTaskRepo() *TaskRepo {
	return &TaskRepo{d.Db}
}

func (d *Database) NewSessionRepo() *SessionRepo {
	return &SessionRepo{d.Db}
}
