package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"{{.Module}}/internal/model"
	"{{.Module}}/internal/repository"
	"{{.Module}}/pkg/cache"
)

// 定义业务错误
var (
	ErrUserNotFound = errors.New("用户不存在")
	ErrUserExists   = errors.New("用户已存在")
)

// UserService 用户服务层
//
// 职责说明：
//   - 实现用户相关的业务逻辑
//   - 协调 repository 层和 cache 层
//   - 处理数据的缓存策略
//   - 实现业务的校验和规则
//
// 使用示例：
//   service := service.NewUserService(userRepo, cacheClient)
//   user, err := service.GetByID(ctx, 1)
type UserService struct {
	repo  repository.UserRepository // 用户数据访问层接口
	cache *cache.Cache              // Redis 缓存客户端
}

// NewUserService 创建用户服务实例
//
// 参数：
//   - repo: 用户数据访问层接口，负责数据库操作
//   - cache: Redis 缓存客户端，负责缓存操作
//
// 返回：
//   - *UserService: 用户服务实例
func NewUserService(repo repository.UserRepository, cache *cache.Cache) *UserService {
	return &UserService{
		repo:  repo,
		cache: cache,
	}
}

// Create 创建新用户
//
// 业务逻辑：
//   1. 检查用户名是否已存在（唯一性校验）
//   2. 如果存在则返回错误
//   3. 不存在则创建新用户
//
// 参数：
//   - ctx: 上下文，用于超时控制和取消操作
//   - user: 要创建的用户对象
//
// 返回：
//   - error: 创建失败时返回错误（如用户已存在）
//
// 使用场景：
//   - 用户注册
//   - 管理员添加用户
//
// 错误处理：
//   - ErrUserExists: 用户名已被占用
func (s *UserService) Create(ctx context.Context, user *model.User) error {
	// 1. 检查用户名是否已存在
	existing, err := s.repo.GetByUsername(ctx, user.Username)
	if err == nil && existing != nil {
		// 用户已存在，返回错误
		return ErrUserExists
	}

	// 2. 创建用户
	if err := s.repo.Create(ctx, user); err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}

	return nil
}

// GetByID 根据 ID 获取用户信息
//
// 缓存策略：
//   1. 先从 Redis 缓存查询
//   2. 缓存命中则直接返回
//   3. 缓存未命中则查询数据库
//   4. 查询成功后写入缓存，过期时间 10 分钟
//
// 参数：
//   - ctx: 上下文
//   - id: 用户 ID
//
// 返回：
//   - *model.User: 用户对象
//   - error: 用户不存在或查询失败时返回错误
//
// 使用场景：
//   - 用户详情页
//   - 用户资料编辑
//   - 验证用户身份
func (s *UserService) GetByID(ctx context.Context, id uint) (*model.User, error) {
	// 1. 构建缓存 key
	cacheKey := fmt.Sprintf("user:%d", id)

	{{if .WithRedis}}
	// 2. 尝试从缓存获取
	cached, err := s.cache.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		// 缓存命中，反序列化并返回
		var user model.User
		if err := cache.Unmarshal(cached, &user); err == nil {
			return &user, nil
		}
	}
	{{end}}

	// 3. 缓存未命中，从数据库查询
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	{{if .WithRedis}}
	// 4. 将用户信息写入缓存，过期时间 10 分钟
	if data, err := cache.Marshal(user); err == nil {
		_ = s.cache.Set(ctx, cacheKey, data, 10*time.Minute)
	}
	{{end}}

	return user, nil
}

// GetByUsername 根据用户名获取用户信息
//
// 参数：
//   - ctx: 上下文
//   - username: 用户名
//
// 返回：
//   - *model.User: 用户对象
//   - error: 用户不存在或查询失败时返回错误
//
// 使用场景：
//   - 用户登录
//   - 检查用户名可用性
func (s *UserService) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// Update 更新用户信息
//
// 业务逻辑：
//   1. 更新数据库中的用户信息
//   2. 清除相关缓存，保证数据一致性
//
// 参数：
//   - ctx: 上下文
//   - user: 要更新的用户对象（需包含 ID）
//
// 返回：
//   - error: 更新失败时返回错误
//
// 使用场景：
//   - 修改用户资料
//   - 更新用户状态
//
// 注意事项：
//   - 更新后会清除缓存，下次查询会重新从数据库读取
func (s *UserService) Update(ctx context.Context, user *model.User) error {
	// 1. 更新数据库
	if err := s.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("更新用户失败: %w", err)
	}

	{{if .WithRedis}}
	// 2. 清除缓存，保证数据一致性
	cacheKey := fmt.Sprintf("user:%d", user.ID)
	_ = s.cache.Del(ctx, cacheKey)
	{{end}}

	return nil
}

// Delete 删除用户
//
// 业务逻辑：
//   1. 从数据库删除用户（软删除）
//   2. 清除相关缓存
//
// 参数：
//   - ctx: 上下文
//   - id: 要删除的用户 ID
//
// 返回：
//   - error: 删除失败时返回错误
//
// 使用场景：
//   - 管理员删除用户
//   - 用户注销账号
//
// 注意事项：
//   - 使用软删除，数据不会物理删除
//   - 删除后会清除缓存
func (s *UserService) Delete(ctx context.Context, id uint) error {
	// 1. 删除用户（软删除）
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}

	{{if .WithRedis}}
	// 2. 清除缓存
	cacheKey := fmt.Sprintf("user:%d", id)
	_ = s.cache.Del(ctx, cacheKey)
	{{end}}

	return nil
}

// List 获取用户列表（分页）
//
// 业务逻辑：
//   1. 参数校验（页码、每页数量）
//   2. 从数据库查询用户列表
//   3. 返回用户列表和总数
//
// 参数：
//   - ctx: 上下文
//   - page: 页码（从 1 开始）
//   - pageSize: 每页数量（最大 100）
//
// 返回：
//   - []*model.User: 用户列表
//   - int64: 用户总数
//   - error: 查询失败时返回错误
//
// 使用场景：
//   - 用户管理页面
//   - 用户列表查询
//
// 分页说明：
//   - page < 1 时，默认为第 1 页
//   - pageSize <= 0 或 > 100 时，默认为 20 条/页
func (s *UserService) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	// 1. 参数校验和默认值设置
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	// 2. 查询用户列表
	users, total, err := s.repo.List(ctx, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("查询用户列表失败: %w", err)
	}

	return users, total, nil
}
