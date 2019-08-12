package roomPkg

import (
	"study/basic_server/chatServer/protocol"
)

type RoomConfig struct {
	StartRoomNumber int32
	MaxRoomCount    int32
	MaxUserCount    int32
}

type roomUser struct {
	netSessionIndex    int32
	netSessionUniqueId uint64

	// 다른 유저에게 알려줘야 하는 정보
	RoomUniqueId   uint64
	IDLen          int8
	ID             [protocol.MAX_USER_ID_BYTE_LENGTH]byte
	packetDataSize int16
}
