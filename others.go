package main

func (l *Lvault) IsHostExists(host string) bool {
	l.L.RLock()
	defer l.L.RUnlock()
	if _, ok := l.others[host]; ok {
		return true
	} else {
		return false
	}
}

func (l *Lvault) AddHost(host string) {
	l.L.Lock()
	defer l.L.Unlock()
	if _, ok := l.others[host]; !ok {
		l.others[host] = ""
		l.Others = append(l.Others, host)
	}
}

func (l *Lvault) GetOthers() []string {
	l.L.RLock()
	defer l.L.RUnlock()
	ret := make([]string, 0, len(l.Others))
	ret = append(ret, l.Others...)
	return ret
}
