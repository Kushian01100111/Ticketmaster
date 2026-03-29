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

func (u *UserHandler) UserRoutes(r *gin.RouterGroup) {
	context := r.Group("/user")
	{
		context.GET("", u.getAllUsers)
		context.PUT("", u.createUser)
		context.GET("/:id", u.getUser)
		context.PATCH("/:id", u.updateUser)
		context.DELETE("/:id", u.deleteUser)
		context.GET("/byEmail", u.getUserByEmail)
	}
}

func (u *UserHandler) getAllUsers(g *gin.Context) {
	ctx, cancel := generateCtx() // Cambiar a g.Request.Context
	defer cancel()

	users, err := u.app.GetAllUsers(ctx)
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

	ctx, cancel := generateCtx() // Cambiar a g.Request.Context
	defer cancel()

	user, err := u.app.CreateUser(user.UserParams{
		Email:      req.Email,
		Role:       req.Role,
		Password:   req.Password,
		AuthMethod: req.AuthMethod,
	}, ctx)

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusCreated, dto.ToUserResponse(user))
}

func (u *UserHandler) getUser(g *gin.Context) {
	id := g.Param("id")

	ctx, cancel := generateCtx() // Cambiar a g.Request.Context
	defer cancel()

	user, err := u.app.GetUser(id, ctx)
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

	ctx, cancel := generateCtx() // Cambiar a g.Request.Context
	defer cancel()

	user, err := u.app.UpdateUser(id, user.UpdateUserParams{
		Role:             req.Role,
		Password:         req.Password,
		AuthMethods:      req.AuthMethods,
		FailedLoginCount: req.FailedLoginCount,
		LastFailedLogin:  req.LastFailedLogin,
		BookedEvents:     req.BookedEvents,
	}, ctx)

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusAccepted, dto.ToUserResponse(user))
}

func (u *UserHandler) deleteUser(g *gin.Context) {
	id := g.Param("id")

	ctx, cancel := generateCtx() // Cambiar a g.Request.Context
	defer cancel()

	if err := u.app.DeleteUser(id, ctx); err != nil {
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
