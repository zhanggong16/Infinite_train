package lock

import "sync"

var GetSchedulerLockFailed = "get scheduler lock failed, skip this task"

type Lock struct {
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
}