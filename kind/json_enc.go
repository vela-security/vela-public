package kind

import (
	"github.com/vela-security/vela-public/auxlib"
	"github.com/vela-security/vela-public/buffer"
	"github.com/vela-security/vela-public/lua"
	"strconv"
	"time"
)

type JsonEncoder struct {
	cache *buffer.Byte
	//cache   []byte
}

func NewJsonEncoder() *JsonEncoder {
	return &JsonEncoder{buffer.Get()}
}

func NewJson(cache []byte) *JsonEncoder {
	return &JsonEncoder{cache: &buffer.Byte{B: cache}}
}

func (enc *JsonEncoder) Char(ch byte) {
	enc.cache.WriteByte(ch)
}

func (enc *JsonEncoder) WriteByte(ch byte) {
	switch ch {
	case '\\':
		enc.cache.WriteString("\\\\")
	case '\r':
		enc.cache.WriteString("\\r")

	case '\n':
		enc.cache.WriteString("\\n")

	case '\t':
		enc.cache.WriteString("\\t")
	case '"':
		enc.cache.WriteString("\\\"")

	default:
		enc.cache.WriteByte(ch)
	}
}

func (enc *JsonEncoder) WriteString(val string) {
	n := len(val)
	if n == 0 {
		return
	}

	for i := 0; i < n; i++ {
		enc.WriteByte(val[i])
	}
}

func (enc *JsonEncoder) Write(val []byte) {
	n := len(val)
	if n == 0 {
		return
	}

	for i := 0; i < n; i++ {
		enc.WriteByte(val[i])
	}
}

func (enc *JsonEncoder) Key(key string) {
	enc.Char('"')
	enc.WriteString(key)
	enc.Char('"')
	enc.WriteByte(':')
}

func (enc *JsonEncoder) Val(v string) {
	enc.Char('"')
	enc.WriteString(v)
	enc.Char('"')
}

func (enc *JsonEncoder) Insert(v []byte) {
	enc.Char('"')
	enc.Write(v)
	enc.Char('"')
}

func (enc *JsonEncoder) Int(n int) {
	enc.WriteString(strconv.Itoa(n))
}

func (enc *JsonEncoder) Bool(v bool) {
	if v {
		enc.Write(True)
	} else {
		enc.Write(False)
	}
}

func (enc *JsonEncoder) Long(n int64) {
	enc.WriteString(auxlib.ToString(n))
}

func (enc *JsonEncoder) ULong(n uint64) {
	enc.WriteString(auxlib.ToString(n))
}

func (enc *JsonEncoder) KT(key string, t time.Time) {
	enc.Key(key)
	enc.Val(t.String())
	enc.WriteByte(',')
}

//func (enc *JsonEncoder) KV(key , val string) {
//	enc.Key(key)
//	enc.Val(val)
//	enc.WriteByte(',')
//}

func (enc *JsonEncoder) KI(key string, n int) {
	enc.Key(key)
	enc.Int(n)
	enc.WriteByte(',')
}

func (enc *JsonEncoder) KF64(key string, v float64) {
	enc.Key(key)
	enc.WriteString(auxlib.ToString(v))
	enc.WriteByte(',')
}

func (enc *JsonEncoder) KL(key string, n int64) {
	enc.Key(key)
	enc.Long(n)
	enc.WriteByte(',')
}

func (enc *JsonEncoder) KUL(key string, n uint64) {
	enc.Key(key)
	enc.ULong(n)
	enc.WriteByte(',')
}

func (enc *JsonEncoder) Join(key string, v []string) {
	enc.Key(key)

	enc.Arr("")
	for _, item := range v {
		enc.Val(item)
		enc.WriteByte(',')
	}

	enc.End("]")
	enc.WriteByte(',')
}

func (enc *JsonEncoder) Join2(key string, v []interface{}) {
	enc.Key(key)

	enc.Arr("")
	for _, item := range v {
		enc.WriteString(auxlib.ToString(item))
		enc.WriteByte(',')
	}

	enc.End("]")
	enc.WriteByte(',')
}

func (enc *JsonEncoder) kv1(key, v string) {
	enc.Key(key)
	enc.WriteString(v)
	enc.WriteByte(',')
}

