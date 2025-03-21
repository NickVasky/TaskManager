package service

import (
	"TaskManager/repo"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const hashCost = 11

type HttpError struct {
	Code    int
	Message string
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("%v - %s", e.Code, e.Message)
}

var (
	ErrNoSession      = errors.New("session not found or expired")
	ErrSessionExpired = errors.New("session expired")
)

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

func SessionAuth(r *http.Request) (*repo.SessionEntity, error) {
	session_cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil, err
	}

	db := repo.NewConnection()
	defer db.Close()
	sr := db.NewSessionRepo()
	session, err := sr.GetByToken(session_cookie.Value)

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

func Authorize(username, password string) (*repo.SessionEntity, error) {
	db := repo.NewConnection()
	ur := db.NewUserRepo()
	defer db.Close()

	user, _ := ur.GetByUsername(username)
	if user == nil {
		return nil, &HttpError{
			Code:    http.StatusNotFound,
			Message: "User not found",
		}
	}

	if err := checkPassword(password, user.Password); err != nil {
		return nil, &HttpError{
			Code:    http.StatusNotFound,
			Message: "Bad password",
		}
	}

	session := repo.SessionEntity{
		UserId:       user.Id,
		SessionToken: generateToken(32),
		//CsrfToken:    generateToken(32),
		ExpiresAt: sql.NullTime{Time: time.Now().Add(time.Hour * 24), Valid: true},
	}

	sr := db.NewSessionRepo()
	sr.Create(&session)

	return &session, nil
}

func Register(username, password string) error {
	db := repo.NewConnection()
	ur := db.NewUserRepo()
	defer db.Close()
	user, _ := ur.GetByUsername(username)

	if user != nil {
		return &HttpError{
			Code:    http.StatusConflict,
			Message: "User already exists",
		}
	}

	password, _ = hashPassword(password)
	user = &repo.UserEntity{
		Username: username,
		Password: password}
	ur.Create(user)

	return nil
}

func RegistrationValidation(username, password string) error {
	if l := len(username); l < 4 || l > 100 {
		return &HttpError{
			Code:    http.StatusNotAcceptable,
			Message: "Invalid username",
		}
	}

	if l := len(password); l < 8 || l > 64 {
		return &HttpError{
			Code:    http.StatusNotAcceptable,
			Message: "Invalid password",
		}
	}

	return nil
}

func LoginValidation(username, password string) error {
	if len(username) < 1 || len(password) < 1 {
		return &HttpError{
			Code:    http.StatusNotAcceptable,
			Message: "Username/password not provided",
		}
	}

	return nil
}

func ServeError(w http.ResponseWriter, err error) {
	var httpErr *HttpError

	if errors.As(err, &httpErr) {
		http.Error(w, httpErr.Message, httpErr.Code)
		return
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
