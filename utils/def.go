package utils

import (
	"log"
	"runtime/debug"
)

// Def 捕获panic
func Def() {
	if err := recover(); err != nil {
		log.Printf("panic recovery：%v\n [stack]:%s\n", err, string(debug.Stack()))
	}
}
