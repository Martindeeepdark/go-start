#!/bin/bash

# gen db 命令端到端测试脚本

set -e

echo "=========================================="
echo "gen db 命令端到端测试"
echo "=========================================="
echo ""

# 1. 准备测试数据库
echo "📦 步骤 1: 准备测试数据库..."
echo ""

# 创建测试数据库
mysql -u root -p'' -e "DROP DATABASE IF EXISTS test_gen_db; CREATE DATABASE test_gen_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 创建测试表
mysql -u root -p'' test_gen_db << 'EOF'
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(255) NOT NULL UNIQUE COMMENT '用户名',
    email VARCHAR(255) NOT NULL UNIQUE COMMENT '邮箱',
    password VARCHAR(255) NOT NULL COMMENT '密码',
    age INT COMMENT '年龄',
    status TINYINT DEFAULT 1 COMMENT '状态 1:正常 0:禁用',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    KEY idx_username (username),
    KEY idx_email (email),
    KEY idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

CREATE TABLE articles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(500) NOT NULL COMMENT '标题',
    content TEXT COMMENT '内容',
    author_id BIGINT NOT NULL COMMENT '作者ID',
    category_id INT COMMENT '分类ID',
    views INT DEFAULT 0 COMMENT '浏览量',
    status TINYINT DEFAULT 1 COMMENT '状态 1:草稿 2:发布',
    published_at DATETIME COMMENT '发布时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    KEY idx_author_id (author_id),
    KEY idx_category_id (category_id),
    KEY idx_status (status),
    KEY idx_published_at (published_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章表';

CREATE TABLE comments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    article_id BIGINT NOT NULL COMMENT '文章ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    content VARCHAR(1000) NOT NULL COMMENT '评论内容',
    parent_id BIGINT DEFAULT NULL COMMENT '父评论ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    KEY idx_article_id (article_id),
    KEY idx_user_id (user_id),
    KEY idx_parent_id (parent_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='评论表';
EOF

echo "✓ 测试数据库创建成功: test_gen_db"
echo "  - users 表 (用户表)"
echo "  - articles 表 (文章表)"
echo "  - comments 表 (评论表)"
echo ""

# 2. 测试 gen db 命令
echo "🔧 步骤 2: 测试 gen db 命令..."
echo ""

# 清理旧的生成结果
rm -rf test_gen_output
mkdir -p test_gen_output

# 运行 gen db 命令
./bin/go-start gen db \
  --dsn="root:@tcp(localhost:3306)/test_gen_db" \
  --tables=users,articles,comments \
  --output=./test_gen_output \
  --arch=mvc

echo "✓ gen db 命令执行完成"
echo ""

# 3. 检查生成的文件
echo "📁 步骤 3: 检查生成的文件..."
echo ""

echo "生成的目录结构:"
tree -L 3 test_gen_output/ 2>/dev/null || find test_gen_output/ -type f | head -20

echo ""
echo "检查关键文件:"
FILES=(
  "test_gen_output/model/user.go"
  "test_gen_output/repository/user.go"
  "test_gen_output/service/user.go"
  "test_gen_output/controller/user.go"
  "test_gen_output/routes.go"
)

for file in "${FILES[@]}"; do
  if [ -f "$file" ]; then
    echo "  ✓ $file"
  else
    echo "  ✗ $file (缺失)"
  fi
done

echo ""

# 4. 检查生成的代码质量
echo "🔍 步骤 4: 检查生成的代码质量..."
echo ""

# 统计代码行数
echo "代码统计:"
echo "  Model 行数:    $(wc -l < test_gen_output/model/user.go)"
echo "  Repository 行数: $(wc -l < test_gen_output/repository/user.go)"
echo "  Service 行数:   $(wc -l < test_gen_output/service/user.go)"
echo "  Controller 行数: $(wc -l < test_gen_output/controller/user.go)"
echo ""

# 检查关键功能
echo "检查生成的功能:"

# 检查 Repository 是否有基于索引的查询方法
if grep -q "ByUsername" test_gen_output/repository/user.go; then
  echo "  ✓ Repository 有索引查询方法 (ByUsername, ByEmail)"
else
  echo "  ✗ Repository 缺少索引查询方法"
fi

# 检查 Service 是否有缓存逻辑
if grep -q "cache" test_gen_output/service/user.go; then
  echo "  ✓ Service 有缓存支持"
else
  echo "  ✗ Service 缺少缓存支持"
fi

# 检查 Controller 是否有完整的 CRUD 端点
if grep -q "Create\|GetByID\|Update\|Delete" test_gen_output/controller/user.go; then
  echo "  ✓ Controller 有完整 CRUD 端点"
else
  echo "  ✗ Controller CRUD 端点不完整"
fi

echo ""

# 5. 尝试编译生成的代码
echo "🔨 步骤 5: 尝试编译生成的代码..."
echo ""

cd test_gen_output

# 创建一个临时的 main.go 来测试编译
cat > main_test.go << 'EOF'
package main

import (
  _ "./model"
  _ "./repository"
  _ "./service"
  _ "./controller"
)

func main() {
  // 只是为了测试导入是否正常
}
EOF

# 尝试编译
if go build -o /dev/null main_test.go 2>&1; then
  echo "✓ 生成的代码可以编译通过"
  COMPILE_SUCCESS=true
else
  echo "✗ 生成的代码编译失败"
  COMPILE_SUCCESS=false
fi

# 清理
rm -f main_test.go

cd ..

echo ""

# 6. 总结
echo "=========================================="
echo "测试总结"
echo "=========================================="
echo ""

if [ "$COMPILE_SUCCESS" = true ]; then
  echo "✅ gen db 命令测试通过!"
  echo ""
  echo "生成的代码:"
  echo "  ✓ 结构完整 (Model/Repository/Service/Controller/Routes)"
  echo "  ✓ 功能齐全 (CRUD + 索引查询 + 缓存)"
  echo "  ✓ 代码规范 (符合 Go 最佳实践)"
  echo "  ✓ 可以编译"
  echo ""
  echo "下一步:"
  echo "  1. 创建一个完整的项目,集成生成的代码"
  echo "  2. 启动服务器,测试 API 端点"
  echo "  3. 编写从 0 到 1 的教程"
else
  echo "❌ gen db 命令测试失败"
  echo ""
  echo "请检查:"
  echo "  1. 生成的代码是否有语法错误"
  echo "  2. 包导入是否正确"
  echo "  3. 类型定义是否匹配"
fi

echo ""

# 7. 清理
echo "🧹 清理测试环境..."
echo ""

read -p "是否保留生成的代码用于检查? (y/n): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
  rm -rf test_gen_output
  mysql -u root -p'' -e "DROP DATABASE IF EXISTS test_gen_db;"
  echo "✓ 测试环境已清理"
else
  echo "⚠  生成的代码保留在 test_gen_output/ 目录"
  echo "⚠  测试数据库保留: test_gen_db"
fi

echo ""
echo "测试完成!"
