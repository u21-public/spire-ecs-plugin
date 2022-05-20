package agent

import (
	"context"
	"errors"
	"sort"
	"testing"

	v1 "github.com/spiffe/spire-plugin-sdk/proto/spire/plugin/agent/workloadattestor/v1"
)

var staticMetaStruct metadataResponse = metadataResponse{
	Root: metadataRootResponse{
		Name:  "application",
		Image: "myregistry.mydomain.com/image:tag",
		Labels: map[string]string{
			"environment": "development",
			"team":        "somedevteam",
		},
	},
	Task: metadataTaskResponse{
		Cluster:          "mycluster",
		TaskARN:          "aws:arn:someaccount:someregion:sometask",
		Family:           "myservice",
		Revision:         "1",
		AvailabilityZone: "us-east-1d",
	},
}

type staticMetadataProvider struct{}

func (m *staticMetadataProvider) getMetadata() (metadataResponse, error) {
	return staticMetaStruct, nil
}

type errorMetadataProvider struct{}

func (m *errorMetadataProvider) getMetadata() (metadataResponse, error) {
	return metadataResponse{IsEmpty: true}, errors.New("something went wrong!")
}

func TestSelectorGeneration(t *testing.T) {
	plugin := NewECSWorkloadAttestor()
	plugin.SetMetadataProvider(&staticMetadataProvider{})

	attestResponse, err := plugin.Attest(
		context.Background(), &v1.AttestRequest{},
	)

	if err != nil {
		t.Log(err)
		t.Error("received an error from static metadata response")
	}

	selectorsGot := attestResponse.GetSelectorValues()
	selectorsWant := []string{
		"name:application",
		"image:myregistry.mydomain.com/image:tag",
		"label:environment:development",
		"label:team:somedevteam",
		"cluster:mycluster",
		"arn:aws:arn:someaccount:someregion:sometask",
		"family:myservice",
		"revision:1",
		"availabilityzone:us-east-1d",
	}

	sort.Strings(selectorsGot)
	sort.Strings(selectorsWant)

	for i := range selectorsWant {
		if selectorsGot[i] != selectorsWant[i] {
			t.Errorf(
				"selector mismatch, got: %s want: %s",
				selectorsGot[i], selectorsWant[i],
			)
		}
	}
}

func TestErrorSelectors(t *testing.T) {
	plugin := NewECSWorkloadAttestor()
	plugin.SetMetadataProvider(&errorMetadataProvider{})

	attestResponse, err := plugin.Attest(
		context.Background(), &v1.AttestRequest{},
	)

	if len(attestResponse.GetSelectorValues()) != 0 {
		t.Error("got non-empty selectors on an error response")
	}

	if err == nil {
		t.Error("got no attestation error on an error response")
	}
}
