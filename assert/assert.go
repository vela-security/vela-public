package assert

import (
	"fmt"
	"github.com/vela-security/vela-public/lua"
	"go.etcd.io/bbolt"
	"net"
	"os"
	"sync"
)

var (
	_G   Environment //缓存全局环境变量
	once sync.Once   //控制设置次数
)

type CallByEnv interface {
	P(*lua.LFunction) lua.P
	Clone(*lua.LState) *lua.LState
	Coroutine() *lua.LState
	Free(*lua.LState)
	DoString(*lua.LState, string) error
	DoFile(*lua.LState, string) error
	Start(*lua.LState, lua.LFace) Start //启动对象的构建
	Call(*lua.LState, *lua.LFunction, ...lua.LValue) error
}

type InjectByEnv interface {
	Set(string, lua.LValue)    //注入接口
	Global(string, lua.LValue) //全局注入接口
}

type NodeByEnv interface {
	ID() string
	Arch() string
	Inet() string
	Inet6() string
	Mac() string
	Edition() string
	LocalAddr() string
	WithBroker(string, net.HardwareAddr, net.IP, net.IP, string) // arch , mac , inet , inet6 , edition
}

type auxiliary interface {
	Register(Closer)
	Name() string            //当前环境的名称
	DB() *bbolt.DB           //当前环境的缓存库
	Prefix() string          //系统前缀
	ExecDir() string         //当前环境目录
	Mode() string            //当前环境模式
	IsDebug() bool           //是否调试模式
	Spawn(int, func()) error //异步执行 (delay int , task func())
	Notify()                 //监控退出信号
	Kill(os.Signal)          //退出
	Bucket(...string) Bucket //缓存
	Adt() interface{}        //审计对象
}

type Environment interface {
	TnlByEnv
	LogByEnv
	CallByEnv
	MimeByEnv
	NodeByEnv
	taskByEnv
	InjectByEnv
	RegionByEnv
	auxiliary
}

func WithEnv(env Environment) {
	if _G == nil {
		_G = env
		fmt.Println("env constructor over..")
		return
	}
	fmt.Println("env already running")
}

func GxEnv() Environment {
	return _G
}
