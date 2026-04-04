package middleware

import (
	"net/http"

	"github.com/Kushian01100111/Tickermaster/internal/domain/session"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwt *session.JWTManager
}

func NewAuthMiddleware(jwt *session.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{jwt: jwt}
}

// Protege rutas validando access tokens y carga claims en el gin.Context
func (a *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(g *gin.Context) {
		tokenStr, err := session.ExtractBearerToken(g.GetHeader("Authorization"))
		if err != nil {
			g.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
			return
		}

		claims, err := a.jwt.ParseAndValidate(tokenStr)
		if err != nil {
			g.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorize"})
			return
		}

		g.Set("userID", claims.Subject)
		g.Set("role", claims.Role)
		g.Set("scopes", claims.Scopes)

		g.Next()
	}
}
