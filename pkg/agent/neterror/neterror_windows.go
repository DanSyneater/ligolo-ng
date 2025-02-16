package neterror

import (
	"errors"
	"syscall"
	"golang.org/x/sys/windows"
)

func HostResponded(err error) bool {
	var se syscall.Errno
	if errors.As(err, &se) {
		return errors.Is(se, syscall.WSAECONNRESET) || 
			   errors.Is(se, syscall.WSAECONNABORTED) ||
			   errors.Is(se, windows.ERROR_NOT_SUPPORTED) ||
			   errors.Is(se, windows.ERROR_INVALID_FUNCTION)
	}
	return false
}
