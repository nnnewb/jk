package gencore

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

type File struct {
	wr   *bytes.Buffer // file content (in memory)
	path string        // absolute path
}

func NewFile(absolutePath string) *File {
	return &File{
		wr:   &bytes.Buffer{},
		path: absolutePath,
	}
}

func (f *File) WriteToDisk() error {
	dst := filepath.Dir(f.path)
	info, err := os.Stat(dst)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dst, 0o755); err != nil {
				return err
			}
		} else {
			return err
		}
	} else if !info.IsDir() {
		return fmt.Errorf("%s is not dir", dst)
	}

	return os.WriteFile(f.path, f.wr.Bytes(), 0o644)
}

func (f *File) P(content ...string) {
	for _, s := range content {
		_, err := f.wr.Write([]byte(s))
		if err != nil {
			panic(err)
		}
	}
}

func (f *File) Write(p []byte) (n int, err error) {
	return f.wr.Write(p)
}
