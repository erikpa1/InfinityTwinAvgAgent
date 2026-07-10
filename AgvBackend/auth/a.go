package auth

import (
	"fmt"
	"net/http"
	"strconv"
	"turtle/core/users"
	"turtle/credentials"
	"turtle/ctrl"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type AuthContext struct {
	GinC *gin.Context
	User *users.User
}

func GetUserFromContext(c *gin.Context) *users.User {
	tmp := users.NewSuperAdmin()

	return tmp
}

func GetUserFromCookies(cookiestring string) *users.User {
	return &users.User{}
}

func LoginRequiredWithUser(fn func(ctx *AuthContext)) gin.HandlerFunc {

	return func(c *gin.Context) {
		str_token, err := c.Cookie("")

		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"why": "no auth token"})
		} else {
			_, err := jwt.Parse(str_token, func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(credentials.AuthInfinityJwtSecret()), nil
			})

			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{"why": "Wrong auth token"})
				return
			}

			ctx := AuthContext{}
			ctx.GinC = c
			ctx.User = users.NewUser()
			fn(&ctx)
		}

	}
}

func IsLoggedInInfinity(c *gin.Context) (*users.User, string) {
	tokenName := credentials.AuthInfinityJwtKey()

	token, err := c.Cookie(tokenName)
	if err != nil {
		token = "" // If the cookie is not found, use an empty string
	}

	// Check the token and retrieve the user
	user, _ := ctrl.CheckInfinityAuth(token)

	if user != nil {
		tokenName := credentials.AuthInfinityJwtKey()

		var expireTime, _ = strconv.Atoi(credentials.AuthInfinityTokenExpire())

		c.SetCookie(tokenName, token, expireTime, "/", "", false, true)

		return user, token
	}

	c.SetCookie(tokenName, "", -1, "/", "", false, true)

	return nil, token
}

type IdpImplementation interface {
	RedirectToIdp(c *gin.Context)
}

func Rbac(module, role string) func(c *gin.Context) {

	return func(c *gin.Context) {
		LoginRequired(c)
	}

}

func LoginRequired(c *gin.Context) {
	if credentials.AuthProvider() == credentials.AUTH_PROVIDER_NONE {
		c.Next()

		return
	}

	if credentials.AuthProvider() == credentials.AUTH_PROVIDER_INFINITY || credentials.AuthProvider() == credentials.AUTH_PROVIDER_ALL {
		user, _ := IsLoggedInInfinity(c)

		if user != nil {
			c.Next()

			return
		}
	}

	c.JSON(http.StatusUnauthorized, gin.H{"why": "Unauthorized"})
}

func AdminRequired(c *gin.Context) {
	c.Next()
}

func LoginOrApp(c *gin.Context) {
	c.Next()
}
