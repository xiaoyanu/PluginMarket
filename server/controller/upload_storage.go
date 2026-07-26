package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var uploadDedupMu sync.Mutex

// saveDeduplicatedUpload 按文件内容去重保存。相同内容即使上传到不同用途目录，也复用已有文件路径。
func saveDeduplicatedUpload(root, subDir, ext string, src io.Reader, expectedSize int64) (string, error) {
	uploadDedupMu.Lock()
	defer uploadDedupMu.Unlock()

	root = filepath.Clean(root)
	dir, err := safeUploadSubdirectory(root, subDir)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(root, 0755); err != nil {
		return "", err
	}
	lock, err := lockUploadRoot(root)
	if err != nil {
		return "", err
	}
	defer unlockUploadRoot(lock)

	temp, err := os.CreateTemp(root, ".upload-*")
	if err != nil {
		return "", err
	}
	tempPath := temp.Name()
	defer os.Remove(tempPath)

	hash := sha256.New()
	written, copyErr := io.Copy(io.MultiWriter(temp, hash), src)
	closeErr := temp.Close()
	if copyErr != nil {
		return "", copyErr
	}
	if closeErr != nil {
		return "", closeErr
	}
	if expectedSize >= 0 && written != expectedSize {
		return "", fmt.Errorf("上传文件不完整")
	}

	digest := hash.Sum(nil)
	existing, err := findUploadByContent(root, tempPath, written, digest)
	if err != nil {
		return "", err
	}
	if existing != "" {
		if err := os.Chtimes(existing, time.Now(), time.Now()); err != nil {
			return "", err
		}
		return uploadURL(root, existing)
	}

	ext = strings.ToLower(ext)
	filename := hex.EncodeToString(digest) + ext
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	destination := filepath.Join(dir, filename)
	if err := os.Link(tempPath, destination); err != nil {
		if !os.IsExist(err) {
			return "", err
		}
		matches, matchErr := pathMatchesDigestAndSize(destination, digest, written)
		if matchErr != nil {
			return "", matchErr
		}
		if !matches {
			return "", fmt.Errorf("上传文件目标路径内容冲突")
		}
	}
	if err := os.Chtimes(destination, time.Now(), time.Now()); err != nil {
		return "", err
	}
	return uploadURL(root, destination)
}

func safeUploadSubdirectory(root, subDir string) (string, error) {
	if subDir == "" || filepath.IsAbs(subDir) {
		return "", fmt.Errorf("上传子目录无效")
	}
	dir := filepath.Join(root, filepath.Clean(subDir))
	relative, err := filepath.Rel(root, dir)
	if err != nil || relative == "." || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("上传子目录超出 uploads 目录")
	}
	return dir, nil
}

func pathMatchesDigestAndSize(path string, digest []byte, size int64) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if !info.Mode().IsRegular() || info.Size() != size {
		return false, nil
	}
	return fileMatchesDigest(path, digest)
}

func findUploadByContent(root, tempPath string, size int64, digest []byte) (string, error) {
	var found string
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if found != "" || path == tempPath || !entry.Type().IsRegular() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.Size() != size {
			return nil
		}
		matches, err := fileMatchesDigest(path, digest)
		if err != nil {
			return err
		}
		if matches {
			found = path
		}
		return nil
	})
	return found, err
}

func fileMatchesDigest(path string, expected []byte) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return false, err
	}
	actual := hash.Sum(nil)
	if len(actual) != len(expected) {
		return false, nil
	}
	for i := range actual {
		if actual[i] != expected[i] {
			return false, nil
		}
	}
	return true, nil
}

func uploadURL(root, path string) (string, error) {
	relative, err := filepath.Rel(root, path)
	if err != nil || relative == "." || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("上传文件路径无效")
	}
	return "/uploads/" + filepath.ToSlash(relative), nil
}
