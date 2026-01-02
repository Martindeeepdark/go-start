package commonadapter

// Auth 提供认证相关能力
// 定义鉴权与用户身份校验的接口，供控制器或中间件调用。
type Auth interface {
	// VerifyToken 校验令牌并返回用户ID
	VerifyToken(token string) (userID string, err error)

	// RequirePermission 校验权限码是否可用
	RequirePermission(userID string, permission string) error
}

// Cache 提供缓存读写能力
// 封装常见的 Get/Set/Delete 与按模式删除，用于服务层缓存策略。
type Cache interface {
	// Get 读取缓存值
	Get(key string) (value interface{}, err error)

	// Set 写入缓存值（带秒级过期）
	Set(key string, value interface{}, ttlSeconds int) error

	// Delete 删除缓存键
	Delete(key string) error

	// DeleteByPattern 按模式批量删除
	DeleteByPattern(pattern string) error
}

// Audit 提供审计日志能力
// 在关键写操作记录操作人、资源、动作与结果。
type Audit interface {
	// Record 记录审计事件
	Record(actor string, resource string, action string, status string, message string) error
}

// Idempotency 提供幂等能力
// 在写操作前后通过幂等键保障请求只生效一次。
type Idempotency interface {
	// CheckAndSet 检查幂等键是否使用过，若未使用则占用
	CheckAndSet(key string, ttlSeconds int) (ok bool, err error)
}

// Lock 提供分布式锁能力
// 用于关键资源的互斥控制。
type Lock interface {
	// Acquire 获取锁
	Acquire(key string, ttlSeconds int) (ok bool, err error)

	// Release 释放锁
	Release(key string) error
}

// EventBus 提供领域事件总线能力
// 支持发布/订阅用于跨域解耦。
type EventBus interface {
	// Publish 发布事件
	Publish(topic string, payload interface{}) error

	// Subscribe 订阅事件
	Subscribe(topic string, handler func(interface{})) error
}
