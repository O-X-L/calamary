package meta

type Proto uint8
type Action uint8

const (
	ActionAccept Action = 0
	ActionDeny   Action = 1

	ProtoL3IP4 Proto = 5
	ProtoL3IP6 Proto = 6

	ProtoL4Tcp Proto = 10
	ProtoL4Udp Proto = 11

	ProtoL5Tls  Proto = 20
	ProtoL5Http Proto = 21
	ProtoL5Dns  Proto = 22
	ProtoL5Ntp  Proto = 23
)
