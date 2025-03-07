package api

import (
	"TaskManager/repo"
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

var patternMethods = map[string]string{
	"/users/register":  http.MethodPost,
	"/users/login":     http.MethodPost,
	"/users/logout":    http.MethodPost,
	"/tasks/protected": http.MethodGet,
}

func Serve() {
	http.HandleFunc("/users/register", Chain(register, mw_methodCheck, mw_errorRecovery))
	http.HandleFunc("/users/login", Chain(login, mw_methodCheck, mw_errorRecovery))
	http.HandleFunc("/users/logout", Chain(logout, mw_auth, mw_methodCheck, mw_errorRecovery))
	http.HandleFunc("/tasks/protected", Chain(protected, mw_auth, mw_methodCheck, mw_errorRecovery))

	http.ListenAndServe(":8080", nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if l := len(username); l < 4 || l > 100 {
		err := http.StatusNotAcceptable
		http.Error(w, "Invalid Username", err)
		return
	}
	if l := len(password); l < 8 || l > 64 {
		err := http.StatusNotAcceptable
		http.Error(w, "Invalid Password", err)
		return
	}

	db := repo.NewConnection()
	ur := db.NewUserRepo()
	defer db.Close()
	user, _ := ur.GetByUsername(username)
	if user != nil {
		err := http.StatusConflict
		http.Error(w, "User already exists", err)
		return
	}

	password, _ = hashPassword(password)
	user = &repo.UserEntity{
		Username: username,
		Password: password}
	ur.Create(user)

	fmt.Fprintf(w, "Registration successfull")
}
func login(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")

	if l := len(username); l < 1 {
		err := http.StatusNotAcceptable
		http.Error(w, "Username not provided", err)
		return
	}
	if l := len(password); l < 1 {
		err := http.StatusNotAcceptable
		http.Error(w, "Password not provided", err)
		return
	}

	db := repo.NewConnection()
	ur := db.NewUserRepo()
	defer db.Close()
	user, _ := ur.GetByUsername(username)
	if user == nil {
		err := http.StatusNotFound
		http.Error(w, "User not found", err)
		return
	}

	if err := checkPassword(password, user.Password); err != nil {
		err := http.StatusOK
		http.Error(w, "Bad password", err)
		return
	}

	session := repo.SessionEntity{
		UserId:       user.Id,
		SessionToken: generateToken(32),
		CsrfToken:    generateToken(32),
		ExpiresAt:    sql.NullTime{Time: time.Now().Add(time.Hour * 24), Valid: true},
	}

	sr := db.NewSessionRepo()
	sr.Create(&session)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.SessionToken,
		Expires:  session.ExpiresAt.Time,
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    session.CsrfToken,
		Expires:  session.ExpiresAt.Time,
		HttpOnly: false,
	})

	fmt.Fprintln(w, "Login successful!")
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(SessionKey).(*repo.SessionEntity)

	if !ok {
		panic("No session data found in context")
	}

	expirationTime := time.Now().Add(-time.Hour)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  expirationTime,
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  expirationTime,
		HttpOnly: false,
	})

	db := repo.NewConnection()
	sr := db.NewSessionRepo()

	sr.Delete(session)

	fmt.Fprintln(w, "Logout successful!")

}

func protected(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "You're accession protected!")
}
