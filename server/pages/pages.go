package pages

import (
	"TaskManager/repo"
	"TaskManager/server/common"
	"TaskManager/server/mw"
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseFiles("templates/auth.html", "templates/dashboard.html"))

var Routes = []common.Route{
	{
		Pattern: "/signin",
		Method:  http.MethodGet,
		Handler: mw.Chain(signIn, mw.Recovery)},
	{
		Pattern: "/signup",
		Method:  http.MethodGet,
		Handler: mw.Chain(signUp, mw.Recovery)},
	{
		Pattern: "/dashboard",
		Method:  http.MethodGet,
		Handler: mw.Chain(dashboard, mw.Auth, mw.Recovery)},
}

type AuthPage struct {
	PageTitle       string
	FormId          string
	FormTitle       string
	PrimaryButton   HtmlButton
	SecondaryButton HtmlButton
}

type HtmlButton struct {
	Text, Href string
}

func renderTemplate(w http.ResponseWriter, tmpl string, p any) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func signIn(w http.ResponseWriter, r *http.Request) {
	authPage := AuthPage{
		PageTitle:       "Sign In",
		FormId:          "signin-form",
		FormTitle:       "Sign In",
		PrimaryButton:   HtmlButton{Text: "Sign In"},
		SecondaryButton: HtmlButton{Text: "Sign Up", Href: "signup"},
	}
	renderTemplate(w, "auth", authPage)
}

func signUp(w http.ResponseWriter, r *http.Request) {
	authPage := AuthPage{
		PageTitle:       "Sign Up",
		FormId:          "signup-form",
		FormTitle:       "Sign Up",
		PrimaryButton:   HtmlButton{Text: "Sign Up"},
		SecondaryButton: HtmlButton{Text: "Sign In", Href: "signup"},
	}
	renderTemplate(w, "auth", &authPage)
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(mw.SessionKey).(*repo.SessionEntity)
	conn := repo.NewConnection()

	if !ok {
		panic("No session data found in context")
	}

	user, _ := conn.NewUserRepo().GetById(session.UserId)

	renderTemplate(w, "dashboard", user)
}
