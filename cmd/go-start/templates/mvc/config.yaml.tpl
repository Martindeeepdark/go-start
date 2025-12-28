server:
  port: {{.ServerPort}}

database:
  driver: {{.Database}}
  host: localhost
  port: 3306
  database: {{.ProjectName}}
  username: root
  password: ""
  charset: utf8mb4
  parse_time: true
  loc: Local
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600
  log_level: info

{{if .WithRedis}}
redis:
  host: localhost
  port: 6379
  password: ""
  db: 0
  pool_size: 10
  min_idle_conns: 5
  dial_timeout: 5
  read_timeout: 3
  write_timeout: 3
  pool_timeout: 4
{{end}}