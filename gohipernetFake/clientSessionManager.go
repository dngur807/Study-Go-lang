package gohipernetFake

import (
	"go.uber.org/zap"
	"sync"
)

type tcpClientSessionManager struct {
	_networkFunctor  SessionNetworkFunctors
	_sessionMap      sync.Map
	_curSessionCount int32 // 멀티 스레드에서 호출된다.
}

func newClientSessionManager(config *NetworkConfig,
	networkFunctor SessionNetworkFunctors) *tcpClientSessionManager {
	sessionMgr := new(tcpClientSessionManager)
	sessionMgr._networkFunctor = networkFunctor
	sessionMgr._sessionMap = sync.Map{}
	return sessionMgr
}

func (sessionMgr *tcpClientSessionManager) addSession(session *TcpSession) bool {
	sessionIndex := session.Index
	sessionUniqueId := session.SeqIndex

	_, result := sessionMgr._findSession(sessionIndex, sessionUniqueId)
	if result {
		return false
	}

	Logger.Info("SessionManager - addSession" , zap.Uint64("unique" , sessionUniqueId))
	sessionMgr._sessionMap.Store(sessionUniqueId , session)

	return true
}

func (sessionMgr *tcpClientSessionManager) removeSession(sessionUniqueId uint64) {
	Logger.Info("SessionManager - removeSession", zap.Uint64("unique" , sessionUniqueId))
	sessionMgr._sessionMap.Delete(sessionUniqueId)
}

func (sessionMgr *tcpClientSessionManager) _findSession(sessionIndex int32,
	sessionUniqueId uint64) (*TcpSession, bool) {
	if session, ok := sessionMgr._sessionMap.Load(sessionUniqueId); ok {
		return session.(*TcpSession), true
	}
	return nil, false
}

func (sessionMgr *tcpClientSessionManager) sendPacket (sessionIndex int32,
	sessionUniqueId uint64,
	sendData []byte) bool {
		session, result := sessionMgr._findSession(sessionIndex , sessionUniqueId)
		if result == false {
			return false
		}

		session.sendPacket(sendData)
		return true
}
