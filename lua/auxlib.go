package lua

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

/* checkType {{{ */

func (ls *LState) CheckAny(n int) LValue {
	if n > ls.GetTop() {
		ls.ArgError(n, "value expected")
	}
	return ls.Get(n)
}

func (ls *LState) CheckInt(n int) int {
	v := ls.Get(n)
	if intv, ok := v.(LNumber); ok {
		return int(intv)
	}

	if intv, ok := v.(LInt); ok {
		return int(intv)
	}

	ls.TypeError(n, LTNumber)
	return 0
}

func (ls *LState) CheckIntOrDefault(n int, d int) int {
	v := ls.Get(n)
	if intv, ok := v.(LNumber); ok {
		return int(intv)
	}

	ls.TypeError(n, LTNumber)
	return d
}

func (ls *LState) CheckInt64(n int) int64 {
	v := ls.Get(n)
	if intv, ok := v.(LNumber); ok {
		return int64(intv)
	}
	ls.TypeError(n, LTNumber)
	return 0
}

func (ls *LState) CheckNumber(n int) LNumber {
	v := ls.Get(n)
	if lv, ok := v.(LNumber); ok {
		return lv
	}
	ls.TypeError(n, LTNumber)
	return 0
}

func (ls *LState) CheckString(n int) string {
	v := ls.Get(n)
	if lv, ok := v.(LString); ok {
		return string(lv)
	} else if LVCanConvToString(v) {
		return ls.ToString(n)
	}
	ls.TypeError(n, LTString)
	return ""
}

func IsString(v LValue) string {
	d, ok := v.AssertString()
	if ok {
		return d
	}

	return ""
}

func IsTrue(v LValue) bool {
	if lv, ok := v.(LBool); ok {
		return bool(lv) == true
	}
	return false
}

func IsFalse(v LValue) bool {
	if lv, ok := v.(LBool); ok {
		return bool(lv) == false
	}
	return false
}

func IsNumber(v LValue) LNumber {
	if lv, ok := v.(LNumber); ok {
		return lv
	}
	return 0
}

func IsInt(v LValue) int {
	if intv, ok := v.(LNumber); ok {
		return int(intv)
	}

	if intv, ok := v.(LInt); ok {
		return int(intv)
	}

	return 0
}

func IsFunc(v LValue) *LFunction {
	fn, _ := v.(*LFunction)
	return fn
}

func (ls *LState) IsTrue(n int) bool {
	return IsTrue(ls.Get(n))
}

func (ls *LState) IsFalse(n int) bool {
	return IsFalse(ls.Get(n))
}

func (ls *LState) IsNumber(n int) LNumber {
	return IsNumber(ls.Get(n))
}

func (ls *LState) IsInt(n int) int {
	return IsInt(ls.Get(n))
}

func (ls *LState) IsFunc(n int) *LFunction {
	return IsFunc(ls.Get(n))
}

func (ls *LState) IsString(n int) string {
	return IsString(ls.Get(n))
}

func (ls *LState) CheckBool(n int) bool {
	v := ls.Get(n)
	if lv, ok := v.(LBool); ok {
		return bool(lv)
	}
	ls.TypeError(n, LTBool)
	return false
}

func (ls *LState) CheckTable(n int) *LTable {
	v := ls.Get(n)
	if lv, ok := v.(*LTable); ok {
		return lv
	}
	ls.TypeError(n, LTTable)
	return nil
}

func (ls *LState) CheckFunction(n int) *LFunction {
	v := ls.Get(n)
	if lv, ok := v.(*LFunction); ok {
		return lv
	}
	ls.TypeError(n, LTFunction)
	return nil
}

func (ls *LState) CheckUserData(n int) *LUserData {
	v := ls.Get(n)
	if lv, ok := v.(*LUserData); ok {
		return lv
	}
	ls.TypeError(n, LTUserData)
	return nil
}

