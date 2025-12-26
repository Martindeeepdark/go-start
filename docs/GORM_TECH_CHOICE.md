# 为什么选择 GORM Gen 而不是 GORM CLI？

## 背景

你提到了 GORM 官方新推出的 CLI 工具（gorm.io/cli），它确实有一些很吸引人的特性，比如：
- ✅ SQL 注释生成类型安全的查询方法
- ✅ Generics-first 设计
- ✅ Plain Go code，无运行时依赖

经过仔细评估和测试，我决定暂时不使用 GORM CLI，而是继续使用 GORM Gen（gorm.io/gen）。原因如下：

## 技术对比

### GORM CLI (gorm.io/cli)

**优势：**
1. **SQL 注释生成方法** - 这是最吸引人的特性
   ```go
   type UserQuery interface {
       // GetByEmail 根据邮箱获取用户
       // SELECT * FROM users WHERE email = @email LIMIT 1
       GetByEmail(ctx context.Context, email string) (*User, error)
   }
   ```
2. **更现代的设计** - 基于 Generics
3. **无运行时依赖** - 纯 Go 代码

**劣势：**
1. **版本太新** - 当前版本 v0.2.4（2024年12月）
2. **不够成熟** - 测试时发现泛型解析有 bug
3. **功能不完整** - 文档缺失，社区资源少
4. **稳定性未知** - 没有大规模生产环境验证

### GORM Gen (gorm.io/gen)

**优势：**
1. **成熟稳定** - v0.3.27，经过大量生产环境验证
2. **功能完整** - 丰富的特性和完善的文档
3. **类型安全** - 生成的代码编译时检查
4. **性能优秀** - 无反射开销
5. **生态支持** - 社区活跃，问题容易解决

**劣势：**
1. **不支持 SQL 注释** - 这是与 GORM CLI 相比的主要缺陷
2. **API 略显复杂** - 需要学习 GORM Gen 的 API

## 测试结果

我在测试 GORM CLI 时遇到了以下问题：

```bash
$ gorm gen -i ./examples -o ./generated
Error: error render template got error: failed to format generated code
got error generated/models.go:26:27: missing ',' in type argument list
```

**问题原因：** GORM CLI v0.2.4 在解析嵌套泛型时存在 bug：

```go
// 生成的代码有问题
Articles field.Slice[gorm-cli-demo.Article]  // ❌ 缺少逗号

// 应该是
Articles field.Slice[gorm-cli-demo.Article]  // ✅ 或者使用完整的导入路径
```

这个 bug 说明 GORM CLI 还不够成熟，不适合用于生产环境。

## 我的方案

基于以上分析，我采用了以下方案：

### 当前方案：GORM Gen + 自定义模板

```
GORM Gen (生成基础代码)
    ↓
Repository 模板 (封装 GORM Gen API)
    ↓
Service 模板 (业务逻辑 + 缓存)
    ↓
Controller 模板 (RESTful API)
```

**优势：**
1. ✅ 稳定可靠 - GORM Gen 经过生产验证
2. ✅ 现在可用 - 不需要等待 GORM CLI 成熟
3. ✅ 灵活可控 - 可以自定义模板满足需求
4. ✅ 渐进增强 - 后续可以平滑切换到 GORM CLI

### 未来方案：GORM CLI（待成熟后）

GORM CLI 达到以下条件后，可以考虑迁移：
1. 版本 >= v1.0.0
2. 修复了泛型解析的 bug
3. 有完整的生产环境案例
4. 文档完善，社区活跃

## 关键差异对比

| 特性 | GORM Gen | GORM CLI |
|------|----------|----------|
| 稳定性 | ✅ 生产验证 | ⚠️ 实验阶段 |
| SQL 注释 | ❌ 不支持 | ✅ 支持 |
| 泛型支持 | ✅ 完整 | ⚠️ 有 bug |
| 文档 | ✅ 完善 | ⚠️ 缺失 |
| 社区 | ✅ 活跃 | ⚠️ 新项目 |
| 类型安全 | ✅ 编译时检查 | ✅ 编译时检查 |
| 性能 | ✅ 无反射 | ✅ 无反射 |

## 实现细节

我们当前的实现结合了两者的优点：

### 1. 使用 GORM Gen 生成基础代码

```bash
go-start gen db --dsn="..." --tables=users
```

生成：
```
internal/dal/
├── query/
│   ├── gen.go           # GORM Gen 主入口
│   ├── users.go         # User 表的查询 API
│   └── ...
└── models.go            # 数据模型
```

### 2. 使用模板生成 Repository

基于 GORM Gen 的 Query API，生成更友好的 Repository：

```go
// Repository 封装 GORM Gen
type UserRepository struct {
    q *query.Query
}

// GetByID 使用 GORM Gen 的字段助手
func (r *UserRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
    return r.q.User.WithContext(ctx).
        Where(r.q.User.ID.Eq(id)).  // ✅ 类型安全，无魔法字符串
        First()
}

// GetByEmail 基于索引的高效查询
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
    return r.q.User.WithContext(ctx).
        Where(r.q.User.Email.Eq(email)).
        First()
}
```

### 3. 自动生成索引查询方法

如果表有索引 `email`，自动生成：
```go
// ByEmail 根据 email 查询（使用索引）
func (r *UserRepository) ByEmail(ctx context.Context, email interface{}) (*model.User, error) {
    return r.q.User.WithContext(ctx).
        Where(r.q.User.Email.Eq(email)).
        First()
}
```

## 总结

**为什么不用 GORM CLI？**
- 不够成熟，有 bug
- 没有生产环境验证
- 文档不完善

**为什么用 GORM Gen？**
- 成熟稳定，生产验证
- 功能完整，文档完善
- 社区活跃，问题好解决

**我们的优势：**
- 结合了 GORM Gen 的稳定性和模板的灵活性
- 自动生成基于索引的查询方法
- 详细的中文注释
- 为新人工程师提供最佳实践

**未来计划：**
GORM CLI 成熟后（v1.0+），可以考虑平滑迁移，享受 SQL 注释的优势。

---

**最终目标：** 帮助新人工程师专注业务逻辑，无需手写重复的 CRUD 代码。

无论用 GORM Gen 还是 GORM CLI，都能实现这个目标。现在 GORM Gen 是更安全的选择。
