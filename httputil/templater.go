package httputil

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"text/template"
	"time"
)

type templater struct {
	fs       http.FileSystem
	includes map[string]bool
	data     interface{}
}

var _ http.FileSystem = (*templater)(nil)

func (t *templater) Open(path string) (http.File, error) {
	// FIXME add a cache

	if t.includes[path] {
		f, err := t.fs.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		var buf = new(bytes.Buffer)
		if _, err = io.Copy(buf, f); err != nil {
			return nil, err
		}

		tmpl, err := template.New(path).Parse(string(buf.Bytes()))
		if err != nil {
			return nil, err
		}

		buf.Reset()

		if err = tmpl.Execute(buf, t.data); err != nil {
			return nil, err
		}

		info, err := f.Stat()
		if err != nil {
			return nil, err
		}

		return newBfile(buf.Bytes(), info), nil
	}

	return t.fs.Open(path)
}

func NewTemplater(fs http.FileSystem, includes []string, data interface{}) http.FileSystem {
	t := &templater{fs, make(map[string]bool, len(includes)), data}

	for _, path := range includes {
		t.includes[path] = true
	}

	return t
}

type bfile struct {
	bytes.Reader
	info os.FileInfo
}

func (*bfile) Close() error {
	return nil
}

func (*bfile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *bfile) Stat() (os.FileInfo, error) {
	return f.info, nil
}

var _ http.File = (*bfile)(nil)

type bfileInfo struct {
	info os.FileInfo
	size int64
}

var _ os.FileInfo = bfileInfo{}

func (i bfileInfo) Name() string {
	return i.info.Name()
}

func (i bfileInfo) Size() int64 {
	return i.size
}

func (i bfileInfo) Mode() os.FileMode {
	return i.info.Mode()
}

func (i bfileInfo) ModTime() time.Time {
	return i.info.ModTime()
}

func (i bfileInfo) IsDir() bool {
	return i.info.IsDir()
}

func (i bfileInfo) Sys() interface{} {
	return i.info.Sys()
}

func newBfile(b []byte, info os.FileInfo) http.File {
	return &bfile{
		*bytes.NewReader(b),
		bfileInfo{
			info,
			int64(len(b)),
		},
	}
}
