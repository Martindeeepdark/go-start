package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"{{.Module}}/internal/application/user"
	"{{.Module}}/pkg/httpx/response"
)

// UserController 用户控制器
type UserController struct {
	userService *user.UserService
}

// NewUserController 创建用户控制器
func NewUserController(userService *user.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// CreateUser 创建用户
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Age      int    `json:"age" binding:"required,min=1"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := c.userService.CreateUser(ctx, req.Username, req.Email, req.Password, req.Age); err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, gin.H{"message": "创建成功"})
}

// UpdateEmail 更新邮箱
func (c *UserController) UpdateEmail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := c.userService.UpdateUserEmail(ctx, uint(id), req.Email); err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, gin.H{"message": "更新成功"})
}
