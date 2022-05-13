package packet

import (
	"github.com/miekg/dns"
)

func Dns(udp *UDPHeader) (dns.Msg, error) {
	var msg dns.Msg
	err := msg.Unpack(udp.Data)
	return msg, err
}
