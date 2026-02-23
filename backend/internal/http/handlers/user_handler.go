package handlers

import (
	"net/http"

	"github.com/Kushian01100111/Tickermaster/internal/app/user"
	"github.com/Kushian01100111/Tickermaster/internal/http/dto"
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
		context.GET("", e.getAllUsers)
		context.POST("/create", e.createUser)
		context.GET("/:userId", e.getUser)
		context.PATCH("/:userId", e.updateUser)
		context.DELETE("/:userid", e.deleteUser)
		context.POST("/login", e.login)
		context.POST("/requestToken", e.requestToken)
		context.POST("/authenticateToken", e.loginPasswordless)
		context.POST("/create/requestToken", e.signUpRequestToken)
		context.POST("/create/authenticateToken", e.signUpPasswordless)
	}

}

func (u *UserHandler) getAllUsers(g *gin.Context) {
	ctx, cancel := generateCtx()
	defer cancel()

	users, err := u.app.GetAllUser(ctx)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, dto.ToUserSliceResponse(users))
}
func (u *UserHandler) createUser(g *gin.Context) {}
func (u *UserHandler) getUser(g *gin.Context)    {}
func (u *UserHandler) updateUser(g *gin.Context) {}
func (u *UserHandler) deleteUser(g *gin.Context) {}
func (u *UserHandler) login(g *gin.Context)      {}

func (u *UserHandler) requestToken(g *gin.Context)       {}
func (u *UserHandler) loginPasswordless(g *gin.Context)  {}
func (u *UserHandler) signUpRequestToken(g *gin.Context) {}
func (u *UserHandler) signUpPasswordless(g *gin.Context) {}
