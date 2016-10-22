package master

type Host struct {
	host    string
	workers map[string]*Slave
}

func (h *Host) addSlave(port string, s *Slave) {
	h.workers[port] = s
}

func (h *Host) delSlave(port string) {
	delete(h.workers, port)
}

func (h *Host) GetWorker() *Slave {
	if len(h.workers) == 0 {
		return nil
	}
	var w *Slave
	for _, v := range h.workers {
		if v.handle.IfBusy() == -1 { // not busy
			w = v
			break
		}
	}
	if w != nil {
		return w
	}
	var max int64
	for _, v := range h.workers {
		max = v.handle.IfBusy()
		w = v
		break
	}
	for _, v := range h.workers {
		if v.handle.IfBusy() > max {
			max = v.handle.IfBusy()
			w = v
		}
	}
	return w
}
