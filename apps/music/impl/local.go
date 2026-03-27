package impl

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type MusicObj struct {
	libraryDir string
}

func NewMusicObj() *MusicObj {
	libraryDir := defaultLibraryDir()
	if err := os.MkdirAll(libraryDir, 0755); err != nil {
		log.Printf("初始化音乐目录失败，err: %v", err)
	}
	return &MusicObj{libraryDir: libraryDir}
}

func (m *MusicObj) LibraryDir() string {
	return m.libraryDir
}

func (m *MusicObj) GetFromLocalPath(localPath string) *LocalMusic {
	targetDir := strings.TrimSpace(localPath)
	if targetDir == "" {
		targetDir = m.libraryDir
	}

	lm := &LocalMusic{LibraryPath: targetDir}
	pE, err := os.ReadDir(targetDir)
	if err != nil {
		log.Printf("读取本地路径失败，err: %v", err)
		return lm
	}

	sort.Slice(pE, func(i, j int) bool {
		return strings.ToLower(pE[i].Name()) < strings.ToLower(pE[j].Name())
	})

	for _, entry := range pE {
		if entry.IsDir() {
			continue
		}

		track, ok := m.buildTrack(targetDir, entry.Name())
		if !ok {
			continue
		}
		lm.MusicList = append(lm.MusicList, track)
	}
	return lm
}

func (m *MusicObj) SaveUploadedFiles(files []*multipart.FileHeader) (*UploadResult, error) {
	if err := os.MkdirAll(m.libraryDir, 0755); err != nil {
		return nil, fmt.Errorf("create music library: %w", err)
	}

	result := &UploadResult{
		Message: "save success",
		Total:   len(files),
	}

	for _, fileHeader := range files {
		track, err := m.saveUploadedFile(fileHeader)
		if err != nil {
			log.Printf("保存文件失败，file=%s err=%v", fileHeader.Filename, err)
			continue
		}

		result.MusicList = append(result.MusicList, *track)
		result.Saved++
	}

	if result.Saved == 0 {
		return nil, errors.New("all files failed to save")
	}

	return result, nil
}

func (m *MusicObj) ResolveTrackPath(fileName string) (string, error) {
	safeName := sanitizeFileName(fileName)
	if safeName == "" {
		return "", errors.New("invalid file name")
	}

	fullPath := filepath.Join(m.libraryDir, safeName)
	info, err := os.Stat(fullPath)
	if err != nil {
		return "", err
	}
	if info.IsDir() || !isAudioFile(info.Name()) {
		return "", errors.New("music not found")
	}

	return fullPath, nil
}

func (m *MusicObj) GetTrackByFileName(fileName string) (*Objector, error) {
	fullPath, err := m.ResolveTrackPath(fileName)
	if err != nil {
		return nil, err
	}

	track, ok := m.buildTrack(filepath.Dir(fullPath), filepath.Base(fullPath))
	if !ok {
		return nil, errors.New("music not found")
	}
	return &track, nil
}

func (m *MusicObj) saveUploadedFile(fileHeader *multipart.FileHeader) (*Objector, error) {
	fileName := sanitizeFileName(fileHeader.Filename)
	if fileName == "" {
		return nil, errors.New("invalid file name")
	}
	if !isAudioFile(fileName) {
		return nil, errors.New("unsupported audio format")
	}

	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	targetName := uniqueFileName(m.libraryDir, fileName)
	targetPath := filepath.Join(m.libraryDir, targetName)

	dst, err := os.Create(targetPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err = dst.ReadFrom(src); err != nil {
		return nil, err
	}

	track, _ := m.buildTrack(m.libraryDir, targetName)
	return &track, nil
}

func (m *MusicObj) buildTrack(baseDir string, fileName string) (Objector, bool) {
	if !isAudioFile(fileName) {
		return Objector{}, false
	}

	musicName, author := parseTrackInfo(fileName)
	return Objector{
		ID:           buildTrackID(fileName),
		MusicName:    musicName,
		Author:       author,
		From:         "local",
		Path:         filepath.Join(baseDir, fileName),
		FileName:     fileName,
		Format:       strings.TrimPrefix(strings.ToLower(filepath.Ext(fileName)), "."),
		AudioURL:     "/music/local/file/" + url.PathEscape(fileName),
		IsLike:       false,
		IsCollection: false,
	}, true
}

func defaultLibraryDir() string {
	if configuredDir := strings.TrimSpace(os.Getenv("VMUSIC_LIBRARY_DIR")); configuredDir != "" {
		return configuredDir
	}

	workDir, err := os.Getwd()
	if err != nil {
		return filepath.Join("storage", "music")
	}
	return filepath.Join(workDir, "storage", "music")
}

func parseTrackInfo(fileName string) (string, string) {
	baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	for _, delimiter := range []string{" - ", "-", "_"} {
		index := strings.LastIndex(baseName, delimiter)
		if index <= 0 || index >= len(baseName)-len(delimiter) {
			continue
		}

		title := strings.TrimSpace(baseName[:index])
		author := strings.TrimSpace(baseName[index+len(delimiter):])
		if title != "" && author != "" {
			return title, author
		}
	}

	return strings.TrimSpace(baseName), "未知艺术家"
}

func isAudioFile(fileName string) bool {
	switch strings.ToLower(filepath.Ext(fileName)) {
	case ".mp3", ".wav", ".flac", ".m4a", ".ogg":
		return true
	default:
		return false
	}
}

func sanitizeFileName(fileName string) string {
	safeName := strings.TrimSpace(filepath.Base(fileName))
	safeName = strings.ReplaceAll(safeName, "\x00", "")
	safeName = strings.Map(func(r rune) rune {
		switch {
		case r == '/' || r == '\\':
			return -1
		case r < 32:
			return -1
		default:
			return r
		}
	}, safeName)

	return strings.TrimSpace(safeName)
}

func uniqueFileName(baseDir string, originalName string) string {
	targetPath := filepath.Join(baseDir, originalName)
	if _, err := os.Stat(targetPath); errors.Is(err, os.ErrNotExist) {
		return originalName
	}

	extension := filepath.Ext(originalName)
	baseName := strings.TrimSuffix(originalName, extension)
	return fmt.Sprintf("%s-%d%s", baseName, time.Now().UnixNano(), extension)
}

func buildTrackID(fileName string) string {
	sum := sha1.Sum([]byte(strings.ToLower(fileName)))
	return hex.EncodeToString(sum[:8])
}
