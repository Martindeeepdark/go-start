package spec

// 内置模板内容

const modelTemplate = `package model

import (
	"time"
	"gorm.io/gorm"
)

// {{.Model.Name}} {{.Model.Comment}}
type {{.Model.Name}} struct {
	{{- range $field := .Model.Fields}}
    {{$field.Name | ToCamelCase}} {{getGoType2 $field}} ` + "`" + `gorm:"{{getGormTag $field}}{{getIndexTags $field.Name $.Model.Indexes}}" json:"{{getJSONTag $field.Name $field.JSON}}"` + "`" + ` // {{$field.Comment}}
	{{- end}}
}

// TableName specifies the table name for {{.Model.Name}} model
func ({{.Model.Name}}) TableName() string {
	return "{{.Model.Table}}"
}
`

const repositoryTemplate = `package repository

import (
    "context"
    "gorm.io/gorm"
    "{{.Spec.Project.Module}}/internal/model"
)

// {{.Model.Name}}Repository {{.Model.Name}}数据访问层
//
// 职责说明：
//   - 负责 {{.Model.Name}} 的数据库操作
//   - 提供基础的 CRUD 方法
//   - 支持复杂查询和事务处理
type {{.Model.Name}}Repository struct {
	db *gorm.DB
}

// New{{.Model.Name}}Repository 创建 {{.Model.Name}} 仓储实例
func New{{.Model.Name}}Repository(db *gorm.DB) *{{.Model.Name}}Repository {
	return &{{.Model.Name}}Repository{db: db}
}

// Create 创建 {{.Model.Name}}
func (r *{{.Model.Name}}Repository) Create(ctx context.Context, {{.Model.Name | ToLowerCamelCase}} *model.{{.Model.Name}}) error {
	return r.db.WithContext(ctx).Create({{.Model.Name | ToLowerCamelCase}}).Error
}

// GetByID 根据 ID 获取 {{.Model.Name}}
func (r *{{.Model.Name}}Repository) GetByID(ctx context.Context, id uint) (*model.{{.Model.Name}}, error) {
	var {{.Model.Name | ToLowerCamelCase}} model.{{.Model.Name}}
	err := r.db.WithContext(ctx).First(&{{.Model.Name | ToLowerCamelCase}}, id).Error
	if err != nil {
		return nil, err
	}
	return &{{.Model.Name | ToLowerCamelCase}}, nil
}

// Update 更新 {{.Model.Name}}
func (r *{{.Model.Name}}Repository) Update(ctx context.Context, {{.Model.Name | ToLowerCamelCase}} *model.{{.Model.Name}}) error {
	return r.db.WithContext(ctx).Save({{.Model.Name | ToLowerCamelCase}}).Error
}

// Delete 删除 {{.Model.Name}}
func (r *{{.Model.Name}}Repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.{{.Model.Name}}{}, id).Error
}

// List 获取 {{.Model.Name}} 列表（分页）
func (r *{{.Model.Name}}Repository) List(ctx context.Context, page, pageSize int) ([]*model.{{.Model.Name}}, int64, error) {
	var {{.Model.Name | ToLowerCamelCase}}s []*model.{{.Model.Name}}
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&model.{{.Model.Name}}{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&{{.Model.Name | ToLowerCamelCase}}s).Error; err != nil {
		return nil, 0, err
	}

	return {{.Model.Name | ToLowerCamelCase}}s, total, nil
}
`

