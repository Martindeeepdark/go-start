package commonadapter

import "errors"

// ErrCommonUnavailable 表示当前未启用或不可用的 common 集成
// 在默认实现中用于指示生成功能不可用的状态。
var ErrCommonUnavailable = errors.New("common package unavailable")