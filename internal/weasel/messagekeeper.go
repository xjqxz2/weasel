package weasel

import (
	"sort"
	"sync"
	"time"
)

const MAX_KEEP_CAPACITY = 10240

type device struct {
	deviceId   string
	message    []byte
	pushedTime int64
	writeOff   bool
}

type sortDevices []device

func (p sortDevices) Len() int           { return len(p) }
func (p sortDevices) Less(i, j int) bool { return p[i].pushedTime < p[j].pushedTime }
func (p sortDevices) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type LRUKeeper struct {
	rmu     sync.RWMutex
	devices map[string]*device

	//	当前可保持的最大设备容量，默认1024
	//	可由业务代码继续扩增
	cap  int
	mcap int
}

func NewLRUKeeper() *LRUKeeper {
	return &LRUKeeper{
		devices: make(map[string]*device),
		cap:     0,
		mcap:    MAX_KEEP_CAPACITY,
	}
}

func (p *LRUKeeper) Push(deviceId string, message []byte) {
	p.rmu.Lock()
	defer p.rmu.Unlock()

	mDevice, ok := p.devices[deviceId]
	if !ok {
		mDevice = &device{
			deviceId: deviceId,
		}

		//	执行扩容监测与扩容操作
		if p.cap >= p.mcap {
			//	获取断连的设备表
			offline := p.getOfflineDevices()
			if len(offline) > 0 {
				//	当断设备 > 0时
				//	则替代(释放掉)其中一个最老的设备
				delete(p.devices, offline[0].deviceId)
				p.cap--
			} else {
				//	如果离线值为0，表示所有设备都在线的情况下，则应当对设备池进行扩容
				p.mcap += MAX_KEEP_CAPACITY
			}
		}
	}

	//	记录当前设备的状态值
	mDevice.message = message
	mDevice.writeOff = true
	mDevice.pushedTime = time.Now().Unix()

	//	将最终的状态写入至设备池
	p.devices[deviceId] = mDevice
	p.cap++
}

func (p *LRUKeeper) Message(deviceId string) []byte {
	p.rmu.RLock()
	defer p.rmu.RUnlock()

	if device, ok := p.devices[deviceId]; ok {
		return device.message
	}

	return nil
}

func (p *LRUKeeper) Offline(deviceId string) {
	p.rmu.Lock()
	defer p.rmu.Unlock()

	if device, ok := p.devices[deviceId]; ok {
		device.writeOff = false
	}
}

func (p *LRUKeeper) getOfflineDevices() (devices sortDevices) {
	for _, device := range p.devices {
		if device.writeOff == false {
			devices = append(devices, *device)
		}
	}

	if len(devices) > 0 {
		sort.Sort(devices)
	}

	return devices
}
