package catch

import (
	"github.com/vela-security/vela-public/auxlib"
	"runtime"
)

func StackTrace(size int) string {
	var buf []byte
	if size == 0 {
		buf = make([]byte, 8192)
	} else {
		buf = make([]byte, size)

	}
	n := runtime.Stack(buf[:], false)
	return auxlib.B2S(buf[:n])
}
