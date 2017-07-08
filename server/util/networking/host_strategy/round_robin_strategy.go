package hoststrategy

import (
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

//RoundRobinHostStrategy is a HostStrategy that just picks the next host in order
type RoundRobinHostStrategy struct {
	Hosts         []*url.URL
	disabledHosts []*url.URL
	nextHost      uint32
	timer         *time.Ticker
	sync.RWMutex
}

//GetNextHost gets the next available host
func (r *RoundRobinHostStrategy) GetNextHost() *url.URL {
	v := atomic.AddUint32(&r.nextHost, 1)
	var h *url.URL
	r.RLock()
	if len(r.Hosts) == 0 {
		h = r.disabledHosts[0]
	} else {
		h = r.Hosts[v%uint32(len(r.Hosts))]
	}
	r.RUnlock()
	return h
}

//CheckAlive periodically determines if the targeted hosts are still functioning
func (r *RoundRobinHostStrategy) CheckAlive() {
	r.timer = time.NewTicker(time.Minute)
	go r.checkAlive()
}

func (r *RoundRobinHostStrategy) checkAlive() {
	for _ = range r.timer.C {
		r.checkDisabledHosts()
		r.checkHosts()
	}
}

func (r *RoundRobinHostStrategy) checkHosts() {
	r.RLock()
	hosts := make([]*url.URL, len(r.Hosts))
	copy(hosts, r.Hosts)
	r.RUnlock()
	var iList []int
	for i, h := range hosts {
		if !hostTCPDialer(h.Host) {
			iList = append(iList, i)
		}
	}

	//we had hosts drop out of the active list
	if iList != nil && len(iList) > 0 {
		//We need to loop over all of the hosts that are unavailable
		for i := len(iList) - 1; i > -1; i-- {
			//Add the unavailable host to the disabled list
			r.disabledHosts = append(r.disabledHosts, hosts[iList[i]])
			//Remove the host from the active list
			hosts = append(hosts[:iList[i]], hosts[iList[i]+1:]...)
		}
		r.Lock()
		r.Hosts = hosts
		r.Unlock()
	}
}

func (r *RoundRobinHostStrategy) checkDisabledHosts() {
	var iList []int
	for i, h := range r.disabledHosts {
		if hostTCPDialer(h.Host) {
			iList = append(iList, i)
		}
	}

	//we are adding hosts back into the active list
	if iList != nil && len(iList) > 0 {
		r.RLock()
		hosts := make([]*url.URL, len(r.Hosts))
		copy(hosts, r.Hosts)
		r.RUnlock()
		for i := len(iList) - 1; i > -1; i-- {
			hosts = append(hosts, r.disabledHosts[iList[i]])
			r.disabledHosts = append(r.disabledHosts[:iList[i]], r.disabledHosts[iList[i]+1:]...)
		}
		r.Lock()
		r.Hosts = hosts
		r.Unlock()
	}
}
