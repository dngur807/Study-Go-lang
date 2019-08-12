package main

import (
	"go.uber.org/zap"
	"strconv"
	"strings"
	"study/basic_server/chatServer/connectedSessions"
	"study/basic_server/chatServer/protocol"
	"study/basic_server/chatServer/roomPkg"
	. "study/basic_server/gohipernetFake"
	"time"
)

type configAppServer struct {
	GameName         string
	RoomMaxCount     int32
	RoomStartNum     int32
	RoomMaxUserCount int32
}

type ChatServer struct {
	ServerIndex int
	IP          string
	Port        int

	PacketChan chan protocol.Packet
	RoomMgr    *roomPkg.RoomManager
}

func createAnsStartServer(netConfig NetworkConfig, appConfig configAppServer) {
	NTELIB_LOG_INFO("Create Server !!!")

	var server ChatServer

	if server.setIPAddress(netConfig.BindAddress) == false {
		NTELIB_LOG_ERROR("fail server address")
		return
	}

	protocol.Init_packet()
	maxUserCount := appConfig.RoomMaxUserCount * appConfig.RoomMaxCount
	connectedSessions.Init(netConfig.MaxSessionCount, maxUserCount)

	server.PacketChan = make(chan protocol.Packet, 256)

	roomConfig := roomPkg.RoomConfig{
		appConfig.RoomStartNum,
		appConfig.RoomMaxCount,
		appConfig.RoomMaxUserCount,
	}
	server.RoomMgr = roomPkg.NewRoomManager(roomConfig)
	go server.PacketProcess_goroutine()


	networkFunctor := SessionNetworkFunctors{}
	networkFunctor.OnConnect = server.OnConnect
	networkFunctor.OnReceive = server.OnReceive
	networkFunctor.OnReceiveBufferData = nil
	networkFunctor.OnClose = server.OnClose
	networkFunctor.PacketTotalSizeFunc = nil
	networkFunctor.PacketHeaderSize = PACKET_HEADER_SIZE
	networkFunctor.IsClientSession = true

	NetLibInitNetwork(PACKET_HEADER_SIZE, PACKET_HEADER_SIZE)
	NetLibStartNetwork(&netConfig, networkFunctor)
	
}

func (server *ChatServer) setIPAddress(ipAddress string) bool {
	results := strings.Split(ipAddress, ":")
	if len(results) != 2 {
		return false
	}
	server.IP = results[0]
	server.Port, _ = strconv.Atoi(results[1])

	NTELIB_LOG_INFO("Server Address", zap.String("IP", server.IP), zap.Int("Port", server.Port))
	return true
}

func (server *ChatServer) OnConnect(sessionIndex int32, sessionUniqueID uint64) {
	NTELIB_LOG_INFO("client OnConnect", zap.Int32("sessionIndex",sessionIndex), zap.Uint64("sessionUniqueId",sessionUniqueID))
	
	connectedSessions.AddSession(sessionIndex, sessionUniqueID)
}

func (server *ChatServer) OnReceive(sessionIndex int32,sessionUniqueID uint64, data []byte) bool {
	NTELIB_LOG_DEBUG("OnReceive", zap.Int32("sessionIndex", sessionIndex),
		zap.Uint64("sessionUniqueID", sessionUniqueID),
		zap.Int("packetSize", len(data)))

	return true
}

func (server *ChatServer) OnClose(sessionIndex int32, sessionUniqueID uint64) {
	NTELIB_LOG_INFO("client OnCloseClientSession" ,zap.Int32("sessionIndex", sessionIndex) , zap.Uint64("sessionUniqueId" , sessionUniqueID))
	server.disConnectClient(sessionIndex , sessionUniqueID)
}

func (server *ChatServer) disConnectClient(sessionIndex int32, sessionUniqueId uint64) {
	// 로그인도 안한 유저라면 그냥 여기서 처리한다.
	// 방 입장을 안한 유저라면 여기서 처리해도 괜찮지만 아래로 넘긴다.
	if connectedSessions.IsLoginUser(sessionIndex) == false {
		NTELIB_LOG_INFO("DisConnectClient - Not Login User", zap.Int32("sessionIndex" , sessionIndex))
		connectedSessions.RemoveSession(sessionIndex, false)
		return
	}

	packet := protocol.Packet {
		sessionIndex ,
		sessionUniqueId,
		protocol.PACKET_ID_SESSION_CLOSE_SYS,
		0,
		nil,
	}

	server.PacketChan <- packet
	NTELIB_LOG_INFO("DisConnectClient - Login User" , zap.Int32("sessionIndex" , sessionIndex))
}


func (server *ChatServer) Stop() {
	NTELIB_LOG_INFO("chatServer Stop !!!")
	
	NetLib_StopServer() // 이 함수가 꼭 제일 먼저 호출 되어야 한다.
	NTELIB_LOG_INFO("chatServer Stop Waitting...")
	time.Sleep(1 * time.Second)
}