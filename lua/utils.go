package lua

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func intMin(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func intMax(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func defaultFormat(v interface{}, f fmt.State, c rune) {
	buf := make([]string, 0, 10)
	buf = append(buf, "%")
	for i := 0; i < 128; i++ {
		if f.Flag(i) {
			buf = append(buf, string(rune(i)))
		}
	}

	if w, ok := f.Width(); ok {
		buf = append(buf, strconv.Itoa(w))
	}
	if p, ok := f.Precision(); ok {
		buf = append(buf, "."+strconv.Itoa(p))
	}
	buf = append(buf, string(c))
	format := strings.Join(buf, "")
	fmt.Fprintf(f, format, v)
}

type flagScanner struct {
	flag       byte
	start      string
	end        string
	buf        []byte
	str        string
	Length     int
	Pos        int
	HasFlag    bool
	ChangeFlag bool
}

func newFlagScanner(flag byte, start, end, str string) *flagScanner {
	return &flagScanner{flag, start, end, make([]byte, 0, len(str)), str, len(str), 0, false, false}
}

func (fs *flagScanner) AppendString(str string) { fs.buf = append(fs.buf, str...) }

func (fs *flagScanner) AppendChar(ch byte) { fs.buf = append(fs.buf, ch) }

func (fs *flagScanner) String() string { return string(fs.buf) }

func (fs *flagScanner) Next() (byte, bool) {
	c := byte('\000')
	fs.ChangeFlag = false
	if fs.Pos == fs.Length {
		if fs.HasFlag {
			fs.AppendString(fs.end)
		}
		return c, true
	} else {
		c = fs.str[fs.Pos]
		if c == fs.flag {
			if fs.Pos < (fs.Length-1) && fs.str[fs.Pos+1] == fs.flag {
				fs.HasFlag = false
				fs.AppendChar(fs.flag)
				fs.Pos += 2
				return fs.Next()
			} else if fs.Pos != fs.Length-1 {
				if fs.HasFlag {
					fs.AppendString(fs.end)
				}
				fs.AppendString(fs.start)
				fs.ChangeFlag = true
				fs.HasFlag = true
			}
		}
	}
	fs.Pos++
	return c, false
}

var cDateFlagToGo = map[byte]string{
	'a': "mon", 'A': "Monday", 'b': "Jan", 'B': "January", 'c': "02 Jan 06 15:04 MST", 'd': "02",
	'F': "2006-01-02", 'H': "15", 'I': "03", 'm': "01", 'M': "04", 'p': "PM", 'P': "pm", 'S': "05",
	'x': "15/04/05", 'X': "15:04:05", 'y': "06", 'Y': "2006", 'z': "-0700", 'Z': "MST"}

func strftime(t time.Time, cfmt string) string {
	sc := newFlagScanner('%', "", "", cfmt)
	for c, eos := sc.Next(); !eos; c, eos = sc.Next() {
		if !sc.ChangeFlag {
			if sc.HasFlag {
				if v, ok := cDateFlagToGo[c]; ok {
					sc.AppendString(t.Format(v))
				} else {
					switch c {
					case 'w':
						sc.AppendString(fmt.Sprint(int(t.Weekday())))
					default:
						sc.AppendChar('%')
						sc.AppendChar(c)
					}
				}
				sc.HasFlag = false
			} else {
				sc.AppendChar(c)
			}
		}
	}

	return sc.String()
}

func isInteger(v LNumber) bool {
	return float64(v) == float64(int64(v))
	//_, frac := math.Modf(float64(v))
	//return frac == 0.0
}

func isArrayKey(v LNumber) bool {
	return isInteger(v) && v < LNumber(int((^uint(0))>>1)) && v > LNumber(0) && v < LNumber(MaxArrayIndex)
}

func parseNumber(number string) (LNumber, error) {
	var value LNumber
	number = strings.Trim(number, " \t\n")
	if v, err := strconv.ParseInt(number, 0, LNumberBit); err != nil {
		if v2, err2 := strconv.ParseFloat(number, LNumberBit); err2 != nil {
			return LNumber(0), err2
		} else {
			value = LNumber(v2)
		}
	} else {
		value = LNumber(v)
	}
	return value, nil
}

func popenArgs(arg string) (string, []string) {
	cmd := "/bin/sh"
	args := []string{"-c"}
	if LuaOS == "windows" {
		cmd = "C:\\Windows\\system32\\cmd.exe"
		args = []string{"/c"}
	}
	args = append(args, arg)
	return cmd, args
}

func isGoroutineSafe(lv LValue) bool {
	switch v := lv.(type) {
	case *LFunction, *LUserData, *LState:
		return false
	case *LTable:
		return v.Metatable == LNil
	default:
		return true
	}
}

func readBufioSize(reader *bufio.Reader, size int64) ([]byte, error, bool) {
	result := []byte{}
	read := int64(0)
	var err error
	var n int
	for read != size {
		buf := make([]byte, size-read)
		n, err = reader.Read(buf)
		if err != nil {
			break
		}
		read += int64(n)
		result = append(result, buf[:n]...)
	}
	e := err
	if e != nil && e == io.EOF {
		e = nil
	}

	return result, e, len(result) == 0 && err == io.EOF
}

func readBufioLine(reader *bufio.Reader) ([]byte, error, bool) {
	result := []byte{}
	var buf []byte
	var err error
	var isprefix bool = true
	for isprefix {
		buf, isprefix, err = reader.ReadLine()
		if err != nil {
			break
		}
		result = append(result, buf...)
	}
	e := err
	if e != nil && e == io.EOF {
		e = nil
	}

	return result, e, len(result) == 0 && err == io.EOF
}

func int2Fb(val int) int {
	e := 0
	x := val
	for x >= 16 {
		x = (x + 1) >> 1
		e++
	}
	if x < 8 {
		return x
	}
	return ((e + 1) << 3) | (x - 8)
}

func strCmp(s1, s2 string) int {
	len1 := len(s1)
	len2 := len(s2)
	for i := 0; ; i++ {
		c1 := -1
		if i < len1 {
			c1 = int(s1[i])
		}
		c2 := -1
		if i != len2 {
			c2 = int(s2[i])
		}
		switch {
		case c1 < c2:
			return -1
		case c1 > c2:
			return +1
		case c1 < 0:
			return 0
		}
	}
}

func unsafeFastStringToReadOnlyBytes(s string) (bs []byte) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&bs))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return
}

