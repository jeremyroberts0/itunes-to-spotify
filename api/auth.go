package api

import (
	"encoding/json"
	"net/http"
	"os"

	"golang.org/x/oauth2"

	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify"
)

const (
	redirectURI = "http://localhost:8081/authorized"
	state       = "secret-state-token"
	cookieName  = "spotify-token"
)

var auth spotify.Authenticator

func init() {
	auth = spotify.NewAuthenticator(
		redirectURI,
		spotify.ScopeUserReadPrivate,
		spotify.ScopePlaylistModifyPrivate,
	)
	auth.SetAuthInfo(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
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
	url := auth.AuthURL(state)

	router.GET("/authorize", func(context *gin.Context) {
		context.Redirect(http.StatusTemporaryRedirect, url)
	})

	router.GET("/authorized", func(context *gin.Context) {
		token, err := auth.Token(state, context.Request)
		if err != nil {
			context.IndentedJSON(500, map[string]string{
				"message": "Couldn't get token",
				"error":   err.Error(),
			})
			return
		}

		cookieValue, err := json.Marshal(token)
		if err != nil {
			context.IndentedJSON(500, map[string]string{
				"message": "Error creating auth cookie",
				"error":   err.Error(),
			})
			return
		}
		context.SetCookie(
			cookieName,
			string(cookieValue),
			60*60*1, // 1 hour
			"/",
			"",
			false,
			true,
		)
		context.IndentedJSON(http.StatusOK, map[string]string{
			"message": "Authentication successful",
		})
	})
}
