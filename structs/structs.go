package structs

import "net"

type StreamType int

const (
	StreamTypeSTDIO StreamType = iota
	StreamTypeTCP
	StreamTypeHTTP
)

type StreamDirection int

const (
	StreamDirectionServer StreamDirection = iota
	StreamDirectionClient
)

type StreamStruct struct {
	Type        StreamType
	Destination string
	Port        int
	Connection  net.Conn
	Direction   StreamDirection
}
