package assert

import (
	opcode "github.com/vela-security/vela-opcode"
	"io"
	"net/http"
)

type HTTPStream interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
}

type HTTPResponse interface {
	Close()
	JSON(interface{}) error
	SaveFile(string) (string, error)
}

type TnlByEnv interface { //tunnel by env
	TnlName() string
	TnlVersion() string
	TnlIsDown() bool
	TnlSend(opcode.Opcode, interface{}) error
	DoHTTP(*http.Request) HTTPResponse
	HTTP(string, string, string, io.Reader, http.Header) HTTPResponse
	PostJSON(string, interface{}, interface{}) error
	Stream(string, interface{}) (HTTPStream, error)
	WithTnl(interface{})
}
