package archiver

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var Verbose bool

type ZipArchive struct {
	file   *os.File
	writer *zip.Writer
}

func Create(filename string) (*ZipArchive, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return &ZipArchive{
		file:   f,
		writer: zip.NewWriter(f),
	}, nil
}

func (z *ZipArchive) AddFile(relPath, absPath string) error {
	info, err := os.Stat(absPath)
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = relPath
	header.Method = zip.Deflate

	writer, err := z.writer.CreateHeader(header)
	if err != nil {
		return err
	}

	file, err := os.Open(absPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(writer, file)
	return err
}

func (z *ZipArchive) AddDir(relPath, absPath string) error {
	return filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Calculate relative path within the zip
		innerPath := filepath.ToSlash(filepath.Join(relPath, strings.TrimPrefix(path, absPath)))
		return z.AddFile(innerPath, path)
	})
}

func (z *ZipArchive) Close() error {
	if err := z.writer.Close(); err != nil {
		return err
	}
	return z.file.Close()
}

func PackTo(srcDir, targetFile string, includeBase bool) error {
	z, err := Create(targetFile)
	if err != nil {
		return err
	}
	defer z.Close()

	base := ""
	if includeBase {
		base = filepath.Base(srcDir)
	}

	return z.AddDir(base, srcDir)
}

func ExtractTo(archivePath, targetDir string) error {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		path := filepath.Join(targetDir, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}
