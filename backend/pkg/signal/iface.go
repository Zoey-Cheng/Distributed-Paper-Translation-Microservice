package signal

// SignalFactory 接口定义了创建信号量的方法。
type SignalFactory interface {
	Semaphore(name string, max int) Semaphore
}

// Semaphore 信号量接口定义了信号量的基本操作。
type Semaphore interface {
	Acquire() (bool, error) // Acquire 方法尝试获取信号量。
	Release() error         // Release 方法释放信号量。
	Reset() error           // Reset 方法重置信号量的状态。
}

// SemaphoreObserver 接口定义了观察信号量的方法。
type SemaphoreObserver interface {
	Observer(Semaphore) (<-chan interface{}, error) // Observer 方法用于观察信号量的状态变化。
}