func (ls *LState) CheckProcData(n int) *ProcData {
	v := ls.Get(n)
	if lv, ok := v.(*ProcData); ok {
		return lv
	}
	ls.TypeError(n, LTProcData)
	return nil
}

func (ls *LState) CheckAnyData(n int) *AnyData {
	v := ls.Get(n)
	if lv, ok := v.(*AnyData); ok {
		return lv
	}
	ls.TypeError(n, LTAnyData)
	return nil
}

func (ls *LState) CheckObject(n int) LValue {
	lv := ls.Get(n)

	if lv.Type() != LTObject {
		ls.TypeError(n, LTObject)
		return nil
	}
	return lv
}

func (ls *LState) CheckThread(n int) *LState {
	v := ls.Get(n)
	if lv, ok := v.(*LState); ok {
		return lv
	}
	ls.TypeError(n, LTThread)
	return nil
}

func (ls *LState) CheckType(n int, typ LValueType) {
	v := ls.Get(n)
	if v.Type() != typ {
		ls.TypeError(n, typ)
	}
}

func (ls *LState) CheckTypes(n int, typs ...LValueType) {
	vt := ls.Get(n).Type()
	for _, typ := range typs {
		if vt == typ {
			return
		}
	}
	buf := []string{}
	for _, typ := range typs {
		buf = append(buf, typ.String())
	}
	ls.ArgError(n, strings.Join(buf, " or ")+" expected, got "+ls.Get(n).Type().String())
}

func (ls *LState) CheckOption(n int, options []string) int {
	str := ls.CheckString(n)
	for i, v := range options {
		if v == str {
			return i
		}
	}
	ls.ArgError(n, fmt.Sprintf("invalid option: %s (must be one of %s)", str, strings.Join(options, ",")))
	return 0
}

func (ls *LState) CheckSocket(n int) string {
	v := ls.CheckString(n)
	if e := CheckSocket(v); e != nil {
		ls.RaiseError("must be socket , got fail , error:%v", e)
		return ""
	}
	return v

}

func (ls *LState) CheckSockets(n int) string {
	v := ls.CheckString(n)
	arr := strings.Split(v, ",")

	var err error
	for _, item := range arr {
		err = CheckSocket(item)
		if err != nil {
			ls.RaiseError("%s error: %v", err)
			return ""
		}
	}

	return v
}

func (ls *LState) checkFile(n int) string {
	v := ls.CheckString(n)

	_, err := os.Stat(v)
	if os.IsNotExist(err) {
		ls.RaiseError("not found %s file", v)
		return ""
	}

	return v
}

/* }}} */

/* optType {{{ */

func (ls *LState) OptInt(n int, d int) int {
	v := ls.Get(n)
	if v == LNil {
		return d
	}
	if intv, ok := v.(LNumber); ok {
		return int(intv)
	}
	ls.TypeError(n, LTNumber)
	return 0
}

func (ls *LState) OptInt64(n int, d int64) int64 {
	v := ls.Get(n)
	if v == LNil {
		return d
	}
	if intv, ok := v.(LNumber); ok {
		return int64(intv)
	}
	ls.TypeError(n, LTNumber)
	return 0
}

func (ls *LState) OptNumber(n int, d LNumber) LNumber {
	v := ls.Get(n)
	if v == LNil {
		return d
	}
	if lv, ok := v.(LNumber); ok {
		return lv
	}
	ls.TypeError(n, LTNumber)
	return 0
}

func (ls *LState) OptString(n int, d string) string {
	v := ls.Get(n)
	if v == LNil {
		return d
	}
	if lv, ok := v.(LString); ok {
		return string(lv)
	}
	ls.TypeError(n, LTString)
	return ""
}

func (ls *LState) OptBool(n int, d bool) bool {
	v := ls.Get(n)
	if v == LNil {
		return d
	}
	if lv, ok := v.(LBool); ok {
		return bool(lv)
	}
	ls.TypeError(n, LTBool)
	return false
}

