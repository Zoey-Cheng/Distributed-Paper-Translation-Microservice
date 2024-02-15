package lock

import (
	"context"
	"time"
)

// 锁接口
type Locker interface {

	// 加锁
	// @param ctx - context
	// @param timeout - 超时时间
	// @return error
	Lock(ctx context.Context, timeout time.Duration) error

	// 尝试加锁
	// @param ctx - context
	// @return 是否成功, error
	TryLock(ctx context.Context) (bool, error)

	// 解锁
	// @param ctx - context
	// @return error
	UnLock(ctx context.Context) error
}
