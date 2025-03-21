package api

import (
	"TaskManager/repo"
	"TaskManager/server/common"
	"TaskManager/server/mw"
	"TaskManager/service"
	"fmt"
	"net/http"
	"time"
)

var Routes = []common.Route{
	{
		Pattern: "/api/users/register",
		Method:  http.MethodPost,
		Handler: mw.Chain(Register, mw.Recovery)},
	{
		Pattern: "/api/users/login",
		Method:  http.MethodPost,
		Handler: mw.Chain(Login, mw.Recovery)},
	{
		Pattern: "/api/users/logout",
		Method:  http.MethodPost,
		Handler: mw.Chain(Logout, mw.Auth, mw.Recovery)},
	{
		Pattern: "/api/tasks/protected",
		Method:  http.MethodGet,
		Handler: mw.Chain(Protected, mw.Auth, mw.Recovery)},
}

func Register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	err := service.RegistrationValidation(username, password)
	if err != nil {
		service.ServeError(w, err)
		return
	}

	err = service.Register(username, password)
	if err != nil {
		service.ServeError(w, err)
		return
	}

	fmt.Fprintf(w, "Registration successfull")
}
func Login(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")

	err := service.LoginValidation(username, password)
	if err != nil {
		service.ServeError(w, err)
		return
	}

	var s *repo.SessionEntity

	s, err = service.Authorize(username, password)
	if err != nil {
		service.ServeError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    s.SessionToken,
		Path:     "/",
		Expires:  s.ExpiresAt.Time,
		HttpOnly: true,
	})

	fmt.Fprintln(w, "Login successful!")
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value("session").(*repo.SessionEntity)

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

func Protected(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "You're accession protected!")
}