const serviceTemplate = `package service

import (
    "context"
    "errors"
    "fmt"
    "time"

    "{{.Spec.Project.Module}}/internal/model"
    "{{.Spec.Project.Module}}/internal/repository"
    {{- if or .GetCacheEnabled .ListCacheEnabled}}
    "github.com/Martindeeepdark/go-start/pkg/commonadapter"
    {{- end}}
)

// 定义业务错误
var (
	Err{{.Model.Name}}NotFound = errors.New("{{.Model.Name}}不存在")
)

// {{.Model.Name}}Service {{.Model.Name}}服务层
//
// 职责说明：
//   - 实现 {{.Model.Name}} 相关的业务逻辑
//   - 协调 repository 层和 cache 层
//   - 处理数据的缓存策略
//   - 实现业务的校验和规则
type {{.Model.Name}}Service struct {
    repo  *repository.{{.Model.Name}}Repository
}
// listCacheEntry 用于列表缓存封装
type listCacheEntry struct {
    List []*model.{{.Model.Name}}
    Total int64
}

// New{{.Model.Name}}Service 创建 {{.Model.Name}} 服务实例
func New{{.Model.Name}}Service(repo *repository.{{.Model.Name}}Repository) *{{.Model.Name}}Service {
    return &{{.Model.Name}}Service{
        repo:  repo,
    }
}

// Create 创建 {{.Model.Name}}
func (s *{{.Model.Name}}Service) Create(ctx context.Context, {{.Model.Name | ToLowerCamelCase}} *model.{{.Model.Name}}) error {
    if err := s.repo.Create(ctx, {{.Model.Name | ToLowerCamelCase}}); err != nil {
        return fmt.Errorf("创建{{.Model.Name}}失败: %w", err)
    }
    _, cache, _, _, _, _ := commonadapter.Abilities()
    _ = cache.Delete(fmt.Sprintf("{{.Model.Name | ToLowerCamelCase}}:%d", {{.Model.Name | ToLowerCamelCase}}.ID))
    _ = cache.DeleteByPattern("{{.Model.Name | ToLowerCamelCase}}:list:*")
    return nil
}

// GetByID 根据 ID 获取 {{.Model.Name}}
func (s *{{.Model.Name}}Service) GetByID(ctx context.Context, id uint) (*model.{{.Model.Name}}, error) {
    {{- if .GetCacheEnabled}}
    cacheKey := fmt.Sprintf("{{.Model.Name | ToLowerCamelCase}}:%d", id)
    _, cache, _, _, _, _ := commonadapter.Abilities()
    if v, err := cache.Get(cacheKey); err == nil {
        if {{.Model.Name | ToLowerCamelCase}}, ok := v.(*model.{{.Model.Name}}); ok {
            return {{.Model.Name | ToLowerCamelCase}}, nil
        }
    }
    {{- end}}
    {{.Model.Name | ToLowerCamelCase}}, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, Err{{.Model.Name}}NotFound
    }
    {{- if .GetCacheEnabled}}
    _, cache, _, _, _, _ := commonadapter.Abilities()
    _ = cache.Set(cacheKey, {{.Model.Name | ToLowerCamelCase}}, {{if .GetCacheTTL}}{{.GetCacheTTL}}{{else}}600{{end}})
    {{- end}}
    return {{.Model.Name | ToLowerCamelCase}}, nil
}

// Update 更新 {{.Model.Name}}
func (s *{{.Model.Name}}Service) Update(ctx context.Context, {{.Model.Name | ToLowerCamelCase}} *model.{{.Model.Name}}) error {
    if err := s.repo.Update(ctx, {{.Model.Name | ToLowerCamelCase}}); err != nil {
        return fmt.Errorf("更新{{.Model.Name}}失败: %w", err)
    }
    _, cache, _, _, _, _ := commonadapter.Abilities()
    _ = cache.Delete(fmt.Sprintf("{{.Model.Name | ToLowerCamelCase}}:%d", {{.Model.Name | ToLowerCamelCase}}.ID))
    _ = cache.DeleteByPattern("{{.Model.Name | ToLowerCamelCase}}:list:*")
    return nil
}

// Delete 删除 {{.Model.Name}}
func (s *{{.Model.Name}}Service) Delete(ctx context.Context, id uint) error {
    if err := s.repo.Delete(ctx, id); err != nil {
        return fmt.Errorf("删除{{.Model.Name}}失败: %w", err)
    }
    _, cache, _, _, _, _ := commonadapter.Abilities()
    _ = cache.Delete(fmt.Sprintf("{{.Model.Name | ToLowerCamelCase}}:%d", id))
    _ = cache.DeleteByPattern("{{.Model.Name | ToLowerCamelCase}}:list:*")
    return nil
}

// List 获取 {{.Model.Name}} 列表（分页）
func (s *{{.Model.Name}}Service) List(ctx context.Context, page, pageSize int) ([]*model.{{.Model.Name}}, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

    {{- if .ListCacheEnabled}}
    cacheKey := fmt.Sprintf("{{.Model.Name | ToLowerCamelCase}}:list:%d:%d", page, pageSize)
    _, cache, _, _, _, _ := commonadapter.Abilities()
    if v, err := cache.Get(cacheKey); err == nil {
        if entry, ok := v.(*listCacheEntry); ok {
            return entry.List, entry.Total, nil
        }
    }
    {{- end}}
    res, total, err := s.repo.List(ctx, page, pageSize)
    if err != nil {
        return nil, 0, err
    }
    {{- if .ListCacheEnabled}}
    _, cache, _, _, _, _ := commonadapter.Abilities()
    _ = cache.Set(cacheKey, &listCacheEntry{List: res, Total: total}, {{if .ListCacheTTL}}{{.ListCacheTTL}}{{else}}300{{end}})
    {{- end}}
    return res, total, nil
}
`

