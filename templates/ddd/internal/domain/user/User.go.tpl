package user

import (
	"errors"
	"time"
)

// User 用户实体（聚合根）
//
// 职责说明：
//   - 封装用户的业务规则和行为
//   - 确保用户对象的一致性和完整性
//   - 不包含任何基础设施细节
type User struct {
	ID        uint
	Username string
	Email    string
	Password string
	Age      int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// 错误定义
var (
	ErrInvalidUsername = errors.New("用户名至少3个字符")
	ErrInvalidEmail    = errors.New("邮箱格式不正确")
	ErrInvalidAge      = errors.New("年龄必须大于0")
)

// NewUser 创建用户聚合根
func NewUser(username, email, password string, age int) (*User, error) {
	user := &User{
		Username:  username,
		Email:     email,
		Password:  password,
		Age:       age,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 业务规则验证
	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

// Validate 验证用户的业务规则
func (u *User) Validate() error {
	if len(u.Username) < 3 {
		return ErrInvalidUsername
	}

	if !u.isValidEmail() {
		return ErrInvalidEmail
	}

	if u.Age <= 0 {
		return ErrInvalidAge
	}

	return nil
}

// ChangeEmail 修改邮箱
//
// 说明：
//   - 封装业务规则
//   - 确保状态变更的一致性
func (u *User) ChangeEmail(newEmail string) error {
	if !u.isValidEmailFormat(newEmail) {
		return ErrInvalidEmail
	}

	u.Email = newEmail
	u.UpdatedAt = time.Now()

	return nil
}

// ChangePassword 修改密码
func (u *User) ChangePassword(newPassword string) error {
	if len(newPassword) < 6 {
		return errors.New("密码至少6个字符")
	}

	u.Password = newPassword
	u.UpdatedAt = time.Now()

	return nil
}

// 辅助方法

func (u *User) isValidEmail() bool {
	return u.isValidEmailFormat(u.Email)
}

func (u *User) isValidEmailFormat(email string) bool {
	// TODO: 实现邮箱格式验证
	return true
}
