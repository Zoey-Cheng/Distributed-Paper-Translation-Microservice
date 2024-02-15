package signal

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// AcquireLuaScript 和 ReleaseLuaScript 是用于 Lua 脚本的常量，用于执行 Redis 操作。
const (
	AcquireLuaScript = `local key = KEYS[1]
local max = tonumber(ARGV[1])
local num = redis.call("INCR", key)
if num > max then
    redis.call("DECR", key)
    return false
else
    return true
end`

	ReleaseLuaScript = `local key = KEYS[1]
local num = redis.call("DECR", key)
if num < 0 then
    redis.call("INCR", key)
    return nil
else
    return nil
end`
)

// ctx 是用于上下文操作的背景上下文。
var (
	ctx = context.Background()
)

// RedisSignalFactoryImpl 结构体实现了 SignalFactory 接口，用于创建 RedisSemaphore。
type RedisSignalFactoryImpl struct {
	rds redis.UniversalClient
}

// NewRedisSignalFactory 函数创建并返回一个新的 RedisSignalFactoryImpl 对象，用于创建 RedisSemaphore。
func NewRedisSignalFactory(uri string) (SignalFactory, error) {
	opt, err := redis.ParseURL(uri)
	if err != nil {
		return nil, err
	}
	return &RedisSignalFactoryImpl{rds: redis.NewClient(opt)}, nil
}

// Semaphore 方法实现了 SignalFactory 接口，用于创建并返回 RedisSemaphore。
func (r *RedisSignalFactoryImpl) Semaphore(name string, max int) Semaphore {
	return &RedisSemaphore{resourceName: name, resourceNumber: max, rds: r.rds}
}

// RedisSemaphore 结构体实现了 Semaphore 接口，使用 Redis 实现信号量。
type RedisSemaphore struct {
	resourceName   string
	resourceNumber int

	rds redis.UniversalClient
}

// Acquire 方法实现了 Semaphore 接口，用于尝试获取信号量。
func (r *RedisSemaphore) Acquire() (bool, error) {
	eval := r.rds.Eval(ctx, AcquireLuaScript, []string{r.resourceName}, r.resourceNumber)
	if eval.Err() != nil {
		if eval.Err() == redis.Nil {
			return false, nil
		}
		return false, eval.Err()
	}

	b, err := eval.Bool()
	if err != nil {
		return false, err
	}
	return b, nil
}

// Release 方法实现了 Semaphore 接口，用于释放信号量。
func (r *RedisSemaphore) Release() error {
	err := r.rds.Eval(ctx, ReleaseLuaScript, []string{r.resourceName}).Err()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}
	return nil
}

// Reset 方法实现了 Semaphore 接口，用于重置信号量的状态。
func (r *RedisSemaphore) Reset() error {
	return r.rds.Del(ctx, r.resourceName).Err()
}

// RedisSemaphoreObserver 结构体用于观察 RedisSemaphore 的状态变化。
type RedisSemaphoreObserver time.Duration

// Observer 方法实现了 SemaphoreObserver 接口，用于观察信号量的状态变化。
func (r RedisSemaphoreObserver) Observer(semaphore Semaphore) (<-chan interface{}, error) {
	ch := make(chan interface{}, 1000)
	timer := time.NewTicker(time.Duration(r))
	go func() {
		for {
			select {
			case <-timer.C:
				acquire, _ := semaphore.Acquire()
				if acquire {
					ch <- struct{}{}
				}
			}
		}
	}()
	return ch, nil
}
