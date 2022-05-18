package agent

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/spiffe/spire-plugin-sdk/pluginsdk"
	workloadattestorv1 "github.com/spiffe/spire-plugin-sdk/proto/spire/plugin/agent/workloadattestor/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ pluginsdk.NeedsLogger = &Plugin{} // Compile time check that ensures we implement the logger interface correctly
)

type Plugin struct {
	workloadattestorv1.UnimplementedWorkloadAttestorServer

	logger           hclog.Logger
	metadataProvider metadataProvider
}

func New() *Plugin {
	return &Plugin{
		metadataProvider: &ecsMetadataProvider{},
	}
}

func (p *Plugin) SetLogger(logger hclog.Logger) {
	p.logger = logger
}

func (p *Plugin) SetMetadataProvider(provider metadataProvider) {
	p.metadataProvider = provider
}

func (p *Plugin) Attest(context context.Context, _ *workloadattestorv1.AttestRequest) (*workloadattestorv1.AttestResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Not implemented")
}
