package roomPkg

import (
	"study/basic_server/chatServer/protocol"
	"sync"
)

/**
sync.Pool은 일종의 메모리 풀이라고 볼 수 있다
자원을 풀에 넣었다가, 필요할때 다시 꺼내 쓰는 것이다.
 */
type baseRoom struct {
	_index               int32
	_number              int32 // 채널 고유 번호
	_config              RoomConfig
	_curUserCount        int32
	_roomUserUniqueIdSeq uint64
	_userPool            *sync.Pool

	// 자료구조를 배열로 바꾸는 것이 좋음
	_userSessionUniqueIdMap map[uint64]*roomUser // range 순회 시 복사 비용 발생해서 포인터 값을 사용한다.

	_funcPackIdlist []int16
	_funclist       []func(*roomUser, protocol.Packet) int16

	enterUserNotify func(int64, int32)
	leaveUserNotify func(int64)
}

func (room *baseRoom) initialize(index int32, config RoomConfig) {
	room._initialize(index, config)
	room._initUserPool()
}

func (room *baseRoom) _initialize(index int32, config RoomConfig) {
	room._number = config.StartRoomNumber + index
	room._index = index
	room._config = config
}

func (room *baseRoom) _initUserPool() {
	room._userPool = &sync.Pool{
		New: func() interface{} {
			user := new(roomUser)
			return user
		},
	}
}
func (room *baseRoom) settingPacketFunction() {
	maxFuncListCount := 16
	room._funclist = make([]func(*roomUser, protocol.Packet) int16 , 0 , maxFuncListCount)
	room._funcPackIdlist = make([]int16, 0 , maxFuncListCount)

	room._addPacketFunction(protocol.PACKET_ID_ROOM_ENTER_REQ, room._packetProcess_EnterUser)
	room._addPacketFunction(protocol.PACKET_ID_ROOM_LEAVE_REQ, room._packetProcess_LeaveUser)
}

func (room *baseRoom) _addPacketFunction(packetID int16 , packetFunc func(*roomUser, protocol.Packet) int16) {
	room._funclist = append(room._funclist, packetFunc)
	room._funcPackIdlist = append(room._funcPackIdlist, packetID)
}

