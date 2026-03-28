package handlers

import (
	"net/http"
	"time"

	"github.com/Kushian01100111/Tickermaster/internal/app/auth"
	"github.com/Kushian01100111/Tickermaster/internal/http/dto"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	app           auth.AuthService
	refreshCookie string
	refreshMaxAge int
	cookieDomain  string
	secureCookies bool
	sameSite      http.SameSite
}

func NewAuthHandler(svc auth.AuthService) *AuthHandler {
	return &AuthHandler{
		app:           svc,
		refreshCookie: "refresh_cookie",
		refreshMaxAge: int((30 * 24 * time.Hour).Seconds()),
		cookieDomain:  "",
		secureCookies: false,
		sameSite:      http.SameSiteLaxMode,
	}
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
	var req *dto.LoginRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind body of request"})
		return
	}

	ctx, cancel := generateCtx() // Cambiar a g.Request.Context
	defer cancel()

	authRes, err := a.app.Login(ctx, auth.LoginParams{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	a.setRefreshCookie(g, authRes.RefreshToken)

	g.JSON(http.StatusCreated, authRes)
}

func (a *AuthHandler) refresh(g *gin.Context) {
	rt, err := g.Cookie(a.refreshCookie)

	if err != nil {
		g.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})

		return
	}

	authRes, err := a.app.Refresh(g.Request.Context(), rt)
	if err != nil {
		g.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	a.setRefreshCookie(g, authRes.RefreshToken)

	g.JSON(http.StatusOK, authRes)
}

func (a *AuthHandler) logout(g *gin.Context) {
	rt, err := g.Cookie(a.refreshCookie)
	if err != nil {
		g.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if rt != "" {
		if err := a.app.Logout(g.Request.Context(), rt); err != nil {
			g.JSON(http.StatusUnauthorized, err.Error())
			return
		}
	}

	a.clearRefreshCookie(g)
	g.Status(http.StatusNoContent)
}

func (a *AuthHandler) signupRequest(g *gin.Context) {
	var req *dto.RequestCode

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.app.SignupRequest(g.Request.Context(), req.Email); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	g.Status(http.StatusNoContent)
}

func (a *AuthHandler) signupVerify(g *gin.Context) {
	var req *dto.VerifyRequest

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authRes, err := a.app.SignupVerify(g.Request.Context(), auth.VerifyParams{
		Email: req.Email,
		Code:  req.Code,
	})
	if err != nil {
		g.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	a.setRefreshCookie(g, authRes.RefreshToken)
	g.JSON(http.StatusOK, authRes)
}

func (a *AuthHandler) loginRequest(g *gin.Context) {
	var req *dto.RequestCode

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.app.LoginRequest(g.Request.Context(), req.Email); err != nil {
		g.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	g.Status(http.StatusNoContent)
}

func (a *AuthHandler) loginVerify(g *gin.Context) {
	var req *dto.VerifyRequest

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authRes, err := a.app.LoginVerify(g.Request.Context(), auth.VerifyParams{
		Email: req.Email,
		Code:  req.Code,
	})
	if err != nil {
		g.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	a.setRefreshCookie(g, authRes.RefreshToken)
	g.JSON(http.StatusOK, authRes)
}

func (a *AuthHandler) setRefreshCookie(g *gin.Context, refreshToken string) {
	http.SetCookie(g.Writer, &http.Cookie{
		Name:     a.refreshCookie,
		Value:    refreshToken,
		Path:     "/auth",
		Domain:   a.cookieDomain,
		MaxAge:   a.refreshMaxAge,
		HttpOnly: true,
		Secure:   a.secureCookies,
		SameSite: a.sameSite,
	})
}

func (a *AuthHandler) clearRefreshCookie(g *gin.Context) {
	http.SetCookie(g.Writer, &http.Cookie{
		Name:     a.refreshCookie,
		Value:    "",
		Path:     "/auth",
		Domain:   a.cookieDomain,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   a.secureCookies,
		SameSite: a.sameSite,
	})
}
