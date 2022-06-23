package lua

import (
	"fmt"
)

type ProcData struct {
	private bool
	code    string
	Data    LFace
}

func NewProcData(v LFace) *ProcData {
	return &ProcData{Data: v, private: false}
}

func (pd *ProcData) String() string                     { return fmt.Sprintf("userdata: %p", pd) }
func (pd *ProcData) Type() LValueType                   { return LTProcData }
func (pd *ProcData) AssertFloat64() (float64, bool)     { return 0, false }
func (pd *ProcData) AssertString() (string, bool)       { return "", false }
func (pd *ProcData) AssertFunction() (*LFunction, bool) { return nil, false }
func (pd *ProcData) Peek() LValue                       { return pd }

func (pd *ProcData) Close() error {
	return pd.Data.Close()
}

func (pd *ProcData) Private(L *LState) {
	if !L.CheckCodeVM(pd.CodeVM()) {
		L.RaiseError("proc private with %s not allow, must be %s", L.CodeVM(), pd.CodeVM())
		return
	}
	pd.private = true
}

func (pd *ProcData) IsPrivate() bool {
	return pd.private
}

func (pd *ProcData) CodeVM() string {
	return pd.code
}

func (pd *ProcData) IsNil() bool {
	return pd.Data == nil
}

func (pd *ProcData) Set(v LFace) {
	pd.Data = v
}

//type SuperIO struct { ProcEx }
//func (s *SuperIO) Typ() string                     { return "superIO" }
//func (s *SuperIO) name() string                     { return "superIO" }
//func (s *SuperIO) Close() error                     { return NotFound  }
//func (s *SuperIO) start() error                     { return NotFound  }
//func (s *SuperIO) Write( v []byte ) error           { return NotFound  }
//func (s *SuperIO) Read() ([]byte , error)           { return nil , NotFound }
