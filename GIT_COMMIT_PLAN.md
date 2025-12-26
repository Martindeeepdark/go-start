# Git 提交计划

## 提交策略
分为 4 个逻辑提交,每个提交聚焦一个主题:

---

## 提交 1: 修复核心 Bug (P0)

**主题**: 修复 create 命令的严重 bug

**文件**:
- `cmd/go-start/create.go`
- `pkg/database/database.go`
- `pkg/cache/serialize.go` (新建)

**提交信息**:
```bash
git add cmd/go-start/create.go
git add pkg/database/database.go
git add pkg/cache/serialize.go

git commit -m "fix: 修复 create 命令的严重 bug

修复了 4 个导致 create 命令完全无法使用的严重 bug:

1. 模板渲染缺少字段错误
   - 改用 generateMVCProjectWithOptions() 传递完整配置
   - 包含 WithRedis, WithAuth, WithSwagger 字段

2. go.mod 语法错误
   - 添加缺少的闭合括号 ')'
   - 修复生成的 go.mod 无法解析问题

3. 硬编码模块路径
   - 移除 pkg/database 对 github.com/yourname/go-start 的导入
   - 内联 TxOptions 和 Stats 类型到 database.go
   - 生成的项目可以在任何模块下编译

4. 缺少序列化函数
   - 新建 pkg/cache/serialize.go
   - 提供 cache.Marshal/Unmarshal JSON 序列化函数

影响:
- ✅ create 命令现在可以生成可编译的项目
- ✅ 生成的代码不再依赖硬编码路径
- ✅ 模板渲染不再报错

修复 #1, #2, #3, #4"
```

---

## 提交 2: 改进模板系统 (P1)

**主题**: 支持条件编译和模板变量

**文件**:
- `templates/mvc/main.go.tpl`
- `templates/mvc/config/config.go.tpl`
- `templates/mvc/config.yaml.tpl`

**提交信息**:
```bash
git add templates/mvc/main.go.tpl
git add templates/mvc/config/config.go.tpl
git add templates/mvc/config.yaml.tpl

git commit -m "feat: 改进模板支持条件编译

增强了模板系统的灵活性,支持功能开关:

1. Main.go 条件编译
   - Swagger 导入和路由支持条件编译 (WithSwagger)
   - Redis 初始化支持条件编译 (WithRedis)
   - 未启用的功能不会包含在代码中

2. Config 模板优化
   - 支持 WithRedis 条件导入
   - 使用 {{.ServerPort}} 和 {{.Database}} 变量
   - Redis 配置和默认值条件包含

3. Config.yaml 模板优化
   - 使用 {{.ServerPort}} 和 {{.Database}} 替代硬编码
   - Redis 配置段条件包含

影响:
- ✅ 生成的代码更简洁,不需要的功能不会包含
- ✅ 配置更灵活,支持动态端口和数据库类型
- ✅ 减少了不必要的依赖导入

改进 #5"
```

---

## 提交 3: 完善项目模板 (P2)

**主题**: 优化 README 和 gitignore 模板

**文件**:
- `templates/mvc/README.md.tpl`
- `templates/mvc/gitignore.tpl`

**提交信息**:
```bash
git add templates/mvc/README.md.tpl
git add templates/mvc/gitignore.tpl

git commit -m "docs: 优化 README 和 gitignore 模板

1. README.md 改进
   - Redis 环境要求支持条件显示
   - Swagger 文档路径提示
   - 配置说明支持条件显示

2. gitignore 改进
   - 添加 .go.version 忽略规则
   - 完善 Go 项目常见的忽略模式

影响:
- ✅ 生成的 README 更准确地反映项目配置
- ✅ 避免提交敏感配置文件

改进 #6"
```

---

## 提交 4: 添加文档 (P2)

**主题**: 添加详细的修复说明和测试文档

**文件**:
- `FIXES_APPLIED.md`
- `TEST_CHECKLIST.md`
- `TODO_DDD.md`
- `WORK_SUMMARY.md`
- `GIT_COMMIT_PLAN.md` (本文件)

