package files

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type SavedFile struct {
	OriginalName string
	SavedName    string
	RelativePath string
	AbsolutePath string
	Size         int64
	Ext          string
}

type LocalStore struct {
	BaseDir string
}

func (s LocalStore) Save(file multipart.File, header *multipart.FileHeader, subdir string) (*SavedFile, error) {
	defer file.Close()

	extension := strings.ToLower(filepath.Ext(header.Filename))
	savedName := uuid.NewString() + extension
	relativePath := filepath.Join(subdir, savedName)
	absolutePath := filepath.Join(s.BaseDir, relativePath)

	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		return nil, fmt.Errorf("mkdir upload dir: %w", err)
	}

	dst, err := os.Create(absolutePath)
	if err != nil {
		return nil, fmt.Errorf("create upload file: %w", err)
	}
	defer dst.Close()

	size, err := io.Copy(dst, file)
	if err != nil {
		return nil, fmt.Errorf("save upload file: %w", err)
	}

	return &SavedFile{
		OriginalName: header.Filename,
		SavedName:    savedName,
		RelativePath: filepath.ToSlash(relativePath),
		AbsolutePath: absolutePath,
		Size:         size,
		Ext:          extension,
	}, nil
}

func (s LocalStore) Delete(relativePath string) error {
	target := filepath.Join(s.BaseDir, relativePath)
	if err := os.Remove(target); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete upload file: %w", err)
	}
	return nil
}
