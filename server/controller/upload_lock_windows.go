package controller

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

func lockUploadRoot(root string) (*os.File, error) {
	root = filepath.Clean(root)
	lock, err := os.OpenFile(root+".dedup.lock", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	overlapped := new(windows.Overlapped)
	if err := windows.LockFileEx(windows.Handle(lock.Fd()), windows.LOCKFILE_EXCLUSIVE_LOCK, 0, 1, 0, overlapped); err != nil {
		_ = lock.Close()
		return nil, err
	}
	return lock, nil
}

func unlockUploadRoot(lock *os.File) {
	overlapped := new(windows.Overlapped)
	_ = windows.UnlockFileEx(windows.Handle(lock.Fd()), 0, 1, 0, overlapped)
	_ = lock.Close()
}
