package dtus

import (
	"net"
	"sync"
	"time"
)

/*ManagedDtu ... */
type ManagedDtu struct {
	sync.RWMutex
	m map[net.Conn]*Dtu
}

/*NewManagedDtu ...*/
func NewManagedDtu() *ManagedDtu {
	return &ManagedDtu{m: make(map[net.Conn]*Dtu)}
}

/*RegisterDtu ...*/
func (managedDtu *ManagedDtu) RegisterDtu(d *Dtu) {
	managedDtu.Lock()
	defer managedDtu.Unlock()
	d.SetRegisteredTime(time.Now())
	managedDtu.m[d.conn] = d
}

/*RetriveDtu ...*/
func (managedDtu *ManagedDtu) RetriveDtu(conn net.Conn) (d *Dtu, ok bool) {
	managedDtu.RLock()
	defer managedDtu.RUnlock()
	d, ok = managedDtu.m[conn]
	return
}

/*RemoveDtu ...*/
func (managedDtu *ManagedDtu) RemoveDtu(conn net.Conn) {
	managedDtu.Lock()
	defer managedDtu.Unlock()
	conn.Close()
	delete(managedDtu.m, conn)
}

/*Size ...*/
func (managedDtu *ManagedDtu) Size() int {
	managedDtu.RLock()
	defer managedDtu.RUnlock()
	return len(managedDtu.m)
}

/*Clean ...*/
func (managedDtu *ManagedDtu) Clean(d time.Duration) {
	managedDtu.RLock()
	defer managedDtu.RUnlock()
	toRemove := make([]net.Conn, 0, managedDtu.Size()/2)
	for conn, dtu := range managedDtu.m {
		if time.Now().Sub(dtu.LastHeartbeatedTime()) >= d {
			toRemove = append(toRemove, conn)
		}
	}
	for _, conn := range toRemove {
		managedDtu.RemoveDtu(conn)
	}
}

/*Send ...*/
func (managedDtu *ManagedDtu) Send(message []byte) error {
	managedDtu.RLock()
	defer managedDtu.RUnlock()
	for conn := range managedDtu.m {
		_, err := conn.Write(message)
		if err != nil {
			managedDtu.RUnlock()
			return err
		}
		if dtu, ok := managedDtu.RetriveDtu(conn); ok {
			dtu.SetLastSent(message)
		}
	}
	return nil
}

/*Dtus ...*/
func (managedDtu *ManagedDtu) Dtus() []*Dtu {
	managedDtu.RLock()
	defer managedDtu.RUnlock()
	v := make([]*Dtu, 0, len(managedDtu.m))
	for _, value := range managedDtu.m {
		v = append(v, value)
	}
	return v
}
