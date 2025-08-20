package api

import (
	"encoding/json"
	"net/http"
	"os"

	"golang.org/x/oauth2"

	"github.com/gin-gonic/gin"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const (
	redirectURI = "http://localhost:8081/authorized"
	state       = "secret-state-token"
	cookieName  = "spotify-token"
)

var authenticator *spotifyauth.Authenticator

func init() {
	authenticator = spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopePlaylistModifyPrivate),
		spotifyauth.WithClientID(os.Getenv("SPOTIFY_CLIENT_ID")),
		spotifyauth.WithClientSecret(os.Getenv("SPOTIFY_CLIENT_SECRET")),
	)
}

func getCookieToken(c *gin.Context) (token *oauth2.Token, err error) {
	notJwtStr, err := c.Cookie(cookieName)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(notJwtStr), &token)

	return
}

// Auth sets up auth routes
func Auth(router *gin.Engine) {
	url := authenticator.AuthURL(state)

	router.GET("/authorize", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, url)
	})

	router.GET("/authorized", func(c *gin.Context) {
		token, err := authenticator.Token(c.Request.Context(), state, c.Request)
		if err != nil {
			c.IndentedJSON(500, map[string]string{
				"message": "Couldn't get token",
				"error":   err.Error(),
			})
			return
		}

		cookieValue, err := json.Marshal(token)
		if err != nil {
			c.IndentedJSON(500, map[string]string{
				"message": "Error creating auth cookie",
				"error":   err.Error(),
			})
			return
		}
		c.SetCookie(
			cookieName,
			string(cookieValue),
			60*60*1, // 1 hour
			"/",
			"",
			false,
			true,
		)
		c.IndentedJSON(http.StatusOK, map[string]string{
			"message": "Authentication successful",
		})
	})
}
