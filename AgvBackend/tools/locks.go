package tools

import (
	"fmt"
	"sync"
	"time"
	"turtle/core/lgr"
)

func TryLock(mutex *sync.Mutex) bool {
	if mutex.TryLock() {
		return true
	} else {
		lgr.ErrorStack("Failed to Lock mutex, because is locked different place")
		mutex.Lock()
		return false
	}
}

// TimedMutex is a custom mutex that detects long-held locks.
type TimedMutex struct {
	mu        sync.Mutex
	lockTime  time.Time
	held      bool
	mutexInfo string
	Name      string
}

// NewTimedMutex creates a new TimedMutex with a fixed stack size.
func NewTimedMutex(name string, stackSize int) *TimedMutex {
	return &TimedMutex{
		Name: name,
	}
}

// Lock acquires the mutex and records the time of acquisition.
func (self *TimedMutex) Lock(mutextLine string) {
	self.mutexInfo = mutextLine
	self.mu.Lock()
	self.lockTime = time.Now()
	self.held = true
}

// Unlock releases the mutex and clears the lock time.
func (self *TimedMutex) Unlock() {
	if self.held {
		self.held = false
		self.mu.Unlock()
	} else {
		lgr.ErrorStack("Trying to unlock unlocked mutex: ", self.Name, self.mutexInfo)
	}

}

// Monitor monitors the TimedMutex and logs if it is held too long.
func (self *TimedMutex) Monitor(duration time.Duration) {
	go func() {
		for {
			time.Sleep(duration / 2) // Check at intervals shorter than the threshold
			if self.held && time.Since(self.lockTime) > duration {
				lgr.Error(fmt.Sprintf("Warning: Mutex [%s] held for more than %v by [%s]", self.Name, self.mutexInfo, duration))
				lgr.Info(self.mutexInfo)

				lgr.Ok("Unlocking mutex")
				self.Unlock()
			}
		}
	}()
}
