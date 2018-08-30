package dtus

import (
	"fmt"
	"net"
	"sync"
	"time"
)

/*Dtu ...*/
type Dtu struct {
	sync.RWMutex
	conn                net.Conn
	id                  string
	simNum              string
	lastSent            []byte
	lastReceived        []byte
	registeredTime      time.Time
	lastHeartbeatedTime time.Time
}

/*NewDtu ...*/
func NewDtu() *Dtu {
	return &Dtu{}
}

/*SetConn ...*/
func (dtu *Dtu) SetConn(conn net.Conn) {
	dtu.Lock()
	defer dtu.Unlock()
	dtu.conn = conn
}

/*Conn ...*/
func (dtu *Dtu) Conn() net.Conn {
	dtu.RLock()
	defer dtu.RUnlock()
	return dtu.conn
}

/*SetID ...*/
func (dtu *Dtu) SetID(id string) {
	dtu.Lock()
	defer dtu.Unlock()
	dtu.id = id
}

/*ID ...*/
func (dtu *Dtu) ID() string {
	dtu.RLock()
	defer dtu.RUnlock()
	return dtu.id
}

/*SetSimNum ...*/
func (dtu *Dtu) SetSimNum(num string) {
	dtu.Lock()
	defer dtu.Unlock()
	dtu.simNum = num
}

/*SimNum ...*/
func (dtu *Dtu) SimNum() string {
	dtu.RLock()
	defer dtu.RUnlock()
	return dtu.simNum
}

/*SetLastSent ...*/
func (dtu *Dtu) SetLastSent(sent []byte) {
	dtu.Lock()
	defer dtu.RUnlock()
	dtu.lastSent = sent
}

/*LastSent ...*/
func (dtu *Dtu) LastSent() []byte {
	dtu.RLock()
	defer dtu.RUnlock()
	return dtu.lastSent
}

/*SetLastReceived ...*/
func (dtu *Dtu) SetLastReceived(received []byte) {
	dtu.Lock()
	defer dtu.RUnlock()
	dtu.lastReceived = received
}

/*LastReceived ...*/
func (dtu *Dtu) LastReceived() []byte {
	dtu.RLock()
	defer dtu.RUnlock()
	return dtu.lastReceived
}

/*SetRegisteredTime ...*/
func (dtu *Dtu) SetRegisteredTime(time time.Time) {
	dtu.Lock()
	defer dtu.RUnlock()
	dtu.registeredTime = time
}

/*RegisteredTime ...*/
func (dtu *Dtu) RegisteredTime() time.Time {
	dtu.RLock()
	defer dtu.RUnlock()
	return dtu.registeredTime
}

/*SetLastHeartbeatedTime ...*/
func (dtu *Dtu) SetLastHeartbeatedTime(time time.Time) {
	dtu.Lock()
	defer dtu.Unlock()
	dtu.lastHeartbeatedTime = time
}

/*LastHeartbeatedTime ...*/
func (dtu *Dtu) LastHeartbeatedTime() time.Time {
	dtu.RLock()
	defer dtu.RUnlock()
	return dtu.lastHeartbeatedTime
}

func (dtu *Dtu) String() string {
	return fmt.Sprintf(
		"Dtu{conn: %v, id: %s, simNum: %s, lastSent: %v, lastReceived: %v, registeredTime: %s, lastHeartbeatedTime: %s}",
		dtu.conn,
		dtu.id,
		dtu.simNum,
		dtu.lastSent,
		dtu.lastReceived,
		dtu.registeredTime.Format(time.RFC3339),
		dtu.lastHeartbeatedTime.Format(time.RFC3339))
}
