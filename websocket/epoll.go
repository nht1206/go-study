package websocket

import (
	"net"
	"reflect"
	"sync"
	"syscall"

	"golang.org/x/sys/unix"
)

// The Epoll struct represents an Epoll instance.
type Epoll struct {
	fd     int              // File descriptor of the epoll instance.
	conns  map[int]net.Conn // Map to store the relationship between file descriptors and net.Conn.
	locker *sync.RWMutex    // Mutex to synchronize access to the conns map.
}

// NewEPoll creates and returns a new epoll instance.
func NewEPoll() (*Epoll, error) {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}

	return &Epoll{
		fd:     fd,
		conns:  make(map[int]net.Conn),
		locker: &sync.RWMutex{},
	}, nil
}

// Add adds a net.Conn to the epoll instance.
func (e *Epoll) Add(conn net.Conn) error {
	fd := getConnFD(conn)
	err := unix.EpollCtl(
		e.fd,
		syscall.EPOLL_CTL_ADD,
		fd,
		&unix.EpollEvent{Events: unix.POLLIN | unix.POLLHUP, Fd: int32(fd)},
	)
	if err != nil {
		return err
	}

	e.locker.Lock()
	defer e.locker.Unlock()
	e.conns[fd] = conn

	return nil
}

// Remove removes a net.Conn from the epoll instance.
func (e *Epoll) Remove(conn net.Conn) error {
	fd := getConnFD(conn)
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		return err
	}

	e.locker.Lock()
	defer e.locker.Unlock()
	delete(e.conns, fd)

	return err
}

// Wait waits for events on the epoll instance and returns a list of net.Conn.
func (e *Epoll) Wait() ([]net.Conn, error) {
	events := make([]unix.EpollEvent, 100)
	n, err := unix.EpollWait(e.fd, events, 100)
	if err != nil {
		return nil, err
	}

	e.locker.RLock()
	defer e.locker.RUnlock()
	var conns []net.Conn
	for i := 0; i < n; i++ {
		conns = append(conns, e.conns[int(events[i].Fd)])
	}

	return conns, nil
}

// getConnFD extracts the file descriptor from a net.Conn.
func getConnFD(conn net.Conn) int {
	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")

	return int(pfdVal.FieldByName("Sysfd").Int())
}
