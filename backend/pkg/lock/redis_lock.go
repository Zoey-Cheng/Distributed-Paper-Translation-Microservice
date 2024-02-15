package lock

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

/**
* Redis解锁脚本
 */
const UnLockLua = `
  if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
  else
    return 0
  end
`

/**
* Redis锁
 */
type RedisLocker struct {
	LockName    string        // 锁名称
	ExpireDur   time.Duration // 过期时间
	Token       string        // 加锁token
	redisClient *redis.Client // Redis客户端
}

/**
* 创建Redis锁
* @param redisClient - Redis客户端
* @param lockName - 锁名称
* @param expireDur - 过期时间
* @return RedisLocker
 */
func NewRedisLocker(
	redisClient *redis.Client,
	lockName string,
	expireDur time.Duration,
) *RedisLocker {

	return &RedisLocker{
		LockName:    lockName,
		ExpireDur:   expireDur,
		redisClient: redisClient,
	}
}

/**
* 加锁
* @param ctx - Context
* @param timeout - 超时时间
* @return error
 */
func (l *RedisLocker) Lock(ctx context.Context, timeout time.Duration) error {

	timer := time.NewTimer(timeout)                  // 定时器
	ticker := time.NewTicker(time.Microsecond * 500) // 打点器

	defer func() {
		timer.Stop()
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C: // 打点期间
			locked, err := l.TryLock(ctx) // 尝试加锁
			if err != nil {
				return err
			}
			if locked {
				return nil
			}
		case <-timer.C: // 超时
			return errors.New("lock timeout")
		}
	}
}

/**
* 尝试加锁
* @param ctx - Context
* @return 是否成功, error
 */
func (l *RedisLocker) TryLock(ctx context.Context) (bool, error) {

	if l.Token == "" {
		l.Token = uuid.NewString() // 生成token
	}

	cmd := l.redisClient.SetNX( // 设置锁
		ctx,
		fmt.Sprintf("redis-lock-%s", l.LockName),
		l.Token,
		l.ExpireDur,
	)

	if cmd.Err() != nil {
		return false, cmd.Err()
	}

	return cmd.Val(), nil
}

/**
* 解锁
* @param ctx - Context
* @return error
 */
func (l *RedisLocker) UnLock(ctx context.Context) error {

	cmd := l.redisClient.Eval(ctx, UnLockLua, []string{fmt.Sprintf("redis-lock-%s", l.LockName)}, l.Token) // 执行解锁脚本

	return cmd.Err()
}
