package main

import (
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"github.com/michep/snap-plugin-collector-crsctl/crsctl"
)

func main() {
	plugin.StartCollector(crsctl.NewCollector(), crsctl.PluginName, crsctl.PluginVersion, plugin.RoutingStrategy(plugin.StickyRouter))
}
