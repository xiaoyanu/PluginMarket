package controller

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const uploadCleanupGracePeriod = time.Hour

var uploadReferencePattern = regexp.MustCompile(`(?i)/uploads/[^\s"'<>?#)\]]+`)

type unusedUploadFile struct {
	Path         string    `json:"path"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
}

type uploadScanResult struct {
	UnusedCount int                `json:"unusedCount"`
	UnusedBytes int64              `json:"unusedBytes"`
	Files       []unusedUploadFile `json:"files"`
}

type uploadCleanupResult struct {
	CleanedCount int   `json:"cleanedCount"`
	CleanedBytes int64 `json:"cleanedBytes"`
	SkippedCount int   `json:"skippedCount"`
}

func extractUploadReferences(values []string) map[string]struct{} {
	result := make(map[string]struct{})
	for _, value := range values {
		for _, match := range uploadReferencePattern.FindAllString(value, -1) {
			result[normalizeUploadReference(match)] = struct{}{}
		}
	}
	return result
}

func normalizeUploadReference(value string) string {
	value = strings.ReplaceAll(value, `\`, "/")
	start := strings.Index(strings.ToLower(value), "/uploads/")
	if start < 0 {
		return ""
	}
	return value[start:]
}

func scanUnusedUploads(root string, used map[string]struct{}, now time.Time, gracePeriod time.Duration) (uploadScanResult, error) {
	result := uploadScanResult{Files: make([]unusedUploadFile, 0)}
	cutoff := now.Add(-gracePeriod)

	err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() || !info.Mode().IsRegular() || info.ModTime().After(cutoff) {
			return nil
		}

		reference, err := uploadReferenceForFile(root, path)
		if err != nil {
			return err
		}
		if _, exists := used[reference]; exists {
			return nil
		}

		result.Files = append(result.Files, unusedUploadFile{
			Path:         reference,
			Size:         info.Size(),
			LastModified: info.ModTime(),
		})
		result.UnusedCount++
		result.UnusedBytes += info.Size()
		return nil
	})
	if errors.Is(err, os.ErrNotExist) {
		return result, nil
	}
	if err != nil {
		return uploadScanResult{}, err
	}

	sort.Slice(result.Files, func(i, j int) bool { return result.Files[i].Path < result.Files[j].Path })
	return result, nil
}

func deleteUnusedUploads(
	root string,
	now time.Time,
	gracePeriod time.Duration,
	loadReferences func() (map[string]struct{}, error),
) (uploadCleanupResult, error) {
	lock, err := lockUploadRoot(root)
	if err != nil {
		return uploadCleanupResult{}, err
	}
	defer unlockUploadRoot(lock)

	used, err := loadReferences()
	if err != nil {
		return uploadCleanupResult{}, err
	}
	candidates, err := scanUnusedUploads(root, used, now, gracePeriod)
	if err != nil {
		return uploadCleanupResult{}, err
	}

	result := uploadCleanupResult{}
	for _, candidate := range candidates.Files {
		// 每个文件删除前重新查询数据库，避免扫描后又被文章或其他记录引用。
		currentReferences, err := loadReferences()
		if err != nil {
			return result, err
		}
		if _, usedNow := currentReferences[candidate.Path]; usedNow {
			result.SkippedCount++
			continue
		}

		path, err := localPathForUpload(root, candidate.Path)
		if err != nil {
			return result, err
		}
		info, err := os.Lstat(path)
		if os.IsNotExist(err) {
			result.SkippedCount++
			continue
		}
		if err != nil {
			return result, err
		}
		if !info.Mode().IsRegular() || info.ModTime().After(now.Add(-gracePeriod)) {
			result.SkippedCount++
			continue
		}
		if err := os.Remove(path); err != nil {
			return result, err
		}
		result.CleanedCount++
		result.CleanedBytes += info.Size()
	}

	return result, nil
}

func uploadReferenceForFile(root, path string) (string, error) {
	relative, err := filepath.Rel(root, path)
	if err != nil {
		return "", err
	}
	if relative == "." || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("上传文件路径超出 uploads 目录: %s", path)
	}
	return "/uploads/" + filepath.ToSlash(relative), nil
}

func localPathForUpload(root, reference string) (string, error) {
	const prefix = "/uploads/"
	if !strings.HasPrefix(reference, prefix) {
		return "", fmt.Errorf("无效上传路径: %s", reference)
	}
	path := filepath.Join(root, filepath.FromSlash(strings.TrimPrefix(reference, prefix)))
	relative, err := filepath.Rel(root, path)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("上传文件路径超出 uploads 目录: %s", reference)
	}
	return path, nil
}
