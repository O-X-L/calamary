package meta

type Proto uint8
type Action uint8
type OptBool uint8
type Match uint8

const (
	ActionAccept Action = 1
	ActionDeny   Action = 2

	ProtoNone  Proto = 0
	ProtoL3IP4 Proto = 5
	ProtoL3IP6 Proto = 6

	ProtoL4Tcp Proto = 10
	ProtoL4Udp Proto = 11

	ProtoL5Tls  Proto = 20
	ProtoL5Http Proto = 21
	ProtoL5Dns  Proto = 22
	ProtoL5Ntp  Proto = 23

	OptBoolNone  OptBool = 0
	OptBoolTrue  OptBool = 1
	OptBoolFalse OptBool = 2

	MatchNeutral  Match = 0
	MatchPositive Match = 1
	MatchNegative Match = 2
)
