package archive

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gen2brain/go-unarr"
	"github.com/pkg/errors"
)

func useUnarr(archivePath string, extractPath string) error {
	a, err := unarr.NewArchive(archivePath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer a.Close()

	_, err = a.Extract(extractPath)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func UnRar(archivePath string, extractPath string) error {
	return useUnarr(archivePath, extractPath)
}

func Un7z(archivePath string, extractPath string) error {
	return useUnarr(archivePath, extractPath)
}

func UnTar(archivePath string, extractPath string) error {
	return useUnarr(archivePath, extractPath)
}

func UnZip(archivePath string, extractPath string) error {
	return useUnarr(archivePath, extractPath)
}

func CompressZip(sourcePath string, distPath string) error {

	file, err := os.Create(distPath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	err = filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.WithStack(err)
		}

		headerPath := strings.TrimPrefix(path, sourcePath)
		if len(headerPath) > 0 && (headerPath[0] == '/' || headerPath[0] == '\\') {
			headerPath = headerPath[1:]
		}
		if len(headerPath) == 0 {
			return nil
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return errors.WithStack(err)
		}
		header.Name = headerPath

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})

	return err
}
