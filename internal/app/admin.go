package app

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

// AdminManager handles admin/root permission checks and elevation
type AdminManager struct{}

// NewAdminManager creates a new AdminManager instance
func NewAdminManager() *AdminManager {
	return &AdminManager{}
}

// IsAdmin checks if the current process is running with admin/root privileges
func (a *AdminManager) IsAdmin() (bool, error) {
	switch runtime.GOOS {
	case "windows":
		return a.isWindowsAdmin()
	case "darwin", "linux":
		return a.isUnixAdmin()
	default:
		return false, errors.New("unsupported operating system")
	}
}

// RequestAdmin requests admin/root privileges for the current process
func (a *AdminManager) RequestAdmin() error {
	switch runtime.GOOS {
	case "windows":
		return a.requestWindowsAdmin()
	case "darwin", "linux":
		return a.requestUnixAdmin()
	default:
		return errors.New("unsupported operating system")
	}
}

// isWindowsAdmin checks for admin privileges on Windows
func (a *AdminManager) isWindowsAdmin() (bool, error) {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false, nil
	}
	return true, nil
}

// isUnixAdmin checks for root privileges on Unix-like systems
func (a *AdminManager) isUnixAdmin() (bool, error) {
	return os.Geteuid() == 0, nil
}

// requestWindowsAdmin attempts to elevate privileges on Windows using runas
func (a *AdminManager) requestWindowsAdmin() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	// Prepare the command to run with elevated privileges
	cmd := exec.Command("runas", "/user:Administrator", exe)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    false,
		CreationFlags: syscall.CREATE_NEW_CONSOLE,
	}

	// Start the new process
	if err := cmd.Start(); err != nil {
		return errors.New("failed to elevate privileges: " + err.Error())
	}

	// Exit the current process
	os.Exit(0)
	return nil
}

// requestUnixAdmin attempts to elevate privileges on Unix-like systems using sudo
func (a *AdminManager) requestUnixAdmin() error {
	// Check if we're already root
	if os.Geteuid() == 0 {
		return nil
	}

	// Get the path to the current executable
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	// Prepare the sudo command
	cmd := exec.Command("sudo", exe, os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the new process
	if err := cmd.Start(); err != nil {
		return errors.New("failed to elevate privileges: " + err.Error())
	}

	// Exit the current process
	os.Exit(0)
	return nil
}

// EnsureAdmin ensures the application is running with admin/root privileges
func (a *AdminManager) EnsureAdmin() error {
	isAdmin, err := a.IsAdmin()
	if err != nil {
		return err
	}

	if !isAdmin {
		return a.RequestAdmin()
	}

	return nil
}