const controllerTemplate = `package controller

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "{{.Spec.Project.Module}}/internal/model"
    "{{.Spec.Project.Module}}/internal/service"
    "{{.Spec.Project.Module}}/pkg/httpx/response"
    "github.com/Martindeeepdark/go-start/pkg/commonadapter"
    {{- if or .CreateValidator .UpdateValidator}}
    "{{.Spec.Project.Module}}/internal/validator"
    {{- end}}
)

// {{.Model.Name}}Controller {{.Model.Name}}控制器
//
// 职责说明：
//   - 处理 {{.Model.Name}} 相关的 HTTP 请求
//   - 参数验证和绑定
//   - 调用服务层处理业务逻辑
//   - 统一返回响应格式
type {{.Model.Name}}Controller struct {
	service *service.{{.Model.Name}}Service
}

// New{{.Model.Name}}Controller 创建 {{.Model.Name}} 控制器实例
func New{{.Model.Name}}Controller(service *service.{{.Model.Name}}Service) *{{.Model.Name}}Controller {
	return &{{.Model.Name}}Controller{
		service: service,
	}
}

// Create 创建 {{.Model.Name}}
// @Summary 创建{{.Model.Name}}
// @Description 创建一个新的{{.Model.Name}}
// @Tags {{.Model.Name}}
// @Accept json
// @Produce json
// @Param {{.Model.Name | ToLowerCamelCase}} body model.{{.Model.Name}} true "{{.Model.Name}}信息"
// @Success 200 {object} response.Response
// @Router /api/v1/{{.Model.Name | ToLowerCamelCase}}s [post]
func (c *{{.Model.Name}}Controller) Create(ctx *gin.Context) {
    _, _, audit, idemp, _, _ := commonadapter.Abilities()
    var userID string
    if v, ok := ctx.Get("UserID"); ok {
        if s, ok2 := v.(string); ok2 { userID = s }
    }
    ik := ctx.GetHeader("Idempotency-Key")
    if ik != "" {
        ok, err := idemp.CheckAndSet(ik, 600)
        if err != nil { response.Error(ctx, http.StatusInternalServerError, "幂等校验失败"); return }
        if !ok { response.Error(ctx, http.StatusConflict, "重复请求"); return }
    }
    {{- if .CreateValidator}}
    var req validator.{{.CreateValidator}}
    if err := ctx.ShouldBindJSON(&req); err != nil {
        response.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
        return
    }
    if err := req.Validate(); err != nil {
        response.Error(ctx, http.StatusBadRequest, err.Error())
        return
    }
    {{- end}}

    var {{.Model.Name | ToLowerCamelCase}} model.{{.Model.Name}}
    if err := ctx.ShouldBindJSON(&{{.Model.Name | ToLowerCamelCase}}); err != nil {
        response.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
        return
    }

    if err := c.service.Create(ctx, &{{.Model.Name | ToLowerCamelCase}}); err != nil {
        response.Error(ctx, http.StatusInternalServerError, err.Error())
        return
    }
    _ = audit.Record(userID, "{{.Model.Name}}", "create", "success", "")
    response.Success(ctx, gin.H{"id": {{.Model.Name | ToLowerCamelCase}}.ID})
}

// GetByID 获取 {{.Model.Name}} 详情
func (c *{{.Model.Name}}Controller) GetByID(ctx *gin.Context) {
    {{- if .GetAuth}}
    var userID string
    if v, ok := ctx.Get("UserID"); ok {
        if s, ok2 := v.(string); ok2 { userID = s }
    }
    {{- end}}
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	{{.Model.Name | ToLowerCamelCase}}, err := c.service.GetByID(ctx, uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}

	response.Success(ctx, {{.Model.Name | ToLowerCamelCase}})
}

// Update 更新 {{.Model.Name}}
func (c *{{.Model.Name}}Controller) Update(ctx *gin.Context) {
    _, _, audit, _, _, _ := commonadapter.Abilities()
    var userID string
    if v, ok := ctx.Get("UserID"); ok {
        if s, ok2 := v.(string); ok2 { userID = s }
    }
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

    {{- if .UpdateValidator}}
    var req validator.{{.UpdateValidator}}
    if err := ctx.ShouldBindJSON(&req); err != nil {
        response.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
        return
    }
    if err := req.Validate(); err != nil {
        response.Error(ctx, http.StatusBadRequest, err.Error())
        return
    }
    {{- end}}

    var {{.Model.Name | ToLowerCamelCase}} model.{{.Model.Name}}
    if err := ctx.ShouldBindJSON(&{{.Model.Name | ToLowerCamelCase}}); err != nil {
        response.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
        return
    }

	{{.Model.Name | ToLowerCamelCase}}.ID = uint(id)
    if err := c.service.Update(ctx, &{{.Model.Name | ToLowerCamelCase}}); err != nil {
        response.Error(ctx, http.StatusInternalServerError, err.Error())
        return
    }
    _ = audit.Record(userID, "{{.Model.Name}}", "update", "success", "")
    response.Success(ctx, nil)
}

// Delete 删除 {{.Model.Name}}
func (c *{{.Model.Name}}Controller) Delete(ctx *gin.Context) {
    _, _, audit, _, _, _ := commonadapter.Abilities()
    var userID string
    if v, ok := ctx.Get("UserID"); ok {
        if s, ok2 := v.(string); ok2 { userID = s }
    }
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

    if err := c.service.Delete(ctx, uint(id)); err != nil {
        response.Error(ctx, http.StatusInternalServerError, err.Error())
        return
    }
    _ = audit.Record(userID, "{{.Model.Name}}", "delete", "success", "")
    response.Success(ctx, nil)
}

// List 获取 {{.Model.Name}} 列表
func (c *{{.Model.Name}}Controller) List(ctx *gin.Context) {
    {{- if .ListAuth}}
    token := ctx.GetHeader("Authorization")
    if len(token) > 7 && (token[:7] == "Bearer " || token[:7] == "bearer ") { token = token[7:] }
    auth, _, _, _, _, _ := commonadapter.Abilities()
    userID, err := auth.VerifyToken(token)
    if err != nil { response.Error(ctx, http.StatusUnauthorized, "未授权"); return }
    {{- if .ListPerm}}
    if err := auth.RequirePermission(userID, "{{.ListPerm}}"); err != nil { response.Error(ctx, http.StatusForbidden, "权限不足"); return }
    {{- end}}
    {{- end}}
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	{{.Model.Name | ToLowerCamelCase}}s, total, err := c.service.List(ctx, page, pageSize)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"list":  {{.Model.Name | ToLowerCamelCase}}s,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}
`

