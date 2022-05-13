package kind

import (
	"github.com/valyala/fastjson"
	"github.com/vela-security/vela-public/lua"
)

type Fast struct {
	value *fastjson.Value
}

func NewFastJson(data []byte) (*Fast, error) {
	v, err := fastjson.ParseBytes(data)
	if err != nil {
		return nil, err
	}

	return &Fast{value: v}, nil
}

func (f *Fast) getInt(L *lua.LState) int {
	key := L.CheckString(1)
	n := f.value.GetInt(key)
	L.Push(lua.LNumber(n))
	return 1
}

func (f *Fast) getStr(L *lua.LState) int {
	key := L.CheckString(1)
	b := f.value.GetStringBytes(key)
	L.Push(lua.LString(lua.B2S(b)))
	return 1
}

func (f *Fast) getBool(L *lua.LState) int {
	key := L.CheckString(1)
	b := f.value.GetBool(key)
	L.Push(lua.LBool(b))
	return 1
}

func (j *Fast) Get(L *lua.LState, key string) lua.LValue {
	switch key {
	case "int":
		return L.NewFunction(j.getInt)
	case "str":
		return L.NewFunction(j.getStr)
	case "bool":
		return L.NewFunction(j.getBool)
	default:
		L.RaiseError("not found %s index", key)
		return lua.LNil
	}
}
