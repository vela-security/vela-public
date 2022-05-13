package lua

import (
	"errors"
)

var (
	AlreadyRun   = errors.New("already running")
	NotFoundCode = errors.New("not found code")
	NotFoundProc = errors.New("not found proc")
	InvalidProc  = errors.New("invalid proc")
	InvalidTree  = errors.New("invalid tree")
	NotFoundTree = errors.New("not found tree")
)

type coder interface {
	Key() string
	AssertCodeLuaStateTagFunc()
	NewProc(*LState, string, string) *ProcData
}

func (ls *LState) newProc() *ProcData {
	return &ProcData{code: ls.CodeVM()}
}

func (ls *LState) CodeVM() string {
	if ls.Code != nil {
		return ls.Code.Key()
	}
	return ""
}

func (ls *LState) CheckCodeVM(name string) bool {
	return ls.CodeVM() == name
}

func (ls *LState) NewProc(key string, typeof string) *ProcData {

	if ls.Code == nil {
		ls.RaiseError("new proc are not allowed in vm without code")
		return nil
	}

	proc := ls.Code.NewProc(ls, key, typeof)
	proc.code = ls.Code.Key()
	proc.private = false
	return proc
}

func (ls *LState) NewP(key string, x func() LFace) {
}
