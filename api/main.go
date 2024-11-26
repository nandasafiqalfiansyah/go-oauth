package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type App struct{
	config *oauth2.Config
}
type UserInfo struct {
    ID            string `json:"id"`
    Email         string `json:"email"`
    VerifiedEmail bool   `json:"verified_email"`
    Name          string `json:"name"`
    GivenName     string `json:"given_name"`
    FamilyName    string `json:"family_name"`
    Picture       string `json:"picture"`
    Locale        string `json:"locale"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "<h1>You SUCCESS DEPLOY GOLANG page 1</h1>")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	conf := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8000/auth/callback",
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

	app := App{config: conf}
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.loginHandler)
	mux.HandleFunc("/auth/oauth", app.oauthHandler)
	mux.HandleFunc("/auth/callback", app.callbackHandler)
	http.ListenAndServe(":8000", mux)
}

// login 
func (a *App) loginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// oauth
func (a *App) oauthHandler(w http.ResponseWriter, r *http.Request) {
	url := a.config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// callback
func (a *App) callbackHandler(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")
    if code == "" {
        http.Error(w, "Authorization code not found", http.StatusBadRequest)
        return
    }

    t, err := a.config.Exchange(context.Background(), code)
    if err != nil {
        http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusBadRequest)
        return
    }

    client := a.config.Client(context.Background(), t)
    resp, err := client.Get("https://openidconnect.googleapis.com/v1/userinfo")
    if err != nil {
        http.Error(w, "Failed to fetch user info: "+err.Error(), http.StatusBadRequest)
        return
    }
    defer resp.Body.Close()

    var userInfo UserInfo
    err = json.NewDecoder(resp.Body).Decode(&userInfo)
    if err != nil {
        http.Error(w, "Failed to decode user info: "+err.Error(), http.StatusInternalServerError)
        return
    }

    tmpl, err := template.ParseFiles("templates/dashboard.html")
    if err != nil {
        http.Error(w, "Failed to load template: "+err.Error(), http.StatusInternalServerError)
        return
    }

    tmpl.Execute(w, userInfo)
}