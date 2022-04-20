package gencore

type GenFunc func(data *PluginData) error

type GenPlugin interface {
	Generate(data *PluginData) error
}

type funcPlugin struct {
	f GenFunc
}

func (p funcPlugin) Generate(data *PluginData) error {
	return p.f(data)
}

func GenFuncPlugin(f GenFunc) funcPlugin {
	return funcPlugin{
		f: f,
	}
}
