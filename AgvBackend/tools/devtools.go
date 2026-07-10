package tools

import "github.com/gin-gonic/gin"

func IsInDevelopment() bool {
	return gin.Mode() == "debug"
}
