package lua

import (
	"bytes"
	"errors"
)

var (
	overflowE = errors.New("index over flow")
	tooSmallE = errors.New("index too small")
	invalidE  = errors.New("invalid slice value")
)

type Slice []LValue

func NewSlice(cap int) Slice {
	if cap == 0 {
		return (Slice)(nil)
	}

	return make(Slice, cap)
}

func (s Slice) Type() LValueType                   { return LTSlice }
func (s Slice) AssertFloat64() (float64, bool)     { return 0, false }
func (s Slice) AssertString() (string, bool)       { return "", false }
func (s Slice) AssertFunction() (*LFunction, bool) { return nil, false }
func (s Slice) Peek() LValue                       { return s }

func (s Slice) String() string {
	n := len(s)
	if n == 0 {
		return "[]"
	}

	var buf bytes.Buffer
	buf.WriteByte('[')
	for i := 0; i < n; i++ {
		if i != 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(s[i].String())
	}
	return buf.String()
}

func (s *Slice) Append(lv LValue) {
	*s = append(*s, lv)
}

func (s *Slice) Len() int {
	return len(*s)
}

func (s *Slice) Get(idx int) LValue {
	a := *s

	if idx < 0 {
		return LNil
	}

	if idx >= len(a) {
		return LNil
	}

	return a[idx]
}

func (s *Slice) Set(idx int, val LValue) error {
	a := *s
	n := len(a)
	if idx < 0 {
		return overflowE
	}

	if idx < n {
		a[idx] = val
		*s = a
		return nil
	}

	switch val.Type() {
	case LTNil:
		return invalidE
	default:
		a = append(a, val)
		*s = a
		return nil
	}
}

func (s Slice) Meta(L *LState, key LValue) LValue {
	i, ok := key.AssertFloat64()
	if !ok {
		return LNil
	}

	var idx int

	if i < 0 {
		idx = s.Len() + int(i)
	} else {
		idx = int(i) - 1
	}
	return s.Get(idx)
}

func (s Slice) NewMeta(L *LState, key LValue, val LValue) {
	i, ok := key.AssertFloat64()
	if !ok {
		return
	}

	var idx int
	if i < 0 {
		idx = s.Len() + int(i)
	} else {
		idx = int(i) - 1
	}

	s.Set(idx, val)
}

func (s Slice) Index(L *LState, key string) LValue {
	switch key {
	case "size":
		return LInt(len(s))
	default:
		return LNil
	}
}

func (s Slice) MetaTable(L *LState, key string) LValue {
	switch key {

	}

	return LNil

}
