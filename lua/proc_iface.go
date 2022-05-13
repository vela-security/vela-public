package lua

import (
	"errors"
	"io"
)

var (
	NULL     = []byte("")
	NotFound = errors.New("not function")
)

type Console interface {
	Println(string)
	Printf(string, ...interface{})
	Invalid(string, ...interface{})
}

type LFace interface {
	Name() string     //获取当前对象名称
	Type() string     //获取对象类型
	State() ProcState //获取状态
	Start() error
	Close() error
	NewMeta(*LState, LValue, LValue)  //设置字段
	Meta(*LState, LValue) LValue      //获取字段
	Index(*LState, string) LValue     //获取字符串字段 __index function
	NewIndex(*LState, string, LValue) //设置字段 __newindex function

	Show(Console) //控制台打印
	Help(Console) //控制台 辅助信息
	//V  设置信息
	V(...interface{})
}

type Writer interface {
	LFace
	io.Writer
}

type IO interface {
	LFace
	io.Writer
	io.Reader
}

type Reader interface {
	LFace
	io.Reader
}

type Closer interface {
	LFace
	io.Closer
}

type ReaderCloser interface {
	LFace
	io.Reader
	io.Closer
}

type WriterCloser interface {
	LFace
	io.Writer
	io.Closer
}

func CheckIO(val *ProcData) IO {
	obj, ok := val.Data.(IO)
	if ok {
		return obj
	}
	return nil
}

func CheckWriter(val *ProcData) Writer {
	obj, ok := val.Data.(Writer)
	if ok {
		return obj
	}
	return nil
}

func CheckReader(val *ProcData) Reader {
	obj, ok := val.Data.(Reader)
	if ok {
		return obj
	}
	return nil
}

func CheckCloser(val *ProcData) Closer {
	obj, ok := val.Data.(Closer)
	if !ok {
		return nil
	}

	return obj
}
