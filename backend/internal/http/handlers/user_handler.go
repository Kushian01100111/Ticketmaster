package handlers

import (
	"net/http"

	"github.com/Kushian01100111/Tickermaster/internal/app/user"
	"github.com/Kushian01100111/Tickermaster/internal/http/dto"
	"github.com/Kushian01100111/Tickermaster/internal/http/middleware"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	app user.UserService
}

func NewUserHandler(svc user.UserService) *UserHandler {
	return &UserHandler{app: svc}
}

func (u *UserHandler) PublicRoutes(r *gin.RouterGroup) {
	context := r.Group("/user")
	{
		context.POST("", u.createUser)
	}
}

func (u *UserHandler) PrivateRoutes(r *gin.RouterGroup) {
	context := r.Group("/user")
	{
		context.GET("", middleware.RequireRole("admin"), u.getAllUsers)
		context.GET("/:id", middleware.RequireRole("admin"), u.getUser)
		context.PATCH("/:id", middleware.RequireRole("costumer", "admin"), u.updateUser)
		context.DELETE("/:id", middleware.RequireRole("costumer", "admin"), u.deleteUser)
		context.GET("/byEmail", middleware.RequireRole("costumer", "admin"), u.getUserByEmail)
	}
}

func (u *UserHandler) getAllUsers(g *gin.Context) {
	users, err := u.app.GetAllUsers(g.Request.Context())
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, dto.ToUserSliceResponse(users))
}

func (u *UserHandler) createUser(g *gin.Context) {
	var req *dto.UserRequest

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind body of request"})
		return
	}

	user, err := u.app.CreateUser(g.Request.Context(), user.UserParams{
		Email:      req.Email,
		Role:       req.Role,
		Password:   req.Password,
		AuthMethod: req.AuthMethod,
	})

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusCreated, dto.ToUserResponse(user))
}

func (u *UserHandler) getUser(g *gin.Context) {
	id := g.Param("id")

	user, err := u.app.GetUser(g.Request.Context(), id)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, dto.ToUserResponse(user))
}

func (u *UserHandler) updateUser(g *gin.Context) {
	var req *dto.UpdateUserRequest

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind body of request"})
		return
	}

	id := g.Param("id")
	user, err := u.app.UpdateUser(g.Request.Context(), id, user.UpdateUserParams{
		Role:             req.Role,
		Password:         req.Password,
		AuthMethods:      req.AuthMethods,
		FailedLoginCount: req.FailedLoginCount,
		LastFailedLogin:  req.LastFailedLogin,
		BookedEvents:     req.BookedEvents,
	})

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusAccepted, dto.ToUserResponse(user))
}

func (u *UserHandler) deleteUser(g *gin.Context) {
	id := g.Param("id")
	if err := u.app.DeleteUser(g.Request.Context(), id); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.Status(http.StatusNoContent)
}

func (u *UserHandler) getUserByEmail(g *gin.Context) {
	var req *dto.RequestCode

	if err := g.ShouldBindJSON(req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := u.app.GetByEmail(g.Request.Context(), req.Email)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, dto.ToUserResponse(user))
}
