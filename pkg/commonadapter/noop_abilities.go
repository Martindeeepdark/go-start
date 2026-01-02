package commonadapter

// NoopAuth 默认的空实现
// 返回未实现错误或固定结果，用于开发期占位。
type NoopAuth struct{}

// VerifyToken 校验令牌并返回用户ID
func (a *NoopAuth) VerifyToken(token string) (string, error) { return "", ErrCommonUnavailable }

// RequirePermission 校验权限码是否可用
func (a *NoopAuth) RequirePermission(userID string, permission string) error {
	return ErrCommonUnavailable
}

// NoopCache 默认的空实现
// 提供基本占位行为。
type NoopCache struct{}

// Get 读取缓存值
func (c *NoopCache) Get(key string) (interface{}, error) { return nil, ErrCommonUnavailable }

// Set 写入缓存值
func (c *NoopCache) Set(key string, value interface{}, ttlSeconds int) error {
	return ErrCommonUnavailable
}

// Delete 删除缓存键
func (c *NoopCache) Delete(key string) error { return ErrCommonUnavailable }

// DeleteByPattern 按模式批量删除
func (c *NoopCache) DeleteByPattern(pattern string) error { return ErrCommonUnavailable }

// NoopAudit 默认的空实现
type NoopAudit struct{}

// Record 记录审计事件
func (a *NoopAudit) Record(actor string, resource string, action string, status string, message string) error {
	return ErrCommonUnavailable
}

// NoopIdempotency 默认的空实现
type NoopIdempotency struct{}

// CheckAndSet 检查并占用幂等键
func (i *NoopIdempotency) CheckAndSet(key string, ttlSeconds int) (bool, error) {
	return false, ErrCommonUnavailable
}

// NoopLock 默认的空实现
type NoopLock struct{}

// Acquire 获取锁
func (l *NoopLock) Acquire(key string, ttlSeconds int) (bool, error) {
	return false, ErrCommonUnavailable
}

// Release 释放锁
func (l *NoopLock) Release(key string) error { return ErrCommonUnavailable }

// NoopEventBus 默认的空实现
type NoopEventBus struct{}

// Publish 发布事件
func (b *NoopEventBus) Publish(topic string, payload interface{}) error { return ErrCommonUnavailable }

// Subscribe 订阅事件
func (b *NoopEventBus) Subscribe(topic string, handler func(interface{})) error {
	return ErrCommonUnavailable
}
