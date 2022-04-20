package gencore

import "path/filepath"

type PluginData struct {
	Request *GenRequest
	files   map[string]*File
}

func NewPluginData(req *GenRequest) *PluginData {
	return &PluginData{
		Request: req,
		files:   map[string]*File{},
	}
}

func (p *PluginData) GetOrCreateFile(path string) *File {
	realpath := filepath.Join(p.Request.GetServiceLocalPath(), path)
	if _, ok := p.files[realpath]; ok {
		return p.files[realpath]
	}

	f := NewFile(realpath)
	p.files[realpath] = f

	return f
}

func (p *PluginData) WriteToDisk() error {
	for _, f := range p.files {
		if err := f.WriteToDisk(); err != nil {
			return err
		}
	}
	return nil
}
