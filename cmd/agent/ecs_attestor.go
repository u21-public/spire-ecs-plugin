package main

import (
	"github.com/spiffe/spire-plugin-sdk/pluginmain"
	workloadattestorv1 "github.com/spiffe/spire-plugin-sdk/proto/spire/plugin/agent/workloadattestor/v1"
	"github.com/u21-public/spire-ecs-plugin/pkg/agent"
)

func main() {
	plugin := agent.NewECSWorkloadAttestor()

	pluginmain.Serve(
		workloadattestorv1.WorkloadAttestorPluginServer(plugin),
	)
}
