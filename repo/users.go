package repo

import (
	"database/sql"
	"log"
	"time"
)

type UserEntity struct {
	Id         int
	Username   string
	Password   string
	FirstName  sql.NullString
	SecondName sql.NullString
	CreatedAt  sql.NullTime
}

func (r *UserRepo) GetById(id int) (*UserEntity, error) {
	query := `
		SELECT
			id,
			username,
			password,
			first_name,
			second_name,
			created_at 
		FROM users
		WHERE id = $1;
	`
	var u UserEntity

	err := r.Db.QueryRow(query, id).Scan(
		&u.Id,
		&u.Username,
		&u.Password,
		&u.FirstName,
		&u.SecondName,
		&u.CreatedAt)

	if err != nil {
		return nil, err
	}

	log.Printf("User with id %d was recieved: %v", id, u)

	return &u, nil
}

func (r *UserRepo) GetByUsername(username string) (*UserEntity, error) {
	query := `
		SELECT
			id,
			username,
			password,
			first_name,
			second_name,
			created_at 
		FROM users
		WHERE username = $1;
	`
	var u UserEntity

	err := r.Db.QueryRow(query, username).Scan(
		&u.Id,
		&u.Username,
		&u.Password,
		&u.FirstName,
		&u.SecondName,
		&u.CreatedAt)

	if err != nil {
		return nil, err
	}

	log.Printf("User with username %s was recieved: %v", username, u)

	return &u, nil
}

func (r *UserRepo) GetAll() ([]UserEntity, error) {
	query := `
		SELECT
			id,
			username,
			password,
			first_name,
			second_name,
			created_at 
		FROM users;
	`

	var users []UserEntity

	rows, err := r.Db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u UserEntity
		err := rows.Scan(
			&u.Id,
			&u.Username,
			&u.Password,
			&u.FirstName,
			&u.SecondName,
			&u.CreatedAt)

		if err != nil {
			log.Fatal(err)
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepo) Create(u *UserEntity) error {
	query := `
		INSERT INTO users (
			username, 
			password, 
			first_name, 
			second_name, 
			created_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	u.CreatedAt.Time = time.Now()
	u.CreatedAt.Valid = true

	err := r.Db.QueryRow(
		query,
		u.Username,
		u.Password,
		u.FirstName,
		u.SecondName,
		u.CreatedAt).Scan(&u.Id)

	if err != nil {
		return err
	}

	log.Printf("User %s created with id %d\n", u.Username, u.Id)

	return nil
}

func (r *UserRepo) Edit(u *UserEntity) error {
	query := `
		UPDATE users
		SET 
			username = $2, 
			password = $3, 
			first_name = $4, 
			second_name = $5,
			created_at = $6
		WHERE id = $1;
	`

	_, err := r.Db.Exec(query, u.Id, u.Username, u.Password, u.FirstName, u.SecondName, u.CreatedAt)

	if err != nil {
		return err
	}
	log.Printf("User %d was updated: %v", u.Id, u)

	return nil
}

func (r *UserRepo) Delete(u *UserEntity) error {
	query := `
		DELETE FROM users
		WHERE id = $1;
	`
	_, err := r.Db.Exec(query, u.Id)

	if err != nil {
		return err
	}

	log.Printf("User %d was deleted", u.Id)

	return nil
}
