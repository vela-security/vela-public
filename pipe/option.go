package pipe

import (
	"github.com/vela-security/vela-public/assert"
)

func Seek(n int) func(*Px) {
	return func(px *Px) {
		if n < 0 {
			return
		}
		px.seek = n
	}
}

func Env(env assert.Environment) func(*Px) {
	return func(px *Px) {
		px.xEnv = env
	}
}
