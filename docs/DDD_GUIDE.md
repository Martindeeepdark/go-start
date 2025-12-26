# DDD 架构使用指南

go-start 现在支持两种架构模式：MVC 和 DDD（领域驱动设计）。

## 架构对比

### MVC 架构（默认）

适合：中小型项目、快速开发

```
internal/
├── dal/           # 数据访问层（GORM Gen）
├── repository/    # 数据访问封装
├── service/       # 业务逻辑层
├── controller/    # HTTP 处理层
└── routes/        # 路由注册
```

**特点：**
- ✅ 简单直观
- ✅ 分层清晰
- ✅ 快速开发
- ❌ 业务复杂时难以维护

### DDD 架构

适合：大型复杂项目、业务逻辑复杂

```
internal/
├── domain/                    # 领域层
│   ├── user/
│   │   ├── User.go           # 用户实体（Entity）
│   │   ├── repository.go     # 仓储接口
│   │   └── service.go        # 领域服务
│   └── ...
│
├── application/              # 应用层
│   ├── user/
│   │   └── service.go        # 应用服务（用例）
│   └── ...
│
├── infrastructure/          # 基础设施层
│   └── persistence/
│       ├── user/
│       │   └── UserRepositoryImpl.go  # 仓储实现
│       └── ...
│
└── interface/               # 接口层
    └── http/
        ├── user/
        │   └── UserController.go
        └── routes.go
```

**特点：**
- ✅ 业务逻辑集中在 Domain 层
- ✅ 技术细节隔离在 Infrastructure 层
- ✅ 用例清晰体现在 Application 层
- ✅ 适合复杂业务领域
- ❌ 学习曲线较陡
- ❌ 代码量相对较多

## 使用 DDD 架构

### 1. 生成代码

```bash
go-start gen db \
  --dsn="root:pass@tcp(localhost:3306)/testdb" \
  --tables=users \
  --arch=ddd \
  --output=./myproject
```

### 2. 生成的代码结构

```
myproject/internal/
├── domain/
│   └── user/
│       ├── User.go              # 用户实体（领域模型）
│       ├── repository.go        # 仓储接口
│       └── service.go           # 领域服务
│
├── application/
│   └── user/
│       └── service.go           # 应用服务（用例）
│
├── infrastructure/
│   └── persistence/
│       └── user/
│           └── UserRepositoryImpl.go  # 仓储实现
│
└── interface/
    └── http/
        ├── user/
        │   └── UserController.go
        └── routes.go
```

## 各层说明

### 1. Domain 层（领域层）

**职责：** 封装业务逻辑和规则

**包含：**
- **Entity（实体）**: 领域模型，包含业务行为
  ```go
  type User struct {
      ID       uint
      Username string
      Email    string
  }

  // 业务方法
  func (u *User) Validate() error {
      // 业务规则验证
  }

  func (u *User) ChangeEmail(newEmail string) error {
      // 状态变更逻辑
  }
  ```

- **Repository Interface（仓储接口）**: 持久化抽象
  ```go
  type UserRepository interface {
      Save(ctx, user) error
      FindByID(ctx, id) (*User, error)
      FindAll(ctx) ([]*User, error)
      Delete(ctx, id) error
  }
  ```

- **Domain Service（领域服务）**: 复杂业务逻辑
  ```go
  type UserService struct {
      // 处理跨聚合根的业务逻辑
  }
  ```

### 2. Application 层（应用层）

**职责：** 编排用例，协调领域对象

**包含：**
- **Application Service（应用服务）**: 用例实现
  ```go
  type UserService struct {
      userRepo domain.UserRepository
  }

  // 用例：创建用户
  func (s *UserService) CreateUser(ctx, cmd) error {
      // 1. 创建领域对象
      user := domain.NewUser(cmd.Username, cmd.Email)

      // 2. 业务验证
      if err := user.Validate(); err != nil {
          return err
      }

      // 3. 保存到仓储
      return s.userRepo.Save(ctx, user)
  }
  ```

**特点：**
- 薄服务层，主要做协调
- 业务逻辑在 Domain 层
- 不直接依赖基础设施

### 3. Infrastructure 层（基础设施层）

**职责：** 提供技术实现

**包含：**
- **Repository Implementation（仓储实现）**: 持久化实现
  ```go
  type UserRepositoryImpl struct {
      db *gorm.DB
  }

  func (r *UserRepositoryImpl) Save(ctx, user) error {
      // 领域模型 → 数据模型
      dataModel := r.toDataModel(user)

      // 保存到数据库
      return r.db.WithContext(ctx).Save(dataModel).Error
  }
  ```

