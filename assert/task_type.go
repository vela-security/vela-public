package assert

import (
	"github.com/vela-security/vela-public/lua"
	"time"
)

type Runner struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Status  string `json:"status"`
	CodeVM  string `json:"code_vm"`
	Private bool   `json:"private"`
}

type Code interface {
	Get(string) *lua.ProcData

	Wrap() error
	Key() string
	Hash() string
	Link() string
	Status() string
	List() []*Runner
	Uptime() time.Time
	Exist(string) bool
	CompareVM(*lua.LState) bool
	NewProc(*lua.LState, string, string) *lua.ProcData
}

type taskByEnv interface {

	//TaskSize 统计数量
	TaskSize() int

	//LoadTask 加载服务
	LoadTask(string, []byte, interface{}) error

	//DoTask 加载文件通过字节 func(name , chunk , env , way)
	DoTask(string, []byte, Way) error

	//DoTaskFile 加载服务通过文件 func(path , env , way)
	DoTaskFile(string, Way) error

	//RegisterTask 注册任务信息 func(name , chunk , env , way)
	RegisterTask(string, []byte, Way) error

	//WakeupTask 唤醒任务 func(way)
	WakeupTask(Way) error

	//RemoveTask  删除任务 func(name , way)
	RemoveTask(string, Way) error

	//FindTask 查找task 对象
	FindTask(string) *Task

	//TaskList 查看服务内容
	TaskList() []*Task

	//FindCode 获取相关服务
	FindCode(string) Code

	//FindProc 获取相关ProcData
	FindProc(string, string) (*lua.ProcData, error)

	//WithTaskTree 设置全局任务
	WithTaskTree(interface{})
}