func (ls *LState) OptTable(n int, d *LTable) *LTable {
	v := ls.Get(n)
	if v == LNil {
		return d
	}
	if lv, ok := v.(*LTable); ok {
		return lv
	}
	ls.TypeError(n, LTTable)
	return nil
}

func (ls *LState) OptFunction(n int, d *LFunction) *LFunction {
	v := ls.Get(n)
	if v == LNil {
		return d
	}
	if lv, ok := v.(*LFunction); ok {
		return lv
	}
	ls.TypeError(n, LTFunction)
	return nil
}

func (ls *LState) OptUserData(n int, d *LUserData) *LUserData {
	v := ls.Get(n)
	if v == LNil {
		return d
	}
	if lv, ok := v.(*LUserData); ok {
		return lv
	}
	ls.TypeError(n, LTUserData)
	return nil
}

func (ls *LState) OptLightUserData(n int, d *ProcData) *ProcData {
	v := ls.Get(n)
	if v == LNil {
		return d
	}
	if lv, ok := v.(*ProcData); ok {
		return lv
	}
	ls.TypeError(n, LTProcData)
	return nil
}

/* }}} */

/* error operations {{{ */

func (ls *LState) ArgError(n int, message string) {
	ls.RaiseError("bad argument #%v to %v (%v)", n, ls.rawFrameFuncName(ls.currentFrame), message)
}

func (ls *LState) TypeError(n int, typ LValueType) {
	ls.RaiseError("bad argument #%v to %v (%v expected, got %v)", n, ls.rawFrameFuncName(ls.currentFrame), typ.String(), ls.Get(n).Type().String())
}

/* }}} */

/* debug operations {{{ */

func (ls *LState) Where(level int) string {
	return ls.where(level, false)
}

/* }}} */

/* table operations {{{ */

func (ls *LState) FindTable(obj *LTable, n string, size int) LValue {
	names := strings.Split(n, ".")
	curobj := obj
	for _, name := range names {
		if curobj.Type() != LTTable {
			return LNil
		}
		nextobj := ls.RawGet(curobj, LString(name))
		if nextobj == LNil {
			tb := ls.CreateTable(0, size)
			ls.RawSet(curobj, LString(name), tb)
			curobj = tb
		} else if nextobj.Type() != LTTable {
			return LNil
		} else {
			curobj = nextobj.(*LTable)
		}
	}
	return curobj
}

/* }}} */

/* register operations {{{ */

func (ls *LState) RegisterModule(name string, funcs map[string]LGFunction) LValue {
	tb := ls.FindTable(ls.Get(RegistryIndex).(*LTable), "_LOADED", 1)
	mod := ls.GetField(tb, name)
	if mod.Type() != LTTable {
		newmod := ls.FindTable(ls.Get(GlobalsIndex).(*LTable), name, len(funcs))
		if newmodtb, ok := newmod.(*LTable); !ok {
			ls.RaiseError("key conflict for module(%v)", name)
		} else {
			for fname, fn := range funcs {
				newmodtb.RawSetString(fname, ls.NewFunction(fn))
			}
			ls.SetField(tb, name, newmodtb)
			return newmodtb
		}
	}
	return mod
}

func (ls *LState) SetFuncs(tb *LTable, funcs map[string]LGFunction, upvalues ...LValue) *LTable {
	for fname, fn := range funcs {
		tb.RawSetString(fname, ls.NewClosure(fn, upvalues...))
	}
	return tb
}

/* }}} */

/* metatable operations {{{ */

func (ls *LState) NewTypeMetatable(typ string) *LTable {
	regtable := ls.Get(RegistryIndex)
	mt := ls.GetField(regtable, typ)
	if tb, ok := mt.(*LTable); ok {
		return tb
	}
	mtnew := ls.NewTable()
	ls.SetField(regtable, typ, mtnew)
	return mtnew
}

