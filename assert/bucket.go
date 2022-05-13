package assert

type Bucket interface {
	Store(string, interface{}, int) error   //key , value , expire
	Replace(string, interface{}, int) error //key , value , expire
	Delete(string) error                    //key
	DeleteBucket(string) error              //bucket name
	Get(string) (interface{}, error)        //key
	Incr(string, int, int) (int, error)
	Int(string) int
	Int64(string) int64
	Bool(string) bool
	Push(string, []byte, int64) error
	Value(string) ([]byte, error)
	Names() string
	String() string
}
