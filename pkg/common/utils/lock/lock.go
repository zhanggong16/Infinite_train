package lock

import (
	"time"
)

/*type Lock struct {
	Key	int
	Mu	sync.Mutex
}

// 当key=0，置为1，返回true
// 当key=1，返回false
func (l *Lock) GetLock() bool {
	l.Mu.Lock()
	defer l.Mu.Unlock()

	if l.Key == 0 {
		l.Key = 1
		return true
	} else {
		return false
	}
}

func (l *Lock) ReleaseLock() {
	l.Mu.Lock()
	l.Key = 0
	l.Mu.Unlock()
}*/

// 使用chan实现互斥锁
type Mutex struct {
	ch chan struct{}
}

// 使用锁需要初始化
func NewMutex() *Mutex {
	mu := &Mutex{make(chan struct{}, 1)}
	mu.ch <- struct{}{}
	return mu
}

// 请求锁，直到获取到
func (m *Mutex) Lock() {
	<-m.ch
}

// 解锁
func (m *Mutex) Unlock() {
	select {
	case m.ch <- struct{}{}:
	default:
		panic("unlock of unlocked mutex")
	}
}

// 尝试获取锁
func (m *Mutex) TryLock() bool {
	select {
	case <-m.ch:
		return true
	default:
	}
	return false
}

// 加入一个超时的设置
func (m *Mutex) LockTimeout(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	select {
	case <-m.ch:
		timer.Stop()
		return true
	case <-timer.C:
	}
	return false
}

// 锁是否已被持有
func (m *Mutex) IsLocked() bool {
	return len(m.ch) == 0
}