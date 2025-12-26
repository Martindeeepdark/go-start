# gen db 命令测试指南

## 测试目标

验证 `go-start gen db` 命令能否:
1. 从数据库读取表结构
2. 生成完整的 CRUD 代码
3. 生成的代码能够编译通过
4. 生成的代码包含预期功能

## 前置条件

1. MySQL 已安装并运行
2. MySQL root 用户无密码或修改脚本中的密码
3. 已构建 go-start CLI 工具: `go build -o bin/go-start ./cmd/go-start`

## 快速测试

### 方法 1: 自动化测试脚本 (推荐)

```bash
# 给脚本添加执行权限
chmod +x test-gen-db.sh

# 运行测试
./test-gen-db.sh
```

脚本会自动:
- 创建测试数据库和表
- 运行 gen db 命令
- 检查生成的文件
- 验证代码质量
- 尝试编译
- 询问是否清理

### 方法 2: 手动测试

#### 步骤 1: 创建测试数据库

```bash
mysql -u root -p -e "CREATE DATABASE test_gen_db; USE test_gen_db; CREATE TABLE users (id BIGINT PRIMARY KEY AUTO_INCREMENT, username VARCHAR(255) NOT NULL UNIQUE, email VARCHAR(255) NOT NULL UNIQUE, created_at DATETIME DEFAULT CURRENT_TIMESTAMP);"
```

#### 步骤 2: 运行 gen db 命令

```bash
./bin/go-start gen db \
  --dsn="root:@tcp(localhost:3306)/test_gen_db" \
  --tables=users \
  --output=./test_gen_output \
  --arch=mvc
```

#### 步骤 3: 检查生成的文件

```bash
ls -la test_gen_output/
# 应该看到:
# - model/user.go
# - repository/user.go
# - service/user.go
# - controller/user.go
# - routes.go
```

#### 步骤 4: 查看生成的代码

```bash
# 查看生成的 Model
cat test_gen_output/model/user.go

# 查看生成的 Repository
cat test_gen_output/repository/user.go

# 查看生成的 Service
cat test_gen_output/service/user.go

# 查看生成的 Controller
cat test_gen_output/controller/user.go

# 查看生成的 Routes
cat test_gen_output/routes.go
```

#### 步骤 5: 验证功能

检查生成的代码是否包含:

**Repository 层**:
- ✅ 基本 CRUD 方法
- ✅ 基于索引的查询方法 (如 `ByUsername`, `ByEmail`)
- ✅ 批量操作方法 (如 `CreateBatch`, `UpdateBatch`)

**Service 层**:
- ✅ 业务逻辑封装
- ✅ 缓存支持 (Redis)
- ✅ 事务支持

**Controller 层**:
- ✅ RESTful API 端点
- ✅ 参数验证
- ✅ 错误处理

**Routes**:
- ✅ 路由注册
- ✅ 分组管理

## 预期结果

### 成功的标志

1. **文件生成完整**
   ```
   test_gen_output/
   ├── model/
   │   └── user.go
   ├── repository/
   │   ├── user.go
   │   └── repository.go
   ├── service/
   │   ├── user.go
   │   └── service.go
   ├── controller/
   │   ├── user.go
   │   └── controller.go
   └── routes.go
   ```

2. **代码质量高**
   - 所有函数有详细的中文注释
   - 错误处理完善
   - 符合 Go 编码规范
   - 有类型安全

3. **功能完整**
   - Model: 使用 GORM Gen 的查询 API
   - Repository: CRUD + 高级查询
   - Service: 缓存 + 事务
   - Controller: RESTful API
   - Routes: 自动注册

### 可能的问题

1. **数据库连接失败**
   ```
   Error: dial tcp: lookup localhost: no such host
   ```
   解决: 检查 MySQL 是否运行,检查 DSN 格式

2. **表不存在**
   ```
   Error: Table 'xxx' doesn't exist
   ```
   解决: 检查表名是否正确,检查数据库是否选择正确

3. **生成的代码编译失败**
   ```
   Error: undefined: xxx
   ```
   解决: 检查包导入,检查类型定义

## 下一步测试

如果 gen db 命令测试通过,接下来需要:

1. **集成测试**: 将生成的代码集成到一个真实项目中
2. **API 测试**: 启动服务器,测试 API 端点
3. **端到端测试**: 从创建数据库到访问 API 的完整流程

## 清理

```bash
# 删除生成的代码
rm -rf test_gen_output

# 删除测试数据库
mysql -u root -p -e "DROP DATABASE IF EXISTS test_gen_db;"
```

## 问题反馈

如果测试失败,请记录:
1. 错误信息
2. 执行的命令
3. 数据库表结构
4. 生成的文件内容

以便诊断问题。