const routesTemplate = `package routes

import (
    "github.com/gin-gonic/gin"
    "{{.Spec.Project.Module}}/internal/controller"
    "github.com/Martindeeepdark/go-start/pkg/httpx/middleware"
)

// RegisterAutoRoutes 自动注册所有路由
//
// 此文件由 spec 工具自动生成，请勿手动修改
func RegisterAutoRoutes(r *gin.Engine, controllers *controller.Controllers) {
    v1 := r.Group("/api/v1")
    {
        {{- range $info := .ModelsInfo}}
        {{$info.ModelVar}} := v1.Group("/{{pluralize $info.ModelVar}}")
        {
            // Create
            {{$info.ModelVar}}.POST("",
                {{- if or $info.CreateAuth $info.CreatePerm}}
                {{- if $info.CreateAuth}}middleware.RequireAuth(),{{end}}
                {{- if $info.CreatePerm}}middleware.RequirePermission("{{$info.CreatePerm}}"),{{end}}
                {{- end}}
                controllers.{{$info.ModelName}}.Create,
            )
            // List
            {{$info.ModelVar}}.GET("",
                {{- if or $info.ListAuth $info.ListPerm}}
                {{- if $info.ListAuth}}middleware.RequireAuth(),{{end}}
                {{- if $info.ListPerm}}middleware.RequirePermission("{{$info.ListPerm}}"),{{end}}
                {{- end}}
                controllers.{{$info.ModelName}}.List,
            )
            // GetByID
            {{$info.ModelVar}}.GET("/:id",
                {{- if or $info.GetAuth $info.GetPerm}}
                {{- if $info.GetAuth}}middleware.RequireAuth(),{{end}}
                {{- if $info.GetPerm}}middleware.RequirePermission("{{$info.GetPerm}}"),{{end}}
                {{- end}}
                controllers.{{$info.ModelName}}.GetByID,
            )
            // Update
            {{$info.ModelVar}}.PUT("/:id",
                {{- if or $info.UpdateAuth $info.UpdatePerm}}
                {{- if $info.UpdateAuth}}middleware.RequireAuth(),{{end}}
                {{- if $info.UpdatePerm}}middleware.RequirePermission("{{$info.UpdatePerm}}"),{{end}}
                {{- end}}
                controllers.{{$info.ModelName}}.Update,
            )
            // Delete
            {{$info.ModelVar}}.DELETE("/:id",
                {{- if or $info.DeleteAuth $info.DeletePerm}}
                {{- if $info.DeleteAuth}}middleware.RequireAuth(),{{end}}
                {{- if $info.DeletePerm}}middleware.RequirePermission("{{$info.DeletePerm}}"),{{end}}
                {{- end}}
                controllers.{{$info.ModelName}}.Delete,
            )
        }
        {{- end}}
    }
}
`

