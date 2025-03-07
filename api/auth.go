package api

import (
	"TaskManager/repo"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNoCsrf         = errors.New("no CSRF token provided")
	ErrNoSession      = errors.New("session not found or expired")
	ErrSessionExpired = errors.New("session expired")
)

const hashCost = 11

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	return string(hash), err
}

func checkPassword(password, hashhedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashhedPassword), []byte(password))
}

func generateToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal("Unable to generate hash")
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

func Authorize(r *http.Request) (*repo.SessionEntity, error) {
	session_cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil, err
	}

	csrf_token := r.Header.Get("X-CSRF-Token")
	if csrf_token == "" {
		return nil, ErrNoCsrf
	}

	db := repo.NewConnection()
	defer db.Close()
	sr := db.NewSessionRepo()
	session, err := sr.GetByTokens(session_cookie.Value, csrf_token)

	if err != nil || session == nil {
		return nil, ErrNoSession
	}

	if session.ExpiresAt.Valid && session.ExpiresAt.Time.Before(time.Now()) {
		log.Printf("Session %d expired", session.Id)
		sr.Delete(session)
		return nil, ErrSessionExpired
	} else {
		log.Printf("Session %d is valid", session.Id)
		return session, nil
	}
}