func (enc *JsonEncoder) kv2(key, v string) {
	enc.Key(key)
	enc.Val(v)
	enc.WriteByte(',')
}

func (enc *JsonEncoder) KV(key string, s interface{}) {
	switch val := s.(type) {
	case bool:
		enc.kv1(key, strconv.FormatBool(val))
	case float64:
		enc.kv1(key, strconv.FormatFloat(val, 'f', -1, 64))
	case float32:
		enc.kv1(key, strconv.FormatFloat(float64(val), 'f', -1, 64))
	case int:
		enc.kv1(key, strconv.Itoa(val))
	case int64:
		enc.kv1(key, strconv.FormatInt(val, 10))
	case int32:
		enc.kv1(key, strconv.Itoa(int(val)))

	case int16:
		enc.kv1(key, strconv.FormatInt(int64(val), 10))
	case int8:
		enc.kv1(key, strconv.FormatInt(int64(val), 10))
	case uint:
		enc.kv1(key, strconv.FormatUint(uint64(val), 10))
	case uint64:
		enc.kv1(key, strconv.FormatUint(val, 10))
	case uint32:
		enc.kv1(key, strconv.FormatUint(uint64(val), 10))
	case uint16:
		enc.kv1(key, strconv.FormatUint(uint64(val), 10))
	case uint8:
		enc.kv1(key, strconv.FormatUint(uint64(val), 10))

	case string:
		enc.kv2(key, val)

	case lua.LString:
		enc.kv2(key, string(val))
	case lua.LBool:
		enc.kv1(key, strconv.FormatBool(bool(val)))
	case lua.LNilType:
		enc.kv2(key, "nil")

	case lua.LNumber:
		enc.kv1(key, strconv.FormatFloat(float64(val), 'f', -1, 64))
	case lua.LInt:
		enc.kv1(key, strconv.Itoa(int(val)))

	case []string:
		enc.Join(key, val)
	case []byte:
		enc.kv2(key, auxlib.B2S(val))

	case []interface{}:
		enc.Join2(key, val)

	case time.Time:
		if y := val.Year(); y < 0 || y >= 10000 {
			// RFC 3339 is clear that years are 4 digits exactly.
			// See golang.org/issue/4556#c15 for more discussion.

			return
			//return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
		}
		enc.kv2(key, val.Format(time.RFC3339Nano))
	case error:
		enc.kv2(key, val.Error())

	default:
		enc.kv2(key, "")

	}

}

var False = []byte("false")
var True = []byte("true")

func (enc *JsonEncoder) KB(key string, b bool) {
	enc.Key(key)

	if b {
		enc.Write(True)
	} else {
		enc.Write(False)
	}

	enc.WriteByte(',')
}

func (enc *JsonEncoder) False(key string) {
	enc.Key(key)
	enc.Write(False)
	enc.WriteByte(',')
}

func (enc *JsonEncoder) True(key string) {
	enc.Key(key)
	enc.Write(True)
	enc.WriteByte(',')
}

func (enc *JsonEncoder) Tab(name string) {
	if len(name) != 0 {
		enc.Val(name)
		enc.WriteByte(':')
	}

	enc.WriteByte('{')
}

func (enc *JsonEncoder) Arr(name string) {
	if len(name) != 0 {
		enc.Val(name)
		enc.WriteByte(':')
	}
	enc.WriteByte('[')
}

func (enc *JsonEncoder) Raw(key string, val []byte) {
	n := len(val)
	if n == 0 {
		return
	}

	enc.Key(key)
	enc.cache.Write(val)
	enc.WriteByte(',')
}

func (enc *JsonEncoder) End(val string) {
	n := enc.cache.Len()

	if n <= 0 {
		return
	}

	if enc.cache.B[n-1] == ',' {
		enc.cache.B = enc.cache.B[:n-1]
	}

	enc.WriteString(val)
}

func (enc *JsonEncoder) Bytes() []byte {
	return enc.cache.Bytes()
}

func (enc *JsonEncoder) Buffer() *buffer.Byte {
	return enc.cache
}
