package templates

import (
	"bytes"
	"embed"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/cockroachdb/errors"

	"gogs.io/gogs/internal/osx"
)

//go:embed *.tmpl **/*
var files embed.FS

type TemplateFile struct {
	name string
	data []byte
	ext  string
}

func (f *TemplateFile) Name() string { return f.name }
func (f *TemplateFile) Data() []byte { return f.data }
func (f *TemplateFile) Ext() string  { return f.ext }

// fileSystem implements a generic TemplateFileSystem interface.
type fileSystem struct {
	files []*TemplateFile
}

func (fs *fileSystem) ListFiles() []*TemplateFile {
	return fs.files
}

func (fs *fileSystem) Get(name string) (io.Reader, error) {
	for i := range fs.files {
		if fs.files[i].Name()+fs.files[i].Ext() == name {
			return bytes.NewReader(fs.files[i].Data()), nil
		}
	}
	return nil, errors.Newf("file %q not found", name)
}

func mustNames(fsys fs.FS) []string {
	var names []string
	walkDirFunc := func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			names = append(names, path)
		}
		return nil
	}
	if err := fs.WalkDir(fsys, ".", walkDirFunc); err != nil {
		panic("assetNames failure: " + err.Error())
	}
	return names
}

// NewTemplateFileSystem returns a file system instance for embedded assets.
// The argument "dir" can be used to serve subset of embedded assets. Template file
// found under the "customDir" on disk has higher precedence over embedded assets.
func NewTemplateFileSystem(dir, customDir string) *fileSystem {
	if dir != "" && !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	var err error
	var tmplFiles []*TemplateFile
	names := mustNames(files)
	for _, name := range names {
		if !strings.HasPrefix(name, dir) {
			continue
		}
		// Check if corresponding custom file exists
		var data []byte
		fpath := path.Join(customDir, name)
		if osx.IsFile(fpath) {
			data, err = os.ReadFile(fpath)
		} else {
			data, err = files.ReadFile(name)
		}
		if err != nil {
			panic(err)
		}
		name = strings.TrimPrefix(name, dir)
		ext := path.Ext(name)
		name = strings.TrimSuffix(name, ext)
		tmplFiles = append(tmplFiles, &TemplateFile{
			name: name,
			data: data,
			ext:  ext,
		})
	}
	return &fileSystem{files: tmplFiles}
}
