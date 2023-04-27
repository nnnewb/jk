package utils

import (
	"github.com/juju/errors"
	"golang.org/x/mod/modfile"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// ResolveFullPackagePath 从给定路径向上搜索go模块根目录，返回完整包名
func ResolveFullPackagePath(origin, dirpath string) (string, error) {
	filename := filepath.Join(dirpath, "go.mod")
	if _, err := os.Stat(filename); err == nil {
		content, err := os.ReadFile(filename)
		if err != nil {
			return "", errors.Trace(err)
		}
		file, err := modfile.Parse(filename, content, nil)
		if err != nil {
			return "", errors.Trace(err)
		}

		relativePath, err := filepath.Rel(dirpath, origin)
		if err != nil {
			return "", errors.Trace(err)
		}

		return path.Join(file.Module.Mod.Path, strings.ReplaceAll(relativePath, "\\", "/")), nil
	} else {
		dirpath = filepath.Dir(dirpath)
		return ResolveFullPackagePath(origin, dirpath)
	}
}
