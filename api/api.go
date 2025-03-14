package api

import (
	"TaskManager/repo"
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

type Route struct {
	Pattern string
	Method  string
	Handler func(http.ResponseWriter, *http.Request)
}

var routes = []Route{
	{
		Pattern: "/api/auth/register",
		Method:  http.MethodPost,
		Handler: Chain(register, mw_errorRecovery)},
	{
		Pattern: "/api/auth/login",
		Method:  http.MethodPost,
		Handler: Chain(login, mw_errorRecovery)},
	{
		Pattern: "/api/auth/logout",
		Method:  http.MethodPost,
		Handler: Chain(logout, mw_auth, mw_errorRecovery)},
	{
		Pattern: "/api/tasks/protected",
		Method:  http.MethodGet,
		Handler: Chain(protected, mw_auth, mw_errorRecovery)},
	{
		Pattern: "/signin",
		Method:  http.MethodGet,
		Handler: Chain(signIn, mw_errorRecovery)},
	{
		Pattern: "/dashboard",
		Method:  http.MethodGet,
		Handler: Chain(dashboard, mw_errorRecovery)},
}

func BuildMux() *http.ServeMux {
	mux := http.NewServeMux()
	for _, r := range routes {
		pattern := fmt.Sprintf("%s %s", r.Method, r.Pattern)
		mux.HandleFunc(pattern, r.Handler)
	}
	return mux
}

func Serve() {
	mux := BuildMux()
	fs := http.FileServer(http.Dir("templates"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", mux)
}

func signIn(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/signin.html")
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/dashboard.html")
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
		err := http.StatusNotFound
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
		Path:     "/api",
		Expires:  session.ExpiresAt.Time,
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    session.CsrfToken,
		Path:     "/api",
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
		Path:     "/api",
		Expires:  expirationTime,
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Path:     "/api",
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
