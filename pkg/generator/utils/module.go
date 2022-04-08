package utils

import (
	"encoding/json"
	"errors"
	"log"
	"os/exec"
	"path/filepath"
)

var (
	goEnv  _GoEnv
	goList _GoList
)

type _GoList struct {
	Path      string `json:"Path"`
	Main      bool   `json:"Main"`
	Dir       string `json:"Dir"`
	GoMod     string `json:"GoMod"`
	GoVersion string `json:"GoVersion"`
}

type _GoEnv struct {
	Ar           string `json:"AR"`
	Cc           string `json:"CC"`
	CgoCflags    string `json:"CGO_CFLAGS"`
	CgoCppflags  string `json:"CGO_CPPFLAGS"`
	CgoCxxflags  string `json:"CGO_CXXFLAGS"`
	CgoEnabled   string `json:"CGO_ENABLED"`
	CgoFflags    string `json:"CGO_FFLAGS"`
	CgoLdflags   string `json:"CGO_LDFLAGS"`
	Cxx          string `json:"CXX"`
	Gccgo        string `json:"GCCGO"`
	Go111Module  string `json:"GO111MODULE"`
	Goamd64      string `json:"GOAMD64"`
	Goarch       string `json:"GOARCH"`
	Gobin        string `json:"GOBIN"`
	Gocache      string `json:"GOCACHE"`
	Goenv        string `json:"GOENV"`
	Goexe        string `json:"GOEXE"`
	Goexperiment string `json:"GOEXPERIMENT"`
	Goflags      string `json:"GOFLAGS"`
	Gogccflags   string `json:"GOGCCFLAGS"`
	Gohostarch   string `json:"GOHOSTARCH"`
	Gohostos     string `json:"GOHOSTOS"`
	Goinsecure   string `json:"GOINSECURE"`
	Gomod        string `json:"GOMOD"`
	Gomodcache   string `json:"GOMODCACHE"`
	Gonoproxy    string `json:"GONOPROXY"`
	Gonosumdb    string `json:"GONOSUMDB"`
	Goos         string `json:"GOOS"`
	Gopath       string `json:"GOPATH"`
	Goprivate    string `json:"GOPRIVATE"`
	Goproxy      string `json:"GOPROXY"`
	Goroot       string `json:"GOROOT"`
	Gosumdb      string `json:"GOSUMDB"`
	Gotmpdir     string `json:"GOTMPDIR"`
	Gotooldir    string `json:"GOTOOLDIR"`
	Govcs        string `json:"GOVCS"`
	Goversion    string `json:"GOVERSION"`
	Gowork       string `json:"GOWORK"`
	PkgConfig    string `json:"PKG_CONFIG"`
}

func getGoList() *_GoList {
	if goList.Dir != "" {
		return &goList
	}

	cmd := exec.Command("go", "list", "-m", "-json")
	r, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	if err := json.NewDecoder(r).Decode(&goList); err != nil {
		log.Fatal(err)
	}

	return &goList
}

func getGoEnv() *_GoEnv {
	if goEnv.Goversion != "" {
		return &goEnv
	}

	cmd := exec.Command("go", "env", "-json")
	r, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	if err := json.NewDecoder(r).Decode(&goEnv); err != nil {
		log.Fatal(err)
	}

	return &goEnv
}

func GetGoModImportPath() string {
	l := getGoList()
	if l.Path == "" {
		panic(errors.New("current working directory is not go module!"))
	}
	return l.Path
}

func GetGoModFolderPath() string {
	l := getGoList()
	if l.Dir == "" {
		panic(errors.New("current working directory is not go module!"))
	}
	return filepath.Dir(l.Dir)
}
