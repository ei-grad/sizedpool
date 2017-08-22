package sizedpool

import "sync"

type Bytes struct {
	pools []*sync.Pool
	mu    sync.RWMutex
}

func (p *Bytes) Get(size int) []byte {
	var (
		i = 0
		j = 4096
	)
	for {
		if i >= len(p.pools) {
			p.mu.Lock()
			if i >= len(p.pools) {
				p.pools = append(p.pools, &sync.Pool{})
			}
			p.mu.Unlock()
		}
		if j >= size {
			p.mu.RLock()
			ret := p.pools[i].Get()
			p.mu.RUnlock()
			if ret == nil {
				ret = make([]byte, j)
			}
			return ret.([]byte)[:size]
		}
		i++
		j *= 2
	}
}

func (p *Bytes) Put(b []byte) {
	var (
		i = 0
		j = 4096
	)
	for {
		if i >= len(p.pools) {
			p.mu.Lock()
			if i >= len(p.pools) {
				p.pools = append(p.pools, &sync.Pool{})
			}
			p.mu.Unlock()
		}
		if j == cap(b) {
			p.mu.RLock()
			p.pools[i].Put(b)
			p.mu.RUnlock()
			return
		}
		if j > cap(b) {
			return
		}
		i++
		j *= 2
	}
}
