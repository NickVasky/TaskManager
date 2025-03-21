package repo

import (
	"database/sql"
	"log"
)

type SessionEntity struct {
	Id           int
	UserId       int
	SessionToken string
	CsrfToken    string
	ExpiresAt    sql.NullTime
}

func (r *SessionRepo) Create(s *SessionEntity) error {
	query := `
		INSERT INTO sessions (
			user_id, 
			session_token, 
			csrf_token, 
			expires_at)
		VALUES ($1, $2, $3, $4) RETURNING id;`

	err := r.Db.QueryRow(
		query,
		s.UserId,
		s.SessionToken,
		s.CsrfToken,
		s.ExpiresAt).Scan(&s.Id)

	if err != nil {
		return err
	}

	log.Printf("Session created for User %d\n", s.Id)

	return nil
}

func (r *SessionRepo) GetByUserId(user_id int) (*SessionEntity, error) {
	query := `
		SELECT
			id,
			user_id,
			session_token,
			csrf_token,
			expires_at
		FROM sessions
		WHERE user_id = $1;
	`
	var s SessionEntity

	err := r.Db.QueryRow(query, user_id).Scan(
		&s.Id,
		&s.UserId,
		&s.SessionToken,
		&s.CsrfToken,
		&s.ExpiresAt)

	if err != nil {
		return nil, err
	}

	log.Printf("Session with id %d was recieved: %v", user_id, s)

	return &s, nil
}

func (r *SessionRepo) GetByToken(sessionToken string) (*SessionEntity, error) {
	query := `
		SELECT
			id,
			user_id,
			session_token,
			csrf_token,
			expires_at
		FROM sessions
		WHERE session_token = $1;
	`
	var s SessionEntity

	err := r.Db.QueryRow(query, sessionToken).Scan(
		&s.Id,
		&s.UserId,
		&s.SessionToken,
		&s.CsrfToken,
		&s.ExpiresAt)

	if err != nil {
		return nil, err
	}

	log.Printf("Session with id %d was recieved: %v", s.Id, s)

	return &s, nil
}

func (r *SessionRepo) Delete(s *SessionEntity) error {
	query := `
		DELETE FROM sessions
		WHERE id = $1;
	`
	_, err := r.Db.Exec(query, s.Id)

	if err != nil {
		return err
	}

	log.Printf("Session %d was deleted", s.Id)

	return nil
}
