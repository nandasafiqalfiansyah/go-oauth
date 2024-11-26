package api

// local package swich main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type App struct {
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

var (
	router *gin.Engine
)

func init() {
	// Initialize Gin router
	router = gin.New()
	router.Use(gin.Logger())

	// env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Configure OAuth2
	conf := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8000/auth/callback",
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

	app := App{config: conf}
	// Define routes
	router.GET("/", func(c *gin.Context) {
		app.oauthHandler(c.Writer, c.Request)
	})
	router.GET("/auth/callback", func(c *gin.Context) {
		app.callbackHandler(c.Writer, c.Request)
	})
}


func (a *App) oauthHandler(w http.ResponseWriter, r *http.Request) {
	url := a.config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

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

	json.NewEncoder(w).Encode(userInfo)
}

// local
// func main() {
// 	if err := router.Run(":8000"); err != nil {
// 		log.Fatal("Failed to start server:", err)
// 	}
// }

func Handler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}
