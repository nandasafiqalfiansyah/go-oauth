package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
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
	router = gin.New()
	router.Use(gin.Logger())
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
	router.GET("/auth", app.oauthHandler)
	router.GET("/auth/callback", app.callbackHandler)
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "server is running",
		})
	})
}	

func (a *App) oauthHandler(c *gin.Context) {
	urlFe := c.Query("url")
	if urlFe == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing URL parameter"})
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:  "urlFe",
		Value: urlFe,
		Path:  "/",
		HttpOnly: true,
	})
	url := a.config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (a *App) callbackHandler(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code not found"})
		return
	}
	t, err := a.config.Exchange(context.Background(), code)
	if err != nil {
		log.Println("Error exchanging token:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange token"})
		return
	}
	client := a.config.Client(context.Background(), t)
	resp, err := client.Get("https://openidconnect.googleapis.com/v1/userinfo")
	if err != nil {
		log.Println("Error fetching user info:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch user info"})
		return
	}
	defer resp.Body.Close()
	var userInfo UserInfo
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		log.Println("Error decoding user info:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info"})
		return
	}
	cookie, err := c.Request.Cookie("urlFe")
	if err != nil {
		log.Println("Error retrieving urlFe from cookie:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL not found in session"})
		return
	}
	userData, err := json.Marshal(userInfo)
	if err != nil {
		log.Println("Error encoding user info:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store user info"})
		return
	}
	redirectURL := fmt.Sprintf("%s?user=%s", cookie.Value, url.QueryEscape(string(userData)))
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}


// func main() {
// 	if err := router.Run(":8000"); err != nil {
// 		log.Fatal("Failed to start server:", err)
// 	}
// }

func Handler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}
