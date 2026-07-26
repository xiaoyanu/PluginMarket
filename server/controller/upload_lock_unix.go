//go:build !windows

package controller

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

func lockUploadRoot(root string) (*os.File, error) {
	root = filepath.Clean(root)
	lock, err := os.OpenFile(root+".dedup.lock", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	if err := unix.Flock(int(lock.Fd()), unix.LOCK_EX); err != nil {
		_ = lock.Close()
		return nil, err
	}
	return lock, nil
}

func unlockUploadRoot(lock *os.File) {
	_ = unix.Flock(int(lock.Fd()), unix.LOCK_UN)
	_ = lock.Close()
}
