package handlers

import (
	"github.com/Kushian01100111/Tickermaster/internal/app/auth"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	app auth.AuthService
}

func NewAuthHandler(svc auth.AuthService) *AuthHandler {
	return &AuthHandler{app: svc}
}

func (a *AuthHandler) AuthRoutes(r *gin.RouterGroup) {
	context := r.Group("/auth")
	{
		context.POST("/login", a.login)
		context.POST("/refresh", a.refresh)
		context.POST("/logout", a.logout)
		context.POST("/signup/request-code", a.signupRequest)
		context.POST("/signup/verify-code", a.signupVerify)
		context.POST("/login/request-code", a.loginRequest)
		context.POST("/login/verify-code", a.loginVerify)
	}
}

func (a *AuthHandler) login(g *gin.Context) {

}

func (a *AuthHandler) refresh(g *gin.Context) {

}

func (a *AuthHandler) logout(g *gin.Context) {

}

func (a *AuthHandler) signupRequest(g *gin.Context) {

}

func (a *AuthHandler) signupVerify(g *gin.Context) {

}

func (a *AuthHandler) loginRequest(g *gin.Context) {

}

func (a *AuthHandler) loginVerify(g *gin.Context) {

}