const validatorTemplate = `package validator

import (
    "fmt"
    "github.com/go-playground/validator/v10"
)

// {{.Request.Name}} {{.Request.Comment}}
// 使用 validator.v10 进行字段规则校验
type {{.Request.Name}} struct {
    {{- range $f := .Request.Fields}}
    {{ToCamelCase $f.Name}} {{getReqFieldType $f.Rules}} ` + "`" + `validate:"{{$f.Rules}}"` + "`" + ` // {{$f.Comment}}
    {{- end}}
}

// Validate 对 {{.Request.Name}} 进行规则校验
// 返回错误以指示具体的规则失败信息
func (r *{{.Request.Name}}) Validate() error {
    v := validator.New()
    if err := v.Struct(r); err != nil {
        return fmt.Errorf("{{.Request.Name}} 校验失败: %w", err)
    }
    return nil
}
`

// getBuiltinTemplate 获取内置模板内容
func getBuiltinTemplate(name string) string {
	templates := map[string]string{
		"model.go.tmpl":      modelTemplate,
		"repository.go.tmpl": repositoryTemplate,
		"service.go.tmpl":    serviceTemplate,
		"controller.go.tmpl": controllerTemplate,
		"routes.go.tmpl":     routesTemplate,
		"validator.go.tmpl":  validatorTemplate,
	}

	return templates[name]
}
