package weasel

const MAX_KEEP_CAPACITY = 10240

type Keeper interface {
	Push(deviceId string, message []byte)
	Message(deviceId string) []byte
	Offline(deviceId string)
}

// NoMessageKeeper 空消息保持记录器
// 该记录器不记录消息
type NoMessageKeeper struct{}

func NewNoMessageKeeper() *NoMessageKeeper {
	return new(NoMessageKeeper)
}

func (p *NoMessageKeeper) Push(deviceId string, message []byte) {}
func (p *NoMessageKeeper) Message(deviceId string) []byte       { return nil }
func (p *NoMessageKeeper) Offline(deviceId string)              {}
