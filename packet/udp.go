package packet

import (
	"bytes"
	"encoding/binary"
)

type UDPHeader struct {
	Source      uint16
	Destination uint16
	Length      uint16
	Checksum    uint16
	Data        []byte
}

func NewUDPHeader(data []byte) *UDPHeader {
	var udp UDPHeader
	r := bytes.NewReader(data)
	binary.Read(r, binary.BigEndian, &udp.Source)
	binary.Read(r, binary.BigEndian, &udp.Destination)
	binary.Read(r, binary.BigEndian, &udp.Length)
	binary.Read(r, binary.BigEndian, &udp.Checksum)
	udp.Data = data[8:udp.Length]
	return &udp
}
