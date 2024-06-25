package connection

import (
	"fmt"
	"log"
	"net"
	"sync"
	"syscall"

	"golang.org/x/sys/unix"
)

type EventLoop struct {
	kqFD  int
	conns map[int]net.Conn
	mu    sync.RWMutex
}

func NewEventLoop() (*EventLoop, error) {
	kq, err := unix.Kqueue()
	if err != nil {
		return nil, fmt.Errorf("cannot create kqueue: %v", err)
	}

	return &EventLoop{
		kqFD:  kq,
		conns: make(map[int]net.Conn, 0),
	}, nil
}

func (el *EventLoop) Add(conn net.Conn) error {
	fd, err := getFD(conn)

	if err != nil {
		return err
	}

	event := unix.Kevent_t{
		//nolint:gosec // fd is safe
		Ident:  uint64(fd),
		Filter: unix.EVFILT_READ,
		Flags:  unix.EV_ADD | unix.EV_ENABLE,
	}

	_, err = unix.Kevent(el.kqFD, []unix.Kevent_t{event}, nil, nil)
	if err != nil {
		return fmt.Errorf("add kqueue failed: %v", err)
	}

	el.mu.Lock()
	el.conns[fd] = conn
	el.mu.Unlock()

	return nil
}

func (el *EventLoop) Start(onMessage func(conn net.Conn)) {
	events := make([]unix.Kevent_t, 128)
	for {
		n, err := unix.Kevent(el.kqFD, nil, events, nil)
		if err != nil {
			if err == unix.EINTR {
				continue
			}
			log.Printf("Lỗi Kqueue Kevent: %v", err)
			break
		}
		for i := 0; i < n; i++ {
			//nolint:gosec // Ident is safe
			fd := int(events[i].Ident)

			el.mu.RLock()
			conn, exists := el.conns[fd]
			el.mu.RUnlock()
			if exists {
				onMessage(conn)
			}
		}
	}
}

func getFD(conn net.Conn) (int, error) {
	sc, ok := conn.(syscall.Conn)
	if !ok {
		return 0, fmt.Errorf("cannot parse to syscall.Conn")
	}

	rc, err := sc.SyscallConn()

	if err != nil {
		return 0, err
	}

	var fd int
	err = rc.Control(func(f uintptr) {
		fd = int(f)
	})

	return fd, err
}
