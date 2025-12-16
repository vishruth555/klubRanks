package middlewares

import (
	"klubRanks/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Aunthenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")

	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "no auth provided"})
		return
	}
	// Expected format: "Bearer <token>"
	parts := strings.SplitN(token, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		return
	}

	token = parts[1] // <-- this is your actual token without "Bearer"
	userId, err := utils.VerifyToken(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "not authorized, " + err.Error()})
		return
	}
	context.Set("userId", userId)
	context.Next()
}
