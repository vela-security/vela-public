package kind

import (
	"context"
	"fmt"
	"github.com/vela-security/vela-public/assert"
	"github.com/vela-security/vela-public/auxlib"
	"github.com/vela-security/vela-public/catch"
	"net"
	"sync/atomic"
	"time"
)

type logger interface {
	Error(...interface{})
	Errorf(string, ...interface{})
	Infof(string, ...interface{})
}

type Listener struct {
	done uint32
	xEnv assert.Environment
	bind auxlib.URL
	fd   []net.Listener
	ch   chan net.Conn
	ctx  context.Context
	stop context.CancelFunc
}

func (ln *Listener) CloseActiveConn() {
	ln.stop()
	ln.ctx, ln.stop = context.WithCancel(context.Background())
}

func (ln *Listener) Done() bool {
	return atomic.LoadUint32(&ln.done) == 1
}

func (ln *Listener) shutdown() {
	atomic.StoreUint32(&ln.done, 1)
}

func (ln *Listener) handle(fn accept, fd net.Listener) error {
	var delay time.Duration

	conn, err := fd.Accept()
	if err == nil {
		ctx, stop := context.WithCancel(ln.ctx) //关闭子线程中的请求
		if e := fn(ctx, conn, stop); e != nil {
			ln.xEnv.Errorf("%s listen handler failure , error %v", fd.Addr().String(), e)
		}
		return nil
	}

	if ne, ok := err.(net.Error); ok && ne.Temporary() {
		if delay == 0 {
			delay = 5 * time.Millisecond
		} else {
			delay *= 2
		}

		if max := 1 * time.Second; delay > max {
			delay = max
		}
		ln.xEnv.Errorf("%s accept error: %v , retrying in %v", fd.Addr().String(), err, delay)
		time.Sleep(delay)
		return nil
	}

	return err
}

func (ln *Listener) singleH(fn accept, fd net.Listener) error {
	defer fd.Close()

	for {
		select {
		case <-ln.ctx.Done():
			ln.xEnv.Errorf("%s exit", fd.Addr().String())
			<-time.After(100 * time.Millisecond) // 等待 历史链接关闭
			return nil

		default:
			if ln.Done() {
				return fmt.Errorf("%s listen is down", ln.bind)
			}
			err := ln.handle(fn, fd)
			if err == nil {
				continue
			}

			if net.ErrClosed.Error() == err.Error() {
				return nil
			}
			ln.xEnv.Errorf("%s handle error %v", fd.Addr().String(), err)
		}
	}
}

func (ln *Listener) multipleH(fn accept) error {
	me := catch.New()
	n := len(ln.fd)
	for i := 0; i < n; i++ {
		go func(k int) {
			fd := ln.fd[k]
			defer fd.Close()
			err := ln.singleH(fn, fd)
			me.Try(fd.Addr().String(), err)
		}(i)
	}

	<-ln.ctx.Done()
	ln.xEnv.Errorf("%s multiple handle exit", ln.bind.String())
	return me.Wrap()
}

type accept func(context.Context, net.Conn, context.CancelFunc) error

func (ln *Listener) OnAccept(fn accept) error {

	n := len(ln.fd)
	if n < 1 {
		return fmt.Errorf("not found ative listen fd")
	}

	if n == 1 {
		return ln.singleH(fn, ln.fd[0])
	} else {
		return ln.multipleH(fn)
	}
}

func (ln *Listener) Close() error {
	if ln == nil {
		return nil
	}

	ln.stop()
	ln.shutdown()
	me := catch.New()
	for _, fd := range ln.fd {
		me.Try(fd.Addr().String(), fd.Close())
	}
	ln.fd = nil
	return me.Wrap()
}

func (ln *Listener) single() error {
	fd, err := net.Listen(ln.bind.Scheme(), ln.bind.Host())
	if err != nil {
		return err
	}
	ln.fd = []net.Listener{fd}
	return nil
}

//multiple tcp://192.168.0.1/?port=1024,65535&exclude=1,2,3,4
func (ln *Listener) multiple() error {
	ps := ln.bind.Ports()
	n := len(ps)
	if n == 0 {
		return fmt.Errorf("%s not found listen", ln.bind.String())
	}

	for i := 0; i < n; i++ {
		port := ps[i]
		fd, e := net.Listen(ln.bind.Scheme(), fmt.Sprintf("%s:%d", ln.bind.Hostname(), port))
		if e != nil {
			return fmt.Errorf("listen %s://%s:%d error %v", ln.bind.Scheme(), ln.bind.Hostname(), port, e)
			continue
		}
		ln.fd = append(ln.fd, fd)
	}

	if len(ln.fd) == 0 {
		return fmt.Errorf("%s listen fail", ln.bind.String())
	}

	return nil
}

func (ln *Listener) Start() error {
	if ln.bind.Port() != 0 {
		return ln.single()
	}

	return ln.multiple()
}

func Listen(env assert.Environment, bind auxlib.URL) (*Listener, error) {
	ctx, stop := context.WithCancel(context.Background())
	ln := &Listener{
		bind: bind,
		ctx:  ctx,
		stop: stop,
		done: 0,
		xEnv: env,
	}

	return ln, ln.Start()
}
