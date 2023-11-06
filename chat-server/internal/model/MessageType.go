package model

type MessageType int

const (
	MsgTypeUndefined = MessageType(0)
	MsgTypeSystem    = MessageType(1)
	MsgTypeUser      = MessageType(2)
)
