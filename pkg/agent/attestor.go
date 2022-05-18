package agent

import (
	"context"
	"fmt"

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

func NewECSWorkloadAttestor() *Plugin {
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
	metadata, err := p.metadataProvider.getMetadata()
	if err != nil {
		p.logger.Error("error fetching ecs metadata:", err)
		return &workloadattestorv1.AttestResponse{}, status.Error(codes.Unavailable, "metadata failed to fetch")
	}

	return &workloadattestorv1.AttestResponse{
		SelectorValues: p.getSelectorValuesFromMetadata(metadata),
	}, nil
}

func (p *Plugin) getSelectorValuesFromMetadata(meta metadataResponse) []string {
	type selector struct{ key, value string }

	var selectorValues []string

	if meta.IsEmpty {
		return selectorValues
	}

	selectors := []selector{
		{"name", meta.Root.Name},
		{"image", meta.Root.Image},
		{"cluster", meta.Task.Cluster},
		{"arn", meta.Task.TaskARN},
		{"family", meta.Task.Family},
		{"revision", meta.Task.Revision},
		{"availabilityzone", meta.Task.AvailabilityZone},
	}

	// Following the pattern of the existing docker plugin, labels
	// will end up in the form ecs:label:some_label:some_value
	for label, value := range meta.Root.Labels {
		labelValue := fmt.Sprintf("%s:%s", label, value)
		selectors = append(selectors, selector{"label", labelValue})
	}

	for _, s := range selectors {
		selectorValues = append(
			selectorValues,
			fmt.Sprintf("%s:%s", s.key, s.value),
		)
	}

	return selectorValues
}
