package gb28181_server

import (
	"fmt"
	"math/rand"
	"strconv"
)

// InviteOptions 拉流参数
type InviteOptions struct {
	Start       int
	End         int
	dump        string
	ssrc        string
	SSRC        uint32
	MediaPort   uint16
	StreamPath  string
	recyclePort func(port uint16) (err error)
}

// IsLive 直播流
func (opt *InviteOptions) IsLive() bool {
	return opt.Start == 0 || opt.End == 0
}

// IsRecord 回放流
func (opt *InviteOptions) IsRecord() bool {
	return !opt.IsLive()
}

func (opt *InviteOptions) CreateSSRC() {
	ssrc := make([]byte, 10)
	if opt.IsLive() {
		ssrc[0] = '0'
	} else {
		ssrc[0] = '1'
	}

	copy(ssrc[1:6], globalGB28181Config.Serial[3:8])
	randNum := 1000 + rand.Intn(8999)
	copy(ssrc[6:], strconv.Itoa(randNum))

	opt.ssrc = string(ssrc)
	_SSRC, _ := strconv.ParseInt(opt.ssrc, 10, 0)
	opt.SSRC = uint32(_SSRC)
}

func (opt *InviteOptions) String() string {
	return fmt.Sprintf("t=%d %d", opt.Start, opt.End)
}
