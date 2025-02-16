//go:build windows && 386
// +build windows,386

package tun

import (
	"fmt"
	"log"

	"golang.org/x/sys/windows"
)

// Variable to hold the device initialization function
var deviceInit func(name string) (windows.Handle, error)

// Default device initialization
func defaultDeviceInit(name string) (windows.Handle, error) {
	const tapGuid = `\\.\Global\{4D36E972-E325-11CE-BFC1-08002BE10318}`
	h, err := windows.CreateFile(
		windows.StringToUTF16Ptr(tapGuid),
		windows.GENERIC_READ|windows.GENERIC_WRITE,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_SYSTEM|windows.FILE_FLAG_OVERLAPPED,
		0)
	if err != nil {
		return 0, fmt.Errorf("failed to create default tap device: %v", err)
	}
	return h, nil
}

// Add Windows Server 2008 compatibility check
func init() {
	// Set default initialization
	deviceInit = defaultDeviceInit

	if err := checkWindowsVersion(); err != nil {
		log.Printf("windows version compatibility warning: %v", err)
	}
}

func checkWindowsVersion() error {
	ver := windows.RtlGetVersion()
	if ver.MajorVersion < 6 {
		return fmt.Errorf("windows version too old: %d.%d", ver.MajorVersion, ver.MinorVersion)
	}

	// Windows Server 2008 is 6.0
	if ver.MajorVersion == 6 && ver.MinorVersion == 0 {
		log.Printf("running on Windows Server 2008 - enabling compatibility mode")
		// Initialize Server 2008 compatibility
		initServer2008Mode()
	}

	return nil
}

func initServer2008Mode() {
	log.Printf("Initializing Server 2008 compatibility mode")
	// Hook into the device creation process
	deviceInit = server2008DeviceInit
}

// Server 2008 specific device initialization
func server2008DeviceInit(name string) (windows.Handle, error) {
	log.Printf("Using Server 2008 device initialization")
	if err := checkTapAdapter(); err != nil {
		log.Printf("TAP adapter check failed: %v", err)
		return 0, err
	}

	h, err := windows.CreateFile(
		windows.StringToUTF16Ptr(`\\.\Global\{4D36E972-E325-11CE-BFC1-08002BE10318}`),
		windows.GENERIC_READ|windows.GENERIC_WRITE,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_SYSTEM,
		0)
	if err != nil {
		log.Printf("Failed to create tap device: %v", err)
		return 0, fmt.Errorf("failed to create tap device: %v", err)
	}
	log.Printf("Successfully created TAP device")
	return h, nil
}

func checkTapAdapter() error {
	log.Printf("Checking TAP adapter")
	const tapGuid = `\\.\Global\{4D36E972-E325-11CE-BFC1-08002BE10318}`
	h, err := windows.CreateFile(
		windows.StringToUTF16Ptr(tapGuid),
		windows.GENERIC_READ|windows.GENERIC_WRITE,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_SYSTEM,
		0)
	if err != nil {
		log.Printf("TAP adapter not found: %v", err)
		return fmt.Errorf("tap adapter not found - please install TAP-Windows")
	}
	windows.CloseHandle(h)
	log.Printf("TAP adapter check successful")
	return nil
}
