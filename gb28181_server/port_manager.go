package gb28181_server

import "errors"

var ErrNoAvailablePort = errors.New("no available port")

type PortManager struct {
	recycle chan uint16
	start   uint16
	max     uint16
	pos     uint16
	Valid   bool
}

func (pm *PortManager) Init(start, end uint16) {
	pm.start = start
	pm.pos = start - 1
	pm.max = end
	if pm.pos > 0 && pm.max > pm.pos {
		pm.Valid = true
		pm.recycle = make(chan uint16, pm.Range())
	}
}

func (pm *PortManager) Range() uint16 {
	return pm.max - pm.pos
}

func (pm *PortManager) Recycle(port uint16) error {
	select {
	case pm.recycle <- port:
		return nil
	default:
		return ErrNoAvailablePort
	}
}

func (pm *PortManager) GetPort() (port uint16, err error) {
	select {
	case port = <-pm.recycle:
		return
	default:
		if pm.Range() > 0 {
			pm.pos++
			port = pm.pos
			return
		} else {
			pm.pos = pm.start - 1
			port = pm.pos
			return
		}
	}
}
