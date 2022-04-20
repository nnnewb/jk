package plugins

import (
	"github.com/nnnewb/jk/pkg/gen/gencore"
	"github.com/nnnewb/jk/pkg/gen/plugins/gensvc"
	"github.com/nnnewb/jk/pkg/gen/plugins/transports/genrpc"
)

func init() {
	gencore.RegisterPlugin("svc", gencore.GenFuncPlugin(gensvc.GenerateEndpoint))
	gencore.RegisterPlugin("netrpc", gencore.GenFuncPlugin(genrpc.GenerateBindings))
}
