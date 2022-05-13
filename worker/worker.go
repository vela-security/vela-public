package worker

import (
	"fmt"
	"github.com/vela-security/vela-public/assert"
	"github.com/vela-security/vela-public/lua"
	"reflect"
)

var (
	workerT = reflect.TypeOf((*Worker)(nil)).String()
)

type Worker struct {
	lua.ProcEx
	env  assert.Environment
	code string
	name string
	task func()
	kill func()
}

func (w *Worker) CodeVM() string {
	return w.code
}

func (w *Worker) Type() string {
	return workerT
}

func (w *Worker) Name() string {
	return fmt.Sprintf("serivce.worker.%s", w.name)
}

func (w *Worker) Close() error {
	if w.kill != nil {
		w.kill()
	}
	return nil
}

func (w *Worker) Kill(kill func()) *Worker {
	if kill != nil {
		w.kill = kill
	}
	return w
}

func (w *Worker) Task(task func()) *Worker {
	if task != nil {
		w.task = task
	}
	return w
}

func (w *Worker) Env(env assert.Environment) *Worker {
	w.env = env
	return w
}

func (w *Worker) Start() error {
	if w.task == nil {
		return fmt.Errorf("%s worker not found task", w.name)
	}

	if w.kill == nil {
		return fmt.Errorf("%s worker not found kill", w.name)
	}
	w.env.Spawn(0, w.task)
	return nil
}

/*
	Worker("wakeup" , func(){} , func(){}).Async()
*/

func New(L *lua.LState, name string) *Worker {
	wk := &Worker{name: name, code: L.CodeVM()}
	proc := L.NewProc(wk.Name(), wk.Type())
	if !proc.IsNil() {
		proc.Close()
	}

	proc.Set(wk)
	return wk
}