**特点：**
- 实现Domain层定义的接口
- 处理技术细节（数据库、缓存等）
- 不包含业务逻辑

### 4. Interface 层（接口层）

**职责：** 对外接口，处理协议

**包含：**
- **HTTP Controller**: HTTP 处理
  ```go
  type UserController struct {
      userService *application.UserService
  }

  func (c *UserController) Create(ctx *gin.Context) {
      var cmd CreateUserCommand
      ctx.ShouldBindJSON(&cmd)

      // 调用应用服务
      if err := c.userService.CreateUser(ctx, &cmd); err != nil {
          response.Error(ctx, http.StatusInternalServerError, err.Error())
          return
      }

      response.Success(ctx, gin.H{"message": "创建成功"})
  }
  ```

**特点：**
- 薄控制器，只处理 HTTP
- 调用 Application 层
- 不包含业务逻辑

## DDD 核心概念

### 1. 聚合根（Aggregate Root）

聚合根是领域模型的入口点：

```go
// User 是聚合根
type User struct {
    ID       uint
    Username string
    Orders   []Order  // 包含的其他实体
}

// 通过聚合根访问内部实体
user.Orders[0].Cancel()
```

### 2. 值对象（Value Object）

通过属性值标识的对象，不可变：

```go
type Email struct {
    value string
}

func NewEmail(value string) (*Email, error) {
    if !isValidEmail(value) {
        return nil, fmt.Errorf("无效邮箱")
    }
    return &Email{value: value}, nil
}
```

### 3. 领域事件（Domain Event）

记录领域内发生的重要事件：

```go
type UserRegisteredEvent struct {
    UserID    uint
    Username  string
    Timestamp time.Time
}

// 触发事件
user.addEvent(UserRegisteredEvent{...})
```

### 4. 仓储（Repository）

持久化抽象，隐藏数据访问细节：

```go
// Domain 层定义接口
type UserRepository interface {
    Save(ctx, user) error
    FindByID(ctx, id) (*User, error)
}

// Infrastructure 层实现
type UserRepositoryImpl struct { ... }
```

## 使用示例

### 完整流程

```bash
# 1. 生成 DDD 架构代码
go-start gen db \
  --dsn="root:pass@tcp(localhost:3306)/testdb" \
  --tables=users \
  --arch=ddd

# 2. 查看生成的代码
tree myproject/internal

# 3. 实现业务逻辑
# 编辑 internal/domain/user/User.go
func (u *User) ChangeEmail(newEmail string) error {
    // 业务规则验证
    if !strings.Contains(newEmail, "@") {
        return fmt.Errorf("邮箱格式不正确")
    }

    // 执行变更
    u.Email = newEmail
    return nil
}

# 4. 实现应用服务
# 编辑 internal/application/user/service.go
func (s *UserService) UpdateUserEmail(ctx, cmd) error {
    // 1. 查询聚合根
    user, err := s.userRepo.FindByID(ctx, cmd.ID)
    if err != nil {
        return err
    }

    // 2. 执行业务操作
    if err := user.ChangeEmail(cmd.NewEmail); err != nil {
        return err
    }

    // 3. 保存
    return s.userRepo.Save(ctx, user)
}

# 5. 实现仓储
# 编辑 internal/infrastructure/persistence/user/UserRepositoryImpl.go
func (r *UserRepositoryImpl) Save(ctx, user) error {
    // 领域模型 → 数据模型
    dataModel := &UserDataModel{
        ID: user.ID,
        Username: user.Username,
        Email: user.Email,
    }

    return r.db.WithContext(ctx).Save(dataModel).Error
}
```

## DDD vs MVC 选择

### 选择 MVC 当：
- ✅ 项目规模较小
- ✅ 业务逻辑简单
- ✅ 追求快速开发
- ✅ 团队对 DDD 不熟悉

### 选择 DDD 当：
- ✅ 业务逻辑复杂
- ✅ 领域概念丰富
- ✅ 长期维护的项目
- ✅ 需要清晰的业务边界
- ✅ 团队熟悉 DDD

## 下一步

1. **学习 DDD 概念**
   - 聚合（Aggregate）
   - 值对象（Value Object）
   - 领域事件（Domain Event）
   - 限界上下文（Bounded Context）

2. **实现业务逻辑**
   - 在 Entity 中添加业务方法
   - 在 Domain Service 中处理复杂逻辑
   - 在 Application Service 中编排用例

3. **完善基础设施**
   - 实现 Repository
   - 添加数据模型转换
   - 处理事务

---

DDD 是一个强大的工具，但需要学习和实践。建议从简单项目开始，逐步掌握！
