# go-start 设计深度反思

## 🔍 核心问题诊断

### 问题1：我们到底在解决什么问题？

**原始目标**：
> 帮助新人工程师快速上手，让高级工程师大展身手
> 让工程师专注业务逻辑，无需手写重复的 CRUD 代码

**现实**：
- ❌ 新人工程师：生成的代码跑不起来，不知道怎么用
- ❌ 高级工程师：代码不灵活，模板有 bug，无法定制
- ❌ CRUD 自动化：命令混乱，不知道什么时候用哪个

**根本问题**：**我们想做的事情太多了！**

试图同时解决：
1. 项目创建（像 `create-react-app`）
2. CRUD 代码生成（像 `Rails scaffold`）
3. 架构模式（MVC/DDD）
4. 规范驱动开发（Spec-Kit）

**结果**：每个都没做好！

### 问题2：命令设计混乱

```bash
go-start create my-api           # 创建项目
go-start gen db --tables=users    # 生成代码
```

**用户的困惑**：
- `create` 和 `gen` 是什么关系？
- `gen` 是生成什么？项目？代码？模型？
- 为什么不能 `create my-api --with-users-table`？

**更好的设计（参考 Rails）**：
```bash
# Rails 的设计
rails new myapi                    # 创建项目
rails generate model User          # 生成模型
rails generate scaffold User        # 生成脚手架
rails db:migrate                   # 迁移数据库
```

**我们应该**：
```bash
go-start new myapi                 # 创建项目
go-start generate model User        # 生成模型
go-start generate api User          # 生成完整 API
go-start db migrate                 # 迁移数据库
```

### 问题3：功能没做完就急着做下一个

**我们现在有**：
- ✅ `create` 命令（但模板有 bug）
- ✅ `gen db` 命令（生成代码）
- ✅ MVC 和 DDD 架构
- ✅ Spec-Kit

**但是**：
- ❌ `create` 命令生成的代码**跑不起来**
- ❌ 从来没有端到端测试过
- ❌ 模块路径是硬编码的假路径

**优先级错了！应该**：
1. 先让 `create` 完美工作
2. 确保生成的代码能运行
3. 再考虑 `gen` 命令

### 问题4：DDD 的过度设计

**DDD 生成的问题**：
- 生成了一堆空模板
- 没有实际生成领域逻辑
- 用户还是要手写大量代码

**真相**：
- DDD 是高级工程师自己用的
- 不应该自动生成 DDD 代码
- 应该只提供 DDD 模板，让工程师自己实现

### 问题5：缺少真实的使用场景

**我们有**：
- 详细的文档
- 完整的架构
- 漂亮的代码

**但是缺少**：
- ❌ 一个真实可运行的示例项目
- ❌ 从 0 到 1 的完整教程
- ❌ 常见问题的解决方案

## 💡 重新定位：我们应该做什么？

### 核心价值主张（重新定义）

**不要**：
- ❌ 什么都做（项目创建 + CRUD + MVC + DDD + Spec）
- ❌ 代替工程师思考
- ❌ 过度设计

**应该**：
- ✅ **专注一件事**：从数据库生成 CRUD 代码
- ✅ **做到极致**：生成的代码开箱即用
- ✅ **简单清晰**：命令直观，文档简明

### 新的产品定位

```
go-start = 从数据库生成 CRUD 代码的 CLI 工具

核心功能：
1. 读取数据库表结构
2. 生成完整的 CRUD 代码
3. 开箱即用，无需修改
```

### 砍掉的功能

**删除 Spec-Kit**：
- 原因：太小众，增加了复杂度
- 如果真需要，应该是单独的项目

**删除 DDD 自动生成**：
- 原因：DDD 需要手动设计，不能自动生成
- 只保留 DDD 模板供高级用户参考

**删除 `create` 命令**：
- 原因：创建项目不是我们的核心价值
- 用户可以用 `mkdir` + `go mod init`

## 🎯 最小可行产品（MVP）

### 核心功能 v1.0

只做一个功能，但要做到完美：

```bash
# 唯一的命令：从数据库生成 CRUD 代码
go-start generate --dsn="user:pass@tcp(localhost:3306)/mydb" --tables=users

# 生成的内容：
✓ Model (GORM 模型)
✓ Repository (数据访问)
✓ Service (业务逻辑)
✓ Handler (HTTP 处理器)
✓ Routes (路由注册)

# 立即可用！
cd generated-project
go mod tidy
go run main.go
curl http://localhost:8080/api/users
```

### 成功标准

1. **能跑起来**
   - 生成的代码能编译
   - 能启动服务器
   - API 能访问

2. **简单清晰**
   - 只有一个命令
   - 参数少而精
   - 文档不超过 3 页

3. **真实测试**
   - 每次修改都要测试
   - 提供 demo 数据库
   - 提供 API 测试脚本

## 📋 具体改进计划

### 第一步：修复致命问题（1小时）

1. 修复 `create` 命令的模板 bug
2. 修复模块路径硬编码问题
3. 生成一个示例项目并测试

### 第二步：重新设计命令（2小时）

```bash
# 旧设计（混乱）
go-start create my-api
go-start gen db --dsn="..." --tables=users
go-start spec generate --file=api.spec.yaml

# 新设计（清晰）
go-start generate --dsn="..." --tables=users
```

只保留一个核心命令！

### 第三步：真实测试（2小时）

```bash
# 1. 创建测试数据库
mysql -u root -p < test-db.sql

# 2. 生成代码
go-start generate --dsn="..." --tables=users

# 3. 编译
cd generated-project
go build

# 4. 运行
./generated-project

# 5. 测试 API
curl http://localhost:8080/api/users
```

### 第四步：简化文档（1小时）

只保留 3 个文档：
1. README.md（5 分钟上手）
2. EXAMPLE.md（完整示例）
3. API.md（API 参考）

删除其他所有冗余文档。

## 🚀 最终目标

**一个命令，解决一个问题**

```bash
go-start generate --dsn="..." --tables=users
```

**生成可用的 CRUD 代码，让工程师专注业务逻辑。**

就这么简单！

---

**问题**：你觉得这个方向对吗？还是你有其他想法？
