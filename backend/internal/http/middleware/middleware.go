package middleware

import (
	"github.com/Kushian01100111/Tickermaster/internal/domain/auth"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwt *auth.JWTManager
}

func NewAuthMiddleware(jwt *auth.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{jwt: jwt}
}

func (a *AuthMiddleware) RequireAuth() gin.HandlerFunc
