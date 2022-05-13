package lua

const (
	PTInit ProcState = iota
	PTRun
	PTErr
	PTClose
	PTPanic
	PTPrivate
	PTMode
)

type ProcState uint32

var ProcStateValue = [...]string{"init", "run", "error", "close", "panic",  "private", "mode"}

func (pv ProcState) String() string {
	return ProcStateValue[int(pv)]
}

