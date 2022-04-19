package driver

import "io"

type GeneratedFile struct {
	Writer       io.Writer
	Filepath     string
	RelativePath string
}

func (g *GeneratedFile) P(text ...string) *GeneratedFile {
	for _, t := range text {
		g.Writer.Write([]byte(t))
	}
	return g
}
