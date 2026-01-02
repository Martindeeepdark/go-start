package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"{{.Module}}/internal/model"
	"{{.Module}}/internal/service"
	"{{.Module}}/pkg/httpx/response"
)

// UserController represents the user controller
type UserController struct {
	service *service.Service
}

// NewUserController creates a new user controller instance
func NewUserController(service *service.Service) *UserController {
	return &UserController{service: service}
}

// CreateRequest represents the create user request
type CreateRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UpdateRequest represents the update user request
type UpdateRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
	Status   *int   `json:"status" binding:"omitempty,oneof=0 1"`
}

// Create creates a new user
func (ctrl *UserController) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Status:   1,
	}
	// In production, hash the password here
	// user.Password = hashPassword(req.Password)

	if err := ctrl.service.User.Create(c.Request.Context(), user); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, user)
}

// Get gets a user by ID
func (ctrl *UserController) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	user, err := ctrl.service.User.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, user)
}

// Update updates a user
func (ctrl *UserController) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := ctrl.service.User.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Status != nil {
		user.Status = *req.Status
	}

	if err := ctrl.service.User.Update(c.Request.Context(), user); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, user)
}

// Delete deletes a user
func (ctrl *UserController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	if err := ctrl.service.User.Delete(c.Request.Context(), uint(id)); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// List lists users with pagination
func (ctrl *UserController) List(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	users, total, err := ctrl.service.User.List(c.Request.Context(), page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Paginated(c, users, total, page, pageSize)
}