func (ls *LState) GetMetaField(obj LValue, event string) LValue {
	return ls.metaOp1(obj, event)
}

func (ls *LState) GetTypeMetatable(typ string) LValue {
	return ls.GetField(ls.Get(RegistryIndex), typ)
}

func (ls *LState) CallMeta(obj LValue, event string) LValue {
	op := ls.metaOp1(obj, event)
	if op.Type() == LTFunction {
		ls.reg.Push(op)
		ls.reg.Push(obj)
		ls.Call(1, 1)
		return ls.reg.Pop()
	}
	return LNil
}

/* }}} */

/* load and function call operations {{{ */

func (ls *LState) LoadFile(path string) (*LFunction, error) {
	var file *os.File
	var err error
	if len(path) == 0 {
		file = os.Stdin
	} else {
		file, err = os.Open(path)
		defer file.Close()
		if err != nil {
			return nil, newApiErrorE(ApiErrorFile, err)
		}
	}

	reader := bufio.NewReader(file)
	// get the first character.
	c, err := reader.ReadByte()
	if err != nil && err != io.EOF {
		return nil, newApiErrorE(ApiErrorFile, err)
	}
	if c == byte('#') {
		// Unix exec. file?
		// skip first line
		_, err, _ = readBufioLine(reader)
		if err != nil {
			return nil, newApiErrorE(ApiErrorFile, err)
		}
	}

	if err != io.EOF {
		// if the file is not empty,
		// unread the first character of the file or newline character(readBufioLine's last byte).
		err = reader.UnreadByte()
		if err != nil {
			return nil, newApiErrorE(ApiErrorFile, err)
		}
	}

	return ls.Load(reader, path)
}

func (ls *LState) LoadString(source string) (*LFunction, error) {
	return ls.Load(strings.NewReader(source), "<string>")
}

func (ls *LState) DoFile(path string) error {
	if fn, err := ls.LoadFile(path); err != nil {
		return err
	} else {
		ls.Push(fn)
		return ls.PCall(0, MultRet, nil)
	}
}

func (ls *LState) DoString(source string) error {
	if fn, err := ls.LoadString(source); err != nil {
		return err
	} else {
		ls.Push(fn)
		return ls.PCall(0, MultRet, nil)
	}
}

/* }}} */

/* GopherLua original APIs {{{ */

// ToStringMeta returns string representation of given LValue.
// This method calls the `__tostring` metaFunc method if defined.
func (ls *LState) ToStringMeta(lv LValue) LValue {
	if fn, ok := ls.metaOp1(lv, "__tostring").AssertFunction(); ok {
		ls.Push(fn)
		ls.Push(lv)
		ls.Call(1, 1)
		return ls.reg.Pop()
	} else {
		return LString(lv.String())
	}
}

// V a module loader to the package.preload table.
func (ls *LState) PreloadModule(name string, loader LGFunction) {
	preload := ls.GetField(ls.GetField(ls.Get(EnvironIndex), "package"), "preload")
	if _, ok := preload.(*LTable); !ok {
		ls.RaiseError("package.preload must be a table")
	}
	ls.SetField(preload, name, ls.NewFunction(loader))
}

// Checks whether the given Index is an LChannel and returns this channel.
func (ls *LState) CheckChannel(n int) chan LValue {
	v := ls.Get(n)
	if ch, ok := v.(LChannel); ok {
		return (chan LValue)(ch)
	}
	ls.TypeError(n, LTChannel)
	return nil
}

// If the given Index is a LChannel, returns this channel. If this argument is absent or is nil, returns ch. Otherwise, raises an error.
func (ls *LState) OptChannel(n int, ch chan LValue) chan LValue {
	v := ls.Get(n)
	if v == LNil {
		return ch
	}
	if ch, ok := v.(LChannel); ok {
		return (chan LValue)(ch)
	}
	ls.TypeError(n, LTChannel)
	return nil
}

/* }}} */

//
