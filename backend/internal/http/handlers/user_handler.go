package handlers

import (
	"github.com/Kushian01100111/Tickermaster/internal/app/user"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	app user.UserService
}

func NewUserHandler(svc user.UserService) *UserHandler {
	return &UserHandler{app: svc}
}

func (e *UserHandler) UserRoutes(r *gin.RouterGroup) {
	context := r.Group("/user")
	{
		context.PUT("", e.createUser)
		context.GET("", e.getAllUsers)
		context.POST("/login", e.login)
		context.POST("/easyLogin", e.easyLogin)
		context.PUT("/createEasyLogin", e.createEasyLoginUser)
		context.GET("/:id", e.getUser)
		context.PATCH("/:id", e.updateUser)
		context.DELETE("/:id", e.deleteUser)
	}

}

func (e *UserHandler) createUser(g *gin.Context)          {}
func (e *UserHandler) getAllUsers(g *gin.Context)         {}
func (e *UserHandler) login(g *gin.Context)               {}
func (e *UserHandler) easyLogin(g *gin.Context)           {}
func (e *UserHandler) createEasyLoginUser(g *gin.Context) {}
func (e *UserHandler) getUser(g *gin.Context)             {}
func (e *UserHandler) updateUser(g *gin.Context)          {}
func (e *UserHandler) deleteUser(g *gin.Context)          {}
