package auxlib

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func Output() (func(string, ...interface{}), *os.File) {
	exe, _ := os.Executable()
	path := filepath.Dir(exe) + "/daemon.log"

	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("start ssc open file %v\n", err)
	}

	return func(format string, args ...interface{}) {
		header := fmt.Sprintf("{\"level\":\"ERROR\" , \"ts\":\"%s\" , \"caller\":\"rock/cmd.go:36\" , %s}\n",
			time.Now().Format("2006-01-02 15:04:05"), format)
		if file == nil {
			fmt.Printf(header, args...)
			return
		}
		file.WriteString(fmt.Sprintf(header, args...))
	}, file

}
