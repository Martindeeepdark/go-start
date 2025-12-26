# go-start 文档索引

最后更新: 2025-12-26

## 📚 用户文档

### 快速开始
- [README.md](README.md) - 项目介绍和快速开始
- [QUICKSTART.md](QUICKSTART.md) - 5分钟快速上手指南

### 功能文档
- [ARCHITECTURE.md](ARCHITECTURE.md) - 架构设计说明 (MVC/DDD)
- [docs/DDD_GUIDE.md](docs/DDD_GUIDE.md) - DDD 架构详细指南
- [docs/GORM_GEN_GUIDE.md](docs/GORM_GEN_GUIDE.md) - GORM Gen 使用指南

## 🔧 开发文档

### 设计文档
- [DESIGN.md](DESIGN.md) - 系统设计文档
- [docs/DESIGN_REVIEW.md](docs/DESIGN_REVIEW.md) - 设计评审记录
- [docs/GORM_TECH_CHOICE.md](docs/GORM_TECH_CHOICE.md) - GORM 技术选型说明
- [SPEC.md](SPEC.md) - Spec-Kit 设计规范

### 项目结构
- [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) - 项目目录结构说明

## 📊 项目状态

### 当前进度
- [TEST_RESULTS.md](TEST_RESULTS.md) - **最新测试结果** (create 命令)
- [GEN_DB_TEST_REPORT.md](GEN_DB_TEST_REPORT.md) - **最新测试报告** (gen db 命令)

### 历史记录（可删除）
- [docs/PROJECT_STATUS.md](docs/PROJECT_STATUS.md) - 项目状态记录（已过时）
- [docs/PROJECT_SUMMARY.md](docs/PROJECT_SUMMARY.md) - 项目总结（已过时）
- [WORK_SUMMARY.md](WORK_SUMMARY.md) - 工作总结（已过时）

## 🐛 问题修复

### 修复记录
- [FIXES_APPLIED.md](FIXES_APPLIED.md) - Bug 修复记录
- [GIT_COMMIT_PLAN.md](GIT_COMMIT_PLAN.md) - Git 提交计划

### 测试检查清单
- [TEST_CHECKLIST.md](TEST_CHECKLIST.md) - 测试检查清单
- [docs/FEATURE_CHECKLIST.md](docs/FEATURE_CHECKLIST.md) - 功能完成度清单
- [docs/TEST_GEN_DB.md](docs/TEST_GEN_DB.md) - gen db 命令测试指南

## 🎯 未来计划

### 待办事项
- [TODO_DDD.md](TODO_DDD.md) - DDD 功能待办
- [WIZARD.md](WIZARD.md) - Wizard 向导功能计划

### 示例项目
- [docs/COMPLETE_EXAMPLE.md](docs/COMPLETE_EXAMPLE.md) - 完整示例项目说明

---

## 🗂️ 推荐文档结构

### 根目录保留（核心文档）
```
README.md                    # 必留 - 项目主页
QUICKSTART.md               # 必留 - 快速开始
ARCHITECTURE.md             # 必留 - 架构说明
DESIGN.md                   # 保留 - 设计文档
TEST_RESULTS.md             # 必留 - 最新测试结果
GEN_DB_TEST_REPORT.md       # 必留 - gen db 测试报告
```

### docs 目录（详细文档）
```
docs/
├── DDD_GUIDE.md            # 必留 - DDD 详细指南
├── GORM_GEN_GUIDE.md       # 必留 - GORM Gen 使用
├── COMPLETE_EXAMPLE.md     # 保留 - 示例项目
└── TEST_GEN_DB.md          # 保留 - gen db 测试指南
```

### 可以删除的文档（冗余/过时）
```
FINAL_STATUS.md             # 删除 - 与 TEST_RESULTS.md 重复
FIXES_APPLIED.md            # 删除 - 历史记录，不重要
GIT_COMMIT_PLAN.md          # 删除 - 临时计划，已完成
PROJECT_STRUCTURE.md        # 删除 - 可合并到 ARCHITECTURE.md
SPEC.md                     # 删除 - Spec-Kit 未实现
TEST_CHECKLIST.md           # 删除 - 已经过时
TODO_DDD.md                 # 删除 - 放到 GitHub Issues
WIZARD.md                   # 删除 - 放到 GitHub Issues
WORK_SUMMARY.md             # 删除 - 历史记录

docs/DESIGN_REVIEW.md       # 删除 - 历史记录
docs/FEATURE_CHECKLIST.md   # 删除 - 已过时
docs/GORM_TECH_CHOICE.md    # 删除 - 可合并到 DESIGN.md
docs/PROJECT_STATUS.md      # 删除 - 已过时
docs/PROJECT_SUMMARY.md     # 删除 - 已过时
```

---

## 📝 文档使用建议

### 新手入门
1. 先看 README.md 了解项目
2. 再看 QUICKSTART.md 快速上手
3. 遇到问题看 TEST_RESULTS.md 了解当前状态

### 深入学习
1. ARCHITECTURE.md - 理解架构设计
2. docs/DDD_GUIDE.md - 学习 DDD 模式
3. docs/GORM_GEN_GUIDE.md - 掌握 GORM Gen

### 贡献代码
1. DESIGN.md - 了解设计思路
2. GEN_DB_TEST_REPORT.md - 查看测试方法
3. docs/TEST_GEN_DB.md - 学习如何测试