var (
	InvalidFormat = errors.New("invalid format")
	InvalidIP     = errors.New("invalid ip addr")
	InvalidPort   = errors.New("expect check socket err: port <1 or port > 65535")
)

func CheckSocket(v string) error {
	s := strings.Split(v, ":")
	if len(s) != 2 {
		return InvalidFormat
	}

	if net.ParseIP(s[0]) == nil {
		return InvalidIP
	}

	port, err := strconv.Atoi(s[1])
	if err != nil {
		return err
	}
	if port < 1 || port > 65535 {
		return InvalidPort
	}
	return nil
}

func IsNotFound(err error) bool {
	if err.Error() == "not found" {
		return true
	}
	return false
}

func IsNull(v []byte) bool {
	if len(v) == 0 {
		return true
	}
	return false
}

func CheckInt(L *LState, lv LValue) int {
	if intv, ok := lv.(LNumber); ok {
		return int(intv)
	}
	L.RaiseError("must be int , got %s", lv.Type().String())
	return 0
}

func CheckIntOrDefault(L *LState, lv LValue, d int) int {
	if intv, ok := lv.(LNumber); ok {
		return int(intv)
	}
	return d
}

func CheckInt64(L *LState, lv LValue) int64 {
	if intv, ok := lv.(LNumber); ok {
		return int64(intv)
	}
	L.RaiseError("must be int64 , got %s", lv.Type().String())
	return 0
}

func CheckNumber(L *LState, lv LValue) LNumber {
	if lv, ok := lv.(LNumber); ok {
		return lv
	}
	L.RaiseError("must be LNumber , got %s", lv.Type().String())
	return 0
}

func CheckString(L *LState, lv LValue) string {
	if lv, ok := lv.(LString); ok {
		return string(lv)
	} else if LVCanConvToString(lv) {
		return LVAsString(lv)
	}
	return ""
}

func CheckBool(L *LState, lv LValue) bool {
	if lv, ok := lv.(LBool); ok {
		return bool(lv)
	}

	L.RaiseError("must be bool , got %s", lv.Type().String())
	return false
}

func CheckTable(L *LState, lv LValue) *LTable {
	if lv, ok := lv.(*LTable); ok {
		return lv
	}
	L.RaiseError("must be LTable, got %s", lv.Type().String())
	return nil
}

func CheckFunction(L *LState, lv LValue) *LFunction {
	if lv, ok := lv.(*LFunction); ok {
		return lv
	}
	L.RaiseError("must be Function, got %s", lv.Type().String())
	return nil
}

func CheckUserData(L *LState, lv LValue) *LUserData {
	if lv, ok := lv.(*LUserData); ok {
		return lv
	}
	L.RaiseError("must be UserData, got %s", lv.Type().String())
	return nil
}

func CheckLightUserData(L *LState, lv LValue) *ProcData {
	if lv, ok := lv.(*ProcData); ok {
		return lv
	}
	L.RaiseError("must be ProcData, got %s", lv.Type().String())
	return nil
}

func CheckThread(L *LState, lv LValue) *LState {
	if lv, ok := lv.(*LState); ok {
		return lv
	}
	L.RaiseError("must be thread , got %s", lv.Type().String())
	return nil
}

func CheckType(L *LState, lv LValue, typ LValueType) {
	if lv.Type() != typ {
		L.RaiseError("must be %", typ.String(), lv.Type().String())
	}
}

func S2B(s string) (b []byte) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return
}

func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func S2L(s string) LString {
	return LString(s)
}

func B2L(b []byte) LString {
	return *(*LString)(unsafe.Pointer(&b))
}

func JsonMarshal(L *LState, i interface{}) LString {
	data, e := json.Marshal(i)
	if e != nil {
		L.RaiseError("%v", e)
		return LSNull
	}
	return *(*LString)(unsafe.Pointer(&data))
}

func FileSuffix(path string) string {
	suffix := filepath.Ext(path)
	return strings.TrimSuffix(path, suffix)
}

func ToLValue(L *LState, v interface{}) LValue {
	switch converted := v.(type) {
	case bool:
		return LBool(converted)
	case float64:
		return LNumber(converted)
	case float32:
		return LNumber(converted)
	case int:
		return LNumber(converted)
	case int32:
		return LNumber(converted)
	case int64:
		return LNumber(converted)
	case string:
		return LString(converted)
	case time.Time:
		tt := float64(converted.UTC().UnixNano()) / float64(time.Second)
		return LNumber(tt)
	case nil:
		return LNil
	default:
		return nil
	}
}
