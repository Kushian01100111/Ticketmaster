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

// RequireRole
func RequireRole(allowed ...string) gin.HandlerFunc {
	allowedSet := map[string]struct{}{}
	for _, r := range allowed {
		allowedSet[r] = struct{}{}
	}

	return func(g *gin.Context) {
		role := g.GetString("role")
		if role == "" {
			g.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if _, ok := allowedSet[role]; !ok {
			g.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		g.Next()
	}
}

// RequireAnyScope verifica que el token tenga al menos uno de los scopes resqueridos.
func RequireAnyScope(required ...string) gin.HandlerFunc {
	req := map[string]struct{}{}
	for _, s := range required {
		req[s] = struct{}{}
	}

	return func(g *gin.Context) {
		any, ok := g.Get("scopes")
		if !ok {
			g.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		scopes, ok := any.([]string)
		if !ok {
			g.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		for _, s := range scopes {
			if _, ok := req[s]; ok {
				g.Next()
				return
			}
		}
		g.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	}
}
