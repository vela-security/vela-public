package assert

const (
	CONSOLE Way = iota + 10
	TRANSPORT
	INLINE
	AGAIN
)

type Way uint8

func (way Way) String() string {
	switch way {
	case CONSOLE:
		return "vela-console"
	case TRANSPORT:
		return "tunnel"
	case INLINE:
		return "inline"
	case AGAIN:
		return "again"
	default:
		return "unknown"
	}
}
