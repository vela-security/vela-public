package xreflect

import (
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf8"

	"github.com/vela-security/vela-public/lua"
)

func check(L *lua.LState, idx int) (ref reflect.Value, mt *Metatable) {
	ud := L.CheckUserData(idx)
	ref = reflect.ValueOf(ud.Data)
	mt = &Metatable{LTable: ud.Metatable.(*lua.LTable)}
	return
}

func tostring(L *lua.LState) int {
	ud := L.CheckUserData(1)
	if stringer, ok := ud.Data.(fmt.Stringer); ok {
		L.Push(lua.LString(stringer.String()))
	} else {
		L.Push(lua.LString(ud.String()))
	}
	return 1
}

func getUnexportedName(name string) string {
	first, n := utf8.DecodeRuneInString(name)
	if n == 0 {
		return name
	}
	return string(unicode.ToLower(first)) + name[n:]
}
