package pipe

import (
	"fmt"
	"github.com/vela-security/vela-public/assert"
	"github.com/vela-security/vela-public/auxlib"
	"github.com/vela-security/vela-public/lua"
	"github.com/vela-security/vela-public/xreflect"
	"io"
)

func (px *Px) Len() int {
	return len(px.chain)
}

func (px *Px) LValue(lv lua.LValue) {
	switch lv.Type() {

	case lua.LTUserData:
		px.Object(lv.Peek().(*lua.LUserData).Data)

	case lua.LTProcData:
		px.Object(lv.Peek().(*lua.ProcData).Data)

	case lua.LTAnyData:
		px.Object(lv.Peek().(*lua.AnyData).Data)

	case lua.LTObject:
		px.Object(lv.Peek())

	case lua.LTFunction:
		px.append(px.LFunc(lv.Peek().(*lua.LFunction)))
	default:
		px.invalid("invalid pipe lua type , got %s", lv.Type().String())
	}
}

func (px *Px) Object(v interface{}) {
	fn := px.Preprocess(v)
	if fn == nil {
		return
	}

	px.append(fn)
}

func (px *Px) Preprocess(v interface{}) Fn {
	switch item := v.(type) {

	case io.Writer:
		return px.Writer(item)

	case *lua.LFunction:
		return px.LFunc(item)

	case lua.Console:
		return px.Console(item)

	case func():
		return func(_ interface{}, _ *lua.LState) error {
			item()
			return nil
		}

	case func(interface{}):
		return func(i interface{}, _ *lua.LState) error {
			item(v)
			return nil
		}

	case func() error:
		return func(i interface{}, _ *lua.LState) error {
			item()
			return nil
		}

	case func(interface{}) error:
		return func(i interface{}, _ *lua.LState) error {
			return item(i)
		}

	default:
		px.invalid("invalid pipe object")
	}

	return nil
}

func (px *Px) LFunc(fn *lua.LFunction) Fn {
	return func(v interface{}, L *lua.LState) error {
		co := px.clone(L)
		cp := px.xEnv.P(fn)
		defer px.xEnv.Free(co)
		return co.CallByParam(cp, xreflect.ToLValue(v, co))
	}
}

func (px *Px) Writer(w io.Writer) Fn {
	return func(v interface{}, _ *lua.LState) error {
		if w == nil {
			return fmt.Errorf("invalid io writer %p", w)
		}

		data, err := auxlib.ToStringE(v)
		if err != nil {
			return err
		}
		w.Write(auxlib.S2B(data))
		return nil
	}
}

func (px *Px) SetEnv(env assert.Environment) *Px {
	if env != nil {
		px.xEnv = env
	}
	return px
}

func (px *Px) Console(out lua.Console) Fn {
	return func(v interface{}, _ *lua.LState) error {
		data, err := auxlib.ToStringE(v)
		if err != nil {
			return err
		}

		out.Println(data)
		return nil
	}
}

func (px *Px) Do(arg interface{}, co *lua.LState, x func(error)) {
	n := len(px.chain)
	if n == 0 {
		return
	}

	for i := 0; i < n; i++ {
		fn := px.chain[i]
		if e := fn(arg, co); e != nil && x != nil {
			x(e)
		}
	}
}
