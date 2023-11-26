package utils

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"
)

func Go(fn func()) {
	go goSafe(fn)
}

func goSafe(fn func(), cleaner ...func()) {
	defer Recover(cleaner...)
	fn()
}

func Recover(cleanups ...func()) {
	if cleanups != nil {
		for _, cleanup := range cleanups {
			cleanup()
		}
	}
	if err := recover(); err != nil {
		fmt.Fprintf(os.Stderr, "time:%d,  recover occurs %+v , stack:%s \n", time.Now().UnixNano(), err, string(debug.Stack()))
	}
}
