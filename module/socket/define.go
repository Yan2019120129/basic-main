package socket

const (
	EventMessageName      = "message"      //	消息
	EventOnlineStatusName = "onlineStatus" //	在线状态
)

// OnlineStatus 在线状态
type OnlineStatus struct {
	UserId int64 `json:"userId"`
	Online bool  `json:"online"`
}

// OutputMessage 输出消息
type OutputMessage struct {
	Event string      `json:"event"` //	事件名称
	Data  interface{} `json:"data"`  //	事件数据
}
