package driver

import (
	"go/token"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type GenerateRequest struct {
	FileSet       *token.FileSet
	Pkg           *types.Package
	Svc           *types.Named
	SvcName       string
	genFiles      map[string]*strings.Builder
	OutDir        string
	TypesImporter types.Importer
}

func NewServiceGenerateRequest(fst *token.FileSet, pkg *types.Package, svcName string, svc *types.Named, outDir string, i types.Importer) *GenerateRequest {
	return &GenerateRequest{
		FileSet:       fst,
		Pkg:           pkg,
		Svc:           svc,
		SvcName:       svcName,
		genFiles:      map[string]*strings.Builder{},
		OutDir:        outDir,
		TypesImporter: i,
	}
}

func (s *GenerateRequest) GenFile(relativePath string) *GeneratedFile {
	realpath, err := filepath.Abs(filepath.Join(s.OutDir, relativePath))
	if err != nil {
		log.Fatal(err)
	}

	if _, ok := s.genFiles[realpath]; ok {
		log.Fatalf("can not generate file %s twice", relativePath)
	}

	sb := &strings.Builder{}
	s.genFiles[realpath] = sb
	return &GeneratedFile{
		Writer:       sb,
		Filepath:     realpath,
		RelativePath: relativePath,
	}
}

func (s *GenerateRequest) SaveGeneratedFiles() error {
	for path, w := range s.genFiles {
		_, err := os.Stat(filepath.Dir(path))
		if err != nil {
			if os.IsNotExist(err) {
				os.MkdirAll(filepath.Dir(path), 0o755)
			} else {
				return err
			}
		}

		err = os.WriteFile(path, []byte(w.String()), 0o644)
		if err != nil {
			return err
		}
	}
	return nil
}
