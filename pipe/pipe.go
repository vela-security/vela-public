package pipe

import (
	"github.com/vela-security/vela-public/assert"
	"github.com/vela-security/vela-public/lua"
)

type Fn func(interface{}, *lua.LState) error

type Px struct {
	chain []Fn
	seek  int
	xEnv  assert.Environment
}

func (px *Px) clone(co *lua.LState) *lua.LState {
	if px.xEnv == nil {
		px.xEnv = assert.GxEnv()
	}

	if co == nil {
		return px.xEnv.Coroutine()
	}

	return px.xEnv.Clone(co)
}
func (px *Px) append(v Fn) {
	if v == nil {
		return
	}

	px.chain = append(px.chain, v)
}

func (px *Px) coroutine() *lua.LState {
	if px.xEnv != nil {
		return px.xEnv.Coroutine()
	}
	return assert.GxEnv().Coroutine()
}

func (px *Px) free(co *lua.LState) {
	if px.xEnv != nil {
		px.xEnv.Free(co)
		return
	}
	assert.GxEnv().Free(co)
}

func (px *Px) invalid(format string, v ...interface{}) {
	if px.xEnv == nil {
		assert.GxEnv().Errorf(format, v...)
		return
	}

	px.xEnv.Errorf(format, v...)
}

func New(opt ...func(*Px)) (px *Px) {
	px = &Px{}

	for _, fn := range opt {
		fn(px)
	}

	return
}
