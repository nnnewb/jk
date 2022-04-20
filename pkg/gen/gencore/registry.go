package gencore

import "fmt"

var registry map[string]GenPlugin = make(map[string]GenPlugin)

func GetPlugin(name string) GenPlugin {
	if plugin, ok := registry[name]; ok {
		return plugin
	}
	return nil
}

func RegisterPlugin(name string, plugin GenPlugin) {
	if _, exist := registry[name]; exist {
		panic(fmt.Errorf("%s already registered as gen plugin", name))
	}
	registry[name] = plugin
}
