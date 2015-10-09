package gotrig

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

/*
#include "skarnet.h"
*/
import "C"

type Listener struct {
	B   <-chan byte
	Err error

	f          *os.File
	fd, fdw    uintptr
	path       string
	shouldQuit chan struct{}
}

// Close stops consuming events, may only be invoked once per *Listener
func (x *Listener) Close() {
	close(x.shouldQuit)
}

func (x *Listener) close() {
	os.Remove(x.path)
	C.fd_close(C.int(x.fdw))
	x.f.Close()
}

func randomName() string {
	return strconv.FormatInt(time.Now().Unix(), 16) + "-" + strconv.FormatInt(rand.Int63(), 16)
}

func subscribe(name string, fifo, fifow uintptr) (*Listener, error) {
	f := os.NewFile(fifo, name)
	ch := make(chan byte)
	ch2 := make(chan byte)
	x := &Listener{
		B:          ch2,
		fd:         fifo,
		fdw:        fifow,
		path:       name,
		shouldQuit: make(chan struct{}),
	}
	go func() {
		buf := make([]byte, 1)
		for {
			n, err := f.Read(buf)
			if n > 0 {
				ch <- buf[0]
			}
			if err != nil {
				if err != io.EOF {
					x.Err = err
				}
				return
			}
		}
	}()
	go func() {
		defer func() {
			close(ch2)
			x.close()
		}()
		for {
			select {
			case <-x.shouldQuit:
				return
			case b := <-ch:
				select {
				case <-x.shouldQuit:
					return
				case ch2 <- b:
				}
			}
		}
	}()
	return x, nil
}

func Listen(path string) (*Listener, error) {
	rname := ".ftrig-" + randomName()
	tmpfifo := filepath.Join(path, rname)
	err := syscall.Mkfifo(tmpfifo, syscall.S_IRUSR|syscall.S_IWUSR|syscall.S_IWGRP|syscall.S_IWOTH)
	if err != nil {
		return nil, err
	}

	fd, fdw, name := C.int(0), C.int(0), filepath.Join(path, rname[1:])
	fd, err = C.open_read(C.CString(tmpfifo))
	if fd == -1 {
		goto err1
	}

	fdw, err = C.open_write(C.CString(tmpfifo))
	if fd == -1 {
		C.fd_close(fd)
		goto err1
	}

	err = os.Rename(tmpfifo, name)
	if err != nil {
		C.fd_close(fdw)
		C.fd_close(fd)
		goto err1
	}

	err = syscall.SetNonblock(int(fd), false)
	if err != nil {
		goto err1
	}

	return subscribe(name, uintptr(fd), uintptr(fdw))
err1:
	os.Remove(tmpfifo)
	if errno, ok := err.(syscall.Errno); ok {
		return nil, errno
	}
	return nil, fmt.Errorf("unknown error opening fifo at %q", tmpfifo)
}
