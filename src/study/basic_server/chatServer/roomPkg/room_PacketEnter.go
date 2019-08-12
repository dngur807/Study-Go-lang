package roomPkg

import (
	"study/basic_server/chatServer/connectedSessions"
	"study/basic_server/chatServer/protocol"
	. "study/basic_server/gohipernetFake"
)

func (room *baseRoom) _packetProcess_EnterUser(inValidUser *roomUser, packet protocol.Packet) int16 {
	curTime := NetLib_GetCurrentUnixTime()
	sessionIndex := packet.UserSessionIndex
	sessionUniqueId := packet.UserSessionUniqueId
	
	NTELIB_LOG_INFO("[[Room _packetProcess_EnterUser]]")
	
	var requestPacket protocol.RoomEnterReqPacket
	(&requestPacket).Decoding(packet.Data)

	userID , ok := connectedSessions.GetUserID(sessionIndex)
	if ok == false {
		_sendRoomEnterResult(sessionIndex, sessionUniqueId, 0, 0, protocol.ERROR_CODE_ENTER_ROOM_INVALID_USER_ID)
	}
	return 0
}

func _sendRoomEnterResult(sessionIndex int32, sessionUniqueId uint64, roomNumber int32, userUniqueId uint64, result int16) {
	response := protocol.RoomEnterResPacket{
		result ,
		roomNumber,
		userUniqueId,
	}

	sendPacket, _ := response.EncodingPacket()

}