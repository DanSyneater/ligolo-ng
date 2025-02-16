//go:build windows && (amd64 || 386)
package tun

import (
	"fmt"
	"golang.org/x/sys/windows"
	"log"
)

const offset = 0

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

func New(tunName string) (*TunInterface, error) {
	// Check for TAP adapter presence
	if err := checkTapAdapter(); err != nil {
		return nil, fmt.Errorf("TAP adapter check failed: %v", err)
	}
	
	// ... rest of the existing code ...
}

func checkTapAdapter() error {
	// Check if TAP Windows Adapter is installed
	const tapGuid = `ROOT\NET\{4D36E972-E325-11CE-BFC1-08002BE10318}`
	_, err := windows.OpenDevice(tapGuid)
	if err != nil {
		return fmt.Errorf("TAP adapter not found - please install TAP-Windows")
	}
	return nil
}
