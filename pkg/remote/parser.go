package remote

import (
	"encoding/binary"

	netpoll "github.com/cloudwego/netpoll"
)

func parseHeader(reader netpoll.Reader) (messageType MessageType, length int, seqID uint32, err error) {
	p, err := reader.Peek(10)
	if err != nil {
		return
	}
	length = int(binary.BigEndian.Uint32(p[0:4]))
	version := p[4]
	err = versionCheck(Version(version))
	if err != nil {
		return
	}
	messageType = MessageType(p[5])
	seqID = binary.BigEndian.Uint32(p[6:])

	return
}

func versionCheck(version Version) error {
	return nil
}
