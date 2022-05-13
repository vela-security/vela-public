package assert

import (
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
	JSON(any) error
	SaveFile(string) (string, error)
}

type TnlByEnv interface { //tunnel by env
	TnlName() string
	TnlVersion() string
	TnlIsDown() bool
	TnlSend(Opcode, any) error
	HTTP(string, string, string, io.Reader, http.Header) HTTPResponse
	PostJSON(string, any, any) error
	Stream(string, any) (HTTPStream, error)
	WithTnl(interface{})
}