**提交信息**:
```bash
git add FIXES_APPLIED.md
git add TEST_CHECKLIST.md
git add TODO_DDD.md
git add WORK_SUMMARY.md
git add GIT_COMMIT_PLAN.md

git commit -m "docs: 添加完整的修复说明和测试文档

添加了详细的文档记录本次 Bug 修复:

1. FIXES_APPLIED.md
   - 详细记录所有 7 个 Bug 的修复过程
   - 包含错误信息、根本原因、修复方案
   - 提供代码示例和影响分析

2. TEST_CHECKLIST.md
   - 完整的测试步骤和预期结果
   - 常见问题排查指南
   - 测试清理步骤

3. TODO_DDD.md
   - DDD 架构实现计划
   - 当前状态和待完成功能
   - 实现优先级划分

4. WORK_SUMMARY.md
   - 本次工作的完整总结
   - 修改文件清单
   - 遇到的挑战和经验教训

5. GIT_COMMIT_PLAN.md
   - Git 提交计划
   - 每个提交的文件和说明

影响:
- ✅ 便于代码审查和后续维护
- ✅ 清晰的测试指南
- ✅ 完整的工作记录

文档 #7"
```

---

## 执行步骤

```bash
# 1. 查看当前状态
git status

# 2. 执行提交 1 - 核心修复
git add cmd/go-start/create.go pkg/database/database.go pkg/cache/serialize.go
git commit -m "fix: 修复 create 命令的严重 bug ..."

# 3. 执行提交 2 - 模板改进
git add templates/mvc/main.go.tpl templates/mvc/config/config.go.tpl templates/mvc/config.yaml.tpl
git commit -m "feat: 改进模板支持条件编译 ..."

# 4. 执行提交 3 - 项目模板
git add templates/mvc/README.md.tpl templates/mvc/gitignore.tpl
git commit -m "docs: 优化 README 和 gitignore 模板 ..."

# 5. 执行提交 4 - 文档
git add FIXES_APPLIED.md TEST_CHECKLIST.md TODO_DDD.md WORK_SUMMARY.md GIT_COMMIT_PLAN.md
git commit -m "docs: 添加完整的修复说明和测试文档 ..."

# 6. 查看提交历史
git log --oneline -4

# 7. (可选) 推送到远程
git push
```

---

## 修改文件总览

### 核心代码 (6 个文件)
1. ✅ `cmd/go-start/create.go` - 修复模板渲染和 go.mod
2. ✅ `pkg/database/database.go` - 移除硬编码路径
3. ✅ `pkg/cache/serialize.go` - 新建序列化函数
4. ✅ `templates/mvc/main.go.tpl` - 条件编译支持
5. ✅ `templates/mvc/config/config.go.tpl` - 条件导入
6. ✅ `templates/mvc/config.yaml.tpl` - 模板变量

### 项目模板 (2 个文件)
7. ✅ `templates/mvc/README.md.tpl` - 条件说明
8. ✅ `templates/mvc/gitignore.tpl` - 添加 .go.version

### 文档 (5 个文件)
9. ✅ `FIXES_APPLIED.md` - Bug 修复详情
10. ✅ `TEST_CHECKLIST.md` - 测试步骤
11. ✅ `TODO_DDD.md` - DDD 计划
12. ✅ `WORK_SUMMARY.md` - 工作总结
13. ✅ `GIT_COMMIT_PLAN.md` - 提交计划

**总计**: 13 个文件

---

## 注意事项

1. **提交顺序很重要**: 先修复 bug,再改进功能,最后添加文档
2. **提交信息要详细**: 每个提交都应该说明问题和影响
3. **测试后再推送**: 虽然代码已修复,但最好先测试再 push
4. **分支策略**: 如果有主分支保护,应该先创建 PR

---

## 后续工作

提交完成后:
1. 执行 TEST_CHECKLIST.md 中的测试步骤
2. 验证生成的项目可以编译和运行
3. 创建一个示例项目作为 demo
4. 实现 TODO_DDD.md 中的 DDD 功能
