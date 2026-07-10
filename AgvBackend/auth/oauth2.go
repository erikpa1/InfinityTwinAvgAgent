package auth

import (
	"fmt"
	"net/http"
	"turtle/core/lgr"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var (
	auth0Domain  = "dev-30fa1v4j5cfgshau.us.auth0.com"                                // Replace with your Auth0 domains
	clientID     = "0H21KyUGDZgWC7tVVBPTBtGn95eR2FEk"                                 // Replace with your Auth0 Client ID
	clientSecret = "YZGIseNK3VREONygBT8kAKn_1C-CAwII4KkFCNSOWd6kjLsrmOxyWM_XxJi-b5aC" // Replace with your Auth0 Client Secret
	redirectURL  = "http://localhost:8080/api/auth0/callback"
	oauth2Config *oauth2.Config
	state        = "random-state" // Replace with a secure, random value
)

func _RedirectToAuth0(c *gin.Context) {
	if oauth2Config == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OAuth2 configuration not initialized"})
		return
	}

	authURL := oauth2Config.AuthCodeURL("random-state", oauth2.AccessTypeOnline)
	c.Redirect(http.StatusFound, authURL)
}

func _Callback(c *gin.Context) {
	// Verify state parameter
	stateParam := c.Query("state")
	if stateParam != state {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
		return
	}

	// Retrieve authorization code
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code not found"})
		return
	}

	// Exchange the code for tokens
	token, err := oauth2Config.Exchange(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token", "details": err.Error()})
		return
	}

	lgr.InfoJson(token)
}

func Init_api_auth0(r *gin.Engine) {
	oauth2Config = &oauth2.Config{
		ClientID:     clientID,     //credentials.Auth0ClientID(),
		ClientSecret: clientSecret, //credentials.Auth0Secret(),
		RedirectURL:  redirectURL,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("https://%s/authorize", auth0Domain),
			TokenURL: fmt.Sprintf("https://%s/oauth/token", auth0Domain),
		},
	}

	r.GET("/api/auth0/login", _RedirectToAuth0)
	r.GET("/api/auth0/callback", _Callback)
}
