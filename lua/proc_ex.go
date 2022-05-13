package lua

import "time"

//ProcEx 防止过多的方法定义
type ProcEx struct {
	Uptime time.Time
	Status ProcState
	TypeOf string
	code   string
}

func (p *ProcEx) Init(typeof string) {
	p.Status = PTInit
	p.TypeOf = typeof
}

func (p *ProcEx) vm(L *LState) {
	p.code = L.CodeVM()
}

func (p *ProcEx) CodeVM() string {
	return p.code
}

func (p *ProcEx) IsRun() bool {
	return p.Status == PTRun
}

func (p *ProcEx) IsPanic() bool {
	return p.Status == PTPanic
}

func (p *ProcEx) IsInit() bool {
	return p.Status == PTInit
}

func (p *ProcEx) IsClose() bool {
	return p.Status == PTClose
}

func (p *ProcEx) V(opts ...interface{}) {
	for _, item := range opts {

		switch v := item.(type) {

		//设置启动时间
		case time.Time:
			p.Uptime = v
		//设置类型
		case string:
			p.TypeOf = v
		//设置数据类型
		case ProcState:
			p.Status = v

		default:

		}
	}
}

func (p *ProcEx) Type() string { return p.TypeOf }
func (p *ProcEx) Name() string { return "" }

func (p *ProcEx) NewMeta(*LState, LValue, LValue)  {}
func (p *ProcEx) Meta(*LState, LValue) LValue      { return LNil }
func (p *ProcEx) Index(*LState, string) LValue     { return LNil }
func (p *ProcEx) NewIndex(*LState, string, LValue) {}

func (p *ProcEx) Show(out Console) {
	out.Println("请定义对象的Show方法 ,如： func(a *A) Show( out lua.Console)")
}
func (p *ProcEx) Help(out Console) {
	out.Println("请定义对象的Help方法 ,如： func(a *A) Help( out lua.Console)")
}

func (p *ProcEx) State() ProcState { return p.Status }

