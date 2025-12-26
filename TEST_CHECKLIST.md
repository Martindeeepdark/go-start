# 测试清单

## 前置准备

```bash
# 1. 构建 CLI 工具
go build -o bin/go-start ./cmd/go-start

# 2. 验证构建成功
./bin/go-start --help
```

## 测试 1: 基本项目创建 (MVC)

```bash
# 创建项目
./bin/go-start create test-mvc --arch=mvc

# 进入项目目录
cd test-mvc

# 下载依赖
go mod tidy

# 检查 go.mod 是否正确
cat go.mod
# 应该看到:
# - module github.com/yourname/test-mvc
# - 闭合括号 )
# - Redis 依赖
# - JWT 依赖
# - Swagger 依赖

# 尝试编译
go build -o server cmd/server/main.go

# 如果编译成功,检查生成的文件
ls -la
```

## 预期结果

✅ go.mod 有正确的闭合括号
✅ 所有依赖下载成功
✅ 项目可以编译

## 测试 2: 检查生成的代码

```bash
# 检查 config/config.go
grep "Redis" config/config.go
# 应该看到 Redis 字段(因为 WithRedis=true)

# 检查 main.go
grep "swagger" cmd/server/main.go
# 应该看到 Swagger 导入和路由(因为 WithSwagger=true)

# 检查 config.yaml.example
cat config.yaml.example
# 应该看到 Redis 配置段
```

## 测试 3: 运行项目(需要数据库)

```bash
# 1. 复制配置文件
cp config.yaml.example config.yaml

# 2. 编辑配置,设置数据库连接
# vim config.yaml

# 3. 运行项目(需要先创建数据库)
go run cmd/server/main.go

# 4. 测试健康检查端点
curl http://localhost:8080/health

# 5. 测试 Swagger 文档(如果启用)
curl http://localhost:8080/swagger/index.html
```

## 测试 4: gen db 命令

```bash
# 准备测试数据库
mysql -u root -p -e "CREATE DATABASE test_gen; USE test_gen; CREATE TABLE users (id INT PRIMARY KEY AUTO_INCREMENT, username VARCHAR(255), email VARCHAR(255));"

# 从数据库生成代码
./bin/go-start gen db --dsn="root:@tcp(localhost:3306)/test_gen" --tables=users

# 检查生成的文件
ls -la generated/
# 应该看到:
# - model/user.go
# - repository/user.go
# - service/user.go
# - controller/user.go
# - routes.go
```

## 常见问题排查

### 问题 1: go.mod 语法错误
**错误**: `go.mod:18: syntax error (unterminated block`
**原因**: 缺少闭合括号
**已修复**: ✅ create.go:499

### 问题 2: 找不到 github.com/yourname/go-start
**错误**: `module github.com/yourname/go-start/pkg/database/defs: git ls-remote failed`
**原因**: 硬编码的模块路径
**已修复**: ✅ database.go 内联类型

### 问题 3: undefined: cache.Unmarshal
**错误**: `undefined: cache.Unmarshal`
**原因**: 缺少序列化函数
**已修复**: ✅ cache/serialize.go

### 问题 4: 模板渲染失败
**错误**: `can't evaluate field WithRedis`
**原因**: 模板数据结构不完整
**已修复**: ✅ 使用 generateMVCProjectWithOptions()

## 完成后清理

```bash
# 清理测试项目
cd ..
rm -rf test-mvc test-gen

# 清理测试数据库
mysql -u root -p -e "DROP DATABASE test_gen;"
```
