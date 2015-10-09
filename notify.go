package gotrig

import (
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
)

/*
#include "skarnet.h"
*/
import "C"

func notifyNosig(path, s string) (int, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return 0, err
	}
	success := 0
	for i := range files {
		if files[i].Mode()&os.ModeNamedPipe == 0 {
			continue
		}
		p := filepath.Join(path, files[i].Name())
		fd, ferr := C.open_write(C.CString(p))
		if fd == -1 {
			if errno, ok := ferr.(syscall.Errno); ok && errno == syscall.ENXIO {
				err = errno
				os.Remove(p)
			}
		} else {
			r, ferr := C.fd_write(fd, C.CString(s), C.uint(len(s)))
			if (r < 0) || uint(r) < uint(len(s)) {
				if errno, ok := ferr.(syscall.Errno); ok && errno == syscall.EPIPE {
					err = errno
					os.Remove(p)
				}
				// what to do if EGAIN ? full fifo -> fix the reader !
				// There's a race condition in extreme cases though ;
				// but it's still better to be nonblocking - the writer
				// shouldn't get in trouble because of a bad reader.
				C.fd_close(fd)
			} else {
				C.fd_close(fd)
				success++
			}
		}
	}
	return success, err
}

func Clean(path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil // yes, we don't care if the dir failed to scan
	}
	for i := range files {
		if files[i].Mode()&os.ModeNamedPipe == 0 {
			continue
		}
		p := filepath.Join(path, files[i].Name())
		fd, ferr := C.open_write(C.CString(p))
		if fd >= 0 {
			C.fd_close(fd)
		} else if errno, ok := ferr.(syscall.Errno); ok && errno == syscall.ENXIO {
			err = os.Remove(p)
		}
	}
	return err
}

// guard signal management
var notifyLock sync.Mutex

func Notify(path, s string) (int, error) {
	notifyLock.Lock()
	signal.Ignore(syscall.SIGPIPE)
	defer func() {
		signal.Reset(syscall.SIGPIPE)
		notifyLock.Unlock()
	}()
	return notifyNosig(path, s)
}
