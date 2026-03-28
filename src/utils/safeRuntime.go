package utils

import (
	"log"
	"runtime/debug"
)

func SafeRuntime(name string, f func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[panic] %s recovered: %v", name, r)
				log.Printf("[stack] %s trace: %s", name, string(debug.Stack()))
			}
		}()
		f()
	}()
}
