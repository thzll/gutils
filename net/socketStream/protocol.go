package socketStream

import "fmt"

const (
	PACKAGE_TAG                  = 0xFFFF
	PACKAGE_MAX_BODY_SIZE        = 60000
	OP_PING                      = 0xf0
	OP_PONG                      = 0xF1
	OP_PROXY_PROTOCOL            = 0x80
	OP_CONNECT                   = 1
	OP_SEND                      = 2
	OP_CLOSE                     = 3
	OP_ACK                       = 4
	OP_CONNECTED                 = 5  //连接服务器成功
	OP_CONNECT_FAILD             = 6  //连接服务器失败
	OP_DISCONNECTED              = 7  //连接断开
	OP_HEART                     = 8  //链接失败
	OP_SOCK_CLOSED               = 9  //sockid已经关闭
	OP_SUBIP_ERROR               = 10 //代理suip配置丢失
	OP_CLEAR_SEND_INDEX          = 11 //服务端关闭
	OP_NEED_INIT_CONNECTION_INFO = 12 //丢失数据包需要重发的
	OP_INIT_CONNECTION_INFO      = 13 //初始化连接信息
	OP_NO_CONNECTION_ID          = 14 //代理反馈无此链接ID
	OP_YES_CONNECTION_ID         = 15 //代理反馈有此链接ID
	OP_FIND_CONNECTION_ID        = 16 //向代理查询SocketId是否存在
	OP_SEND_TIMEOUT              = 17 //代理发送数据超时

	OP_INSTANCE_GET_INFO    = 0x20 //获取实例信息
	OP_INSTANCE_USER_DATA   = 0x21 //用户实例数据用户发送到数据量之类的
	OP_INSTANCE_NODE_STAT   = 0x22 //用户实例数据用户发送到数据量之类的
	OP_INSTANCE_WAIGUA_JSON = 0x23 //外挂窗口Json
	OP_INSTANCE_INFO        = 0x30 //节点返回实例信息
	OP_INSTANCE_SHOW        = 0x32 //通知盾显示信息
	OP_INSTANCE_EXIT        = 0x33 //通知盾退出进程
	OP_INSTANCE_CLEAR       = 0x34 //清除外挂窗口发送标志

	OP_KCP_DATA = 100
)

var opMap map[uint16]string

func init() {
	opMap = make(map[uint16]string)
	opMap[OP_PING] = "OP_PING"
	opMap[OP_PONG] = "OP_PONG"
	opMap[OP_PROXY_PROTOCOL] = "OP_PROXY_PROTOCOL"
	opMap[OP_CONNECT] = "OP_CONNECT"
	opMap[OP_SEND] = "OP_SEND"
	opMap[OP_CLOSE] = "OP_CLOSE"
	opMap[OP_ACK] = "OP_ACK"
	opMap[OP_CONNECTED] = "OP_CONNECTED"
	opMap[OP_CONNECT_FAILD] = "OP_CONNECT_FAILD"
	opMap[OP_DISCONNECTED] = "OP_DISCONNECTED"
	opMap[OP_HEART] = "OP_HEART"
	opMap[OP_SOCK_CLOSED] = "OP_SOCK_CLOSED"
	opMap[OP_SUBIP_ERROR] = "OP_SUBIP_ERROR"
	opMap[OP_CLEAR_SEND_INDEX] = "OP_CLEAR_SEND_INDEX"
	opMap[OP_NEED_INIT_CONNECTION_INFO] = "OP_NEED_INIT_CONNECTION_INFO"
	opMap[OP_INIT_CONNECTION_INFO] = "OP_INIT_CONNECTION_INFO"
	opMap[OP_NO_CONNECTION_ID] = "OP_NO_CONNECTION_ID"
	opMap[OP_YES_CONNECTION_ID] = "OP_YES_CONNECTION_ID"
	opMap[OP_FIND_CONNECTION_ID] = "OP_FIND_CONNECTION_ID"
	opMap[OP_SEND_TIMEOUT] = "OP_SEND_TIMEOUT"
	opMap[OP_INSTANCE_GET_INFO] = "OP_INSTANCE_GET_INFO"
	opMap[OP_INSTANCE_USER_DATA] = "OP_INSTANCE_USER_DATA"
	opMap[OP_INSTANCE_NODE_STAT] = "OP_INSTANCE_NODE_STAT"
	opMap[OP_INSTANCE_WAIGUA_JSON] = "OP_INSTANCE_WAIGUA_JSON"
	opMap[OP_INSTANCE_INFO] = "OP_INSTANCE_INFO"
	opMap[OP_INSTANCE_SHOW] = "OP_INSTANCE_SHOW"
	opMap[OP_INSTANCE_EXIT] = "OP_INSTANCE_EXIT"
	opMap[OP_INSTANCE_CLEAR] = "OP_INSTANCE_CLEAR"
	opMap[OP_KCP_DATA] = "OP_KCP_DATA"
}

func ParseOP(op uint16) string {
	if name, ok := opMap[op]; ok {
		return name
	} else {
		return fmt.Sprintf("uknow %d", op)
	}
}
