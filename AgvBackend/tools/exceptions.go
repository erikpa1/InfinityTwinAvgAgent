package tools

import (
	"runtime/debug"
	"sync"
	"turtle/core/lgr"
)

func Recover(message string, onErrr ...func(any)) {
	if r := recover(); r != nil {
		lgr.ErrorStack(message, r)
		debug.PrintStack()

		for _, receiver := range onErrr {
			receiver(r)
		}

	}
}

func RecoverWithMutUnlock(mutex sync.Mutex, message string) {
	defer mutex.Unlock()
	r := recover()
	if r != nil {
		lgr.ErrorStack(message, r)
		debug.PrintStack()
	}

}
