package socket

import (
	"sync"

	"github.com/gorilla/websocket"
)

var SocketInstance *Sockeet

type WebsocketControl struct {
	Conn *websocket.Conn
	Mux  sync.RWMutex
}

// WebsocketControlMap 用户socket
type WebsocketControlMap struct {
	sync.RWMutex
	Map map[string]*WebsocketControl
}

// UserIdStringsMap 用户ID 绑定数组字符串
type UserIdStringsMap struct {
	sync.RWMutex
	Map map[int64][]string
}

type Sockeet struct {
	SocketConn *WebsocketControlMap //	用户socket		{connKey: conn}
	BindUserId *UserIdStringsMap    //	绑定用户ID 		{用户ID: connKey}
}

func init() {
	SocketInstance = &Sockeet{
		SocketConn: &WebsocketControlMap{Map: map[string]*WebsocketControl{}},
		BindUserId: &UserIdStringsMap{Map: map[int64][]string{}},
	}
}

// SendChatMessage 发送客服消息
func (_Sockeet *Sockeet) SendChatMessage(userId int64, data interface{}) {
	msg := &OutputMessage{Event: EventMessageName, Data: data}
	_Sockeet.SendWebsocketData(userId, msg)
}

// SendChatOnlineStatus 发送用户在线状态
func (_Sockeet *Sockeet) SendChatOnlineStatus(adminUserId, userId int64, online bool) {
	if adminUserId == userId {
		return
	}
	msg := &OutputMessage{Event: EventOnlineStatusName, Data: &OnlineStatus{UserId: userId, Online: online}}
	_Sockeet.SendWebsocketData(adminUserId, msg)
}

// SendWebsocketData 发送消息
func (_Sockeet *Sockeet) SendWebsocketData(userId int64, data *OutputMessage) {
	_Sockeet.BindUserId.RLock()
	defer _Sockeet.BindUserId.RUnlock()

	if _, ok := _Sockeet.BindUserId.Map[userId]; ok {
		for i := 0; i < len(_Sockeet.BindUserId.Map[userId]); i++ {
			_Sockeet.SendConnKeyData(_Sockeet.BindUserId.Map[userId][i], data)
		}
	}
}

// SendConnKeyData 发送connKey 数据
func (_Sockeet *Sockeet) SendConnKeyData(connKey string, data *OutputMessage) {
	_Sockeet.SocketConn.RLock()
	defer _Sockeet.SocketConn.RUnlock()

	if _, ok := _Sockeet.SocketConn.Map[connKey]; ok {
		_Sockeet.SocketConn.Map[connKey].Mux.Lock()
		_ = _Sockeet.SocketConn.Map[connKey].Conn.WriteJSON(data)
		_Sockeet.SocketConn.Map[connKey].Mux.Unlock()
	}
}

// IsOnline 是否在线
func (_Sockeet *Sockeet) IsOnline(userId int64) bool {
	if _, ok := _Sockeet.BindUserId.Map[userId]; ok {
		if len(_Sockeet.BindUserId.Map[userId]) > 0 {
			return true
		}
	}
	return false
}

// SetWebsocketConn 赋值conn
func (_Sockeet *Sockeet) SetWebsocketConn(userId int64, connKey string, conn *WebsocketControl) {
	_Sockeet.BindUserId.Lock()
	_Sockeet.SocketConn.Lock()

	defer _Sockeet.BindUserId.Unlock()
	defer _Sockeet.SocketConn.Unlock()

	_Sockeet.SocketConn.Map[connKey] = conn
	_Sockeet.BindUserId.Map[userId] = append(_Sockeet.BindUserId.Map[userId], connKey)
}

// DelWebsocketConn 删除conn
func (_Sockeet *Sockeet) DelWebsocketConn(userId int64, connKey string) {
	_Sockeet.BindUserId.Lock()
	_Sockeet.SocketConn.Lock()
	defer _Sockeet.BindUserId.Unlock()
	defer _Sockeet.SocketConn.Unlock()

	delete(_Sockeet.SocketConn.Map, connKey)

	// 删除绑定的UID
	connKeyIndex := _Sockeet.GetRawDataIndexOf(_Sockeet.BindUserId.Map[userId], connKey)
	if connKeyIndex > -1 {
		_Sockeet.BindUserId.Map[userId] = append(_Sockeet.BindUserId.Map[userId][:connKeyIndex], _Sockeet.BindUserId.Map[userId][connKeyIndex+1:]...)
	}
}

// GetRawDataIndexOf 获取原数据Index
func (_Sockeet *Sockeet) GetRawDataIndexOf(rawData []string, connKey string) int {
	indexOf := -1
	for i := 0; i < len(rawData); i++ {
		if rawData[i] == connKey {
			indexOf = i
			break
		}
	}
	return indexOf
}
