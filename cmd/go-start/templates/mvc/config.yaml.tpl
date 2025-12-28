# go-start 项目配置文件
#
# 💡 使用说明:
#   1. 复制此文件为 config.yaml: cp config.yaml.example config.yaml
#   2. 根据你的实际环境修改配置
#   3. 密码等敏感信息请务必修改
#   4. 生产环境建议使用环境变量覆盖敏感配置

server:
  port: {{.ServerPort}}
  # 服务器监听端口
  # 默认: 8080
  # 生产环境建议使用反向代理 (如 Nginx)

database:
  driver: {{.Database}}
  # 数据库类型: mysql 或 postgresql

  host: localhost
  # 数据库主机地址
  # 本地开发: localhost
  # Docker 环境: 使用容器名称 (如 db)

  port: 3306
  # 数据库端口
  # MySQL: 3306
  # PostgreSQL: 5432

  database: {{.ProjectName}}
  # 数据库名称
  # 请确保数据库已创建

  username: root
  # 数据库用户名
  # 建议创建专用用户而不是使用 root

  password: "YOUR_DATABASE_PASSWORD_HERE"
  # ⚠️  数据库密码 (必须修改)
  # 生产环境请使用强密码
  # 可以使用环境变量: DATABASE_PASSWORD

  charset: utf8mb4
  # 字符集 (仅 MySQL)
  # 推荐: utf8mb4 (支持 emoji 等字符)

  parse_time: true
  # 将数据库时间解析为 time.Time 类型

  loc: Local
  # 时区设置
  # Local: 使用本地时区

  max_idle_conns: 10
  # 最大空闲连接数
  # 根据并发量调整

  max_open_conns: 100
  # 最大打开连接数
  # 根据数据库服务器性能调整

  conn_max_lifetime: 3600
  # 连接最大生命周期 (秒)
  # 建议: 1 小时

  log_level: info
  # 数据库日志级别
  # 开发: info (显示所有 SQL)
  # 生产: warn (仅显示慢查询和错误)

{{if .WithRedis}}
redis:
  host: localhost
  # Redis 主机地址

  port: 6379
  # Redis 端口

  password: "YOUR_REDIS_PASSWORD_HERE"
  # ⚠️  Redis 密码 (如果 Redis 设置了密码)
  # 生产环境建议启用密码认证

  db: 0
  # Redis 数据库编号 (0-15)

  pool_size: 10
  # 连接池大小

  min_idle_conns: 5
  # 最小空闲连接数

  dial_timeout: 5
  # 连接超时 (秒)

  read_timeout: 3
  # 读超时 (秒)

  write_timeout: 3
  # 写超时 (秒)

  pool_timeout: 4
  # 连接池超时 (秒)
{{end}}

# 💡 环境变量配置示例:
#
# 在 .env 文件中或直接设置环境变量:
#
# export DATABASE_PASSWORD="your_secure_password"
# export REDIS_PASSWORD="your_redis_password"
# export SERVER_PORT="9000"
#
# 然后在代码中使用 os.Getenv() 读取这些值
