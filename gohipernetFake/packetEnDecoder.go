package gohipernetFake

import "encoding/binary"

func packetTotalSize(data []byte) int16 {
	totalsize := binary.LittleEndian.Uint16(data)
	return int16(totalsize)
}