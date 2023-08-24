package fileutils

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func GetFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}

func SearchArcives(basedir string, findName string) ([]string, error) {
	var files []string
	err := filepath.Walk(basedir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !IsArcive(path) {
			return nil
		}

		matched, err := filepath.Match(findName, info.Name())
		if err != nil {
			return err
		}
		if matched {
			fullpath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			files = append(files, fullpath)
		}
		return nil
	})

	if err != nil {
		return nil, errors.WithStack(err)
	}
	return files, nil
}

func IsArcive(filePath string) bool {
	// 拡張子が以下以外の場合はfalse
	// zip, rar, 7zip, tar
	ext := filepath.Ext(filePath)
	switch ext {
	case ".zip", ".rar", ".7z", ".tar":
		break
	default:
		return false
	}

	// // 拡張子とMIMEタイプが一致しない場合はfalse
	// mimeType, err := GetMimeType(filePath)
	// if err != nil {
	// 	return false
	// }
	// if ext == ".zip" && mimeType != "application/zip" {
	// 	return false
	// }
	// if ext == ".rar" && mimeType != "application/x-rar-compressed" {
	// 	return false
	// }
	// if ext == ".7z" && mimeType != "application/x-7z-compressed" {
	// 	return false
	// }
	// if ext == ".tar" && mimeType != "application/x-tar" {
	// 	return false
	// }
	return true
}

func GetMimeType(filePath string) (string, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", errors.WithStack(err)
	}
	mimeType := http.DetectContentType(bytes)
	return mimeType, nil
}
