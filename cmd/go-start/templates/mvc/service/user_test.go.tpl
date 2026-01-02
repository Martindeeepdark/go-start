package service

import (
	"context"
	"errors"
	"testing"

	"{{.Module}}/internal/model"
	"{{.Module}}/internal/repository"
)

// ============================================
// Mock Repository（用于测试）
// ============================================

// MockUserRepository 模拟的 User Repository
// 实现了 repository.UserRepository 接口
type MockUserRepository struct {
	Users map[uint]*model.User
	// 可以添加测试辅助字段
	CreateCalled bool
	FindByIDCalled bool
}

// NewMockUserRepository 创建 mock repository
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		Users: make(map[uint]*model.User),
	}
}

// Create 创建用户（mock 实现）
func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	m.CreateCalled = true
	if user.ID == 0 {
		user.ID = uint(len(m.Users) + 1)
	}
	m.Users[user.ID] = user
	return nil
}

// FindByID 根据 ID 查找用户（mock 实现）
func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
	m.FindByIDCalled = true
	user, ok := m.Users[id]
	if !ok {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}

// FindByUsername 根据用户名查找（mock 实现）
func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	for _, user := range m.Users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("用户不存在")
}

// FindByEmail 根据邮箱查找（mock 实现）
func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	for _, user := range m.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("用户不存在")
}

// Update 更新用户（mock 实现）
func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	if _, ok := m.Users[user.ID]; !ok {
		return errors.New("用户不存在")
	}
	m.Users[user.ID] = user
	return nil
}

// Delete 删除用户（mock 实现）
func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	if _, ok := m.Users[id]; !ok {
		return errors.New("用户不存在")
	}
	delete(m.Users, id)
	return nil
}

// List 分页查询（mock 实现）
func (m *MockUserRepository) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	users := make([]*model.User, 0, len(m.Users))
	for _, user := range m.Users {
		users = append(users, user)
	}
	return users, int64(len(users)), nil
}

// ============================================
// 单元测试示例
// ============================================

// TestUserService_Create 测试创建用户
func TestUserService_Create(t *testing.T) {
	// 1. 创建 mock repository
	mockRepo := NewMockUserRepository()

	// 2. 创建 service（注入 mock，cache 为 nil）
	svc := NewUserService(mockRepo, nil)

	// 3. 准备测试数据
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
	}

	// 4. 执行测试
	ctx := context.Background()
	err := svc.Create(ctx, user)

	// 5. 验证结果
	if err != nil {
		t.Errorf("创建用户失败: %v", err)
	}

	if user.ID == 0 {
		t.Error("用户 ID 应该被设置")
	}

	if !mockRepo.CreateCalled {
		t.Error("Repository.Create 应该被调用")
	}

	// 验证用户被保存
	if mockRepo.Users[user.ID] == nil {
		t.Error("用户应该被保存到 repository")
	}
}

// TestUserService_GetByID 测试根据 ID 获取用户
func TestUserService_GetByID(t *testing.T) {
	// 1. 创建 mock repository
	mockRepo := NewMockUserRepository()

	// 2. 准备测试数据
	testUser := &model.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}
	mockRepo.Users[1] = testUser

	// 3. 创建 service
	svc := NewUserService(mockRepo, nil)

	// 4. 执行测试
	ctx := context.Background()
	user, err := svc.GetByID(ctx, 1)

	// 5. 验证结果
	if err != nil {
		t.Errorf("获取用户失败: %v", err)
	}

	if user == nil {
		t.Fatal("用户不应该为空")
	}

	if user.ID != 1 {
		t.Errorf("用户 ID 不匹配，期望 1，实际 %d", user.ID)
	}

	if user.Username != "testuser" {
		t.Errorf("用户名不匹配，期望 testuser，实际 %s", user.Username)
	}

	if !mockRepo.FindByIDCalled {
		t.Error("Repository.FindByID 应该被调用")
	}
}

// TestUserService_GetByID_NotFound 测试获取不存在的用户
func TestUserService_GetByID_NotFound(t *testing.T) {
	// 1. 创建空的 mock repository
	mockRepo := NewMockUserRepository()

	// 2. 创建 service
	svc := NewUserService(mockRepo, nil)

	// 3. 执行测试
	ctx := context.Background()
	user, err := svc.GetByID(ctx, 999)

	// 4. 验证结果
	if err == nil {
		t.Error("应该返回错误")
	}

	if err != ErrUserNotFound {
		t.Errorf("错误类型不匹配，期望 ErrUserNotFound，实际 %v", err)
	}

	if user != nil {
		t.Error("用户应该为空")
	}
}

// TestUserService_Update 测试更新用户
func TestUserService_Update(t *testing.T) {
	// 1. 创建 mock repository
	mockRepo := NewMockUserRepository()

	// 2. 准备测试数据
	testUser := &model.User{
		ID:       1,
		Username: "oldname",
		Email:    "old@example.com",
	}
	mockRepo.Users[1] = testUser

	// 3. 创建 service
	svc := NewUserService(mockRepo, nil)

	// 4. 修改用户
	testUser.Username = "newname"
	ctx := context.Background()
	err := svc.Update(ctx, testUser)

	// 5. 验证结果
	if err != nil {
		t.Errorf("更新用户失败: %v", err)
	}

	// 验证 repository 中的用户也被更新
	savedUser := mockRepo.Users[1]
	if savedUser.Username != "newname" {
		t.Errorf("用户名没有被更新，期望 newname，实际 %s", savedUser.Username)
	}
}

// TestUserService_Delete 测试删除用户
func TestUserService_Delete(t *testing.T) {
	// 1. 创建 mock repository
	mockRepo := NewMockUserRepository()

	// 2. 准备测试数据
	testUser := &model.User{
		ID:       1,
		Username: "testuser",
	}
	mockRepo.Users[1] = testUser

	// 3. 创建 service
	svc := NewUserService(mockRepo, nil)

	// 4. 执行删除
	ctx := context.Background()
	err := svc.Delete(ctx, 1)

	// 5. 验证结果
	if err != nil {
		t.Errorf("删除用户失败: %v", err)
	}

	// 验证用户被删除
	if _, ok := mockRepo.Users[1]; ok {
		t.Error("用户应该被从 repository 中删除")
	}
}

// ============================================
// 测试说明
// ============================================
//
// 如何运行测试：
//   go test ./internal/service
//
// 运行测试并查看覆盖率：
//   go test -cover ./internal/service
//
// 运行特定测试：
//   go test -run TestUserService_Create ./internal/service
//
// 查看详细输出：
//   go test -v ./internal/service
//
// 优势：
//   ✅ 不需要真实数据库
//   ✅ 测试运行速度快
//   ✅ 可以精确控制测试场景
//   ✅ 易于维护和调试
