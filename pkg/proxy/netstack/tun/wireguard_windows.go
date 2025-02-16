//go:build windows && (amd64 || 386)

package tun

import (
	"fmt"
	"log"

	"golang.org/x/sys/windows"
)

const offset = 0

type TunInterface struct {
	// Add required fields
	fd windows.Handle
}

// Add Windows Server 2008 compatibility check
func init() {
	// Check Windows version
	if err := checkWindowsVersion(); err != nil {
		log.Printf("Windows version compatibility warning: %v", err)
	}
}

func checkWindowsVersion() error {
	ver := windows.RtlGetVersion()
	if ver.MajorVersion < 6 {
		return fmt.Errorf("Windows version too old: %d.%d", ver.MajorVersion, ver.MinorVersion)
	}

	// Windows Server 2008 is 6.0
	if ver.MajorVersion == 6 && ver.MinorVersion == 0 {
		log.Printf("Running on Windows Server 2008 - some features may have limited functionality")
	}

	return nil
}

// New creates a new TUN interface
func New(tunName string) (*TunInterface, error) {
	// Check for TAP adapter presence
	if err := checkTapAdapter(); err != nil {
		return nil, fmt.Errorf("TAP adapter check failed: %v", err)
	}

	// Create basic TUN interface
	tun := &TunInterface{}
	return tun, nil
}

func checkTapAdapter() error {
	// For Windows Server 2008, use a different method to check TAP adapter
	const tapGuid = `\\.\Global\{4D36E972-E325-11CE-BFC1-08002BE10318}`
	h, err := windows.CreateFile(
		windows.StringToUTF16Ptr(tapGuid),
		windows.GENERIC_READ|windows.GENERIC_WRITE,
		0,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_SYSTEM|windows.FILE_FLAG_OVERLAPPED,
		0)
	if err != nil {
		return fmt.Errorf("TAP adapter not found - please install TAP-Windows")
	}
	windows.CloseHandle(h)
	return nil
}
