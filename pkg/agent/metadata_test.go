package agent

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMissingEnvVar(t *testing.T) {
	provider := ecsMetadataProvider{}

	meta, err := provider.getMetadata()

	if !meta.IsEmpty || err == nil {
		t.Error("no environment variable was set, but no error response was returned")
	}
}

func TestMetadataResponse(t *testing.T) {
	provider := ecsMetadataProvider{}
	server := httptest.NewServer(http.HandlerFunc(staticMetadataResponder))
	defer server.Close()

	t.Setenv(metadataEnvVarName, server.URL)

	meta, err := provider.getMetadata()

	if meta.IsEmpty || err != nil {
		t.Fatal(
			"error response received from static data.",
		)
	}

	var subTests = []struct {
		name, got, want string
	}{
		{"Name", meta.Root.Name, "application"},
		{"Image", meta.Root.Image, "awsid.dkr.ecr.us-east-1.amazonaws.com/myapp:latest"},
		{"Cluster", meta.Task.Cluster, "arn:aws:ecs:us-east-1:awsid:cluster/my-cluster"},
		{"TaskARN", meta.Task.TaskARN, "arn:aws:ecs:us-east-1:awsid:task/my-cluster/958f1db5294a45a6a5f1b2445d4d7362"},
		{"EnvironmentLabel", meta.Root.Labels["environment"], "development"},
		{"OwnerLabel", meta.Root.Labels["owning_team"], "dev-team-a"},
		{"Family", meta.Task.Family, "service-a"},
		{"Revision", meta.Task.Revision, "11"},
		{"AvailabilityZone", meta.Task.AvailabilityZone, "us-east-1d"},
	}

	for _, test := range subTests {
		testName := fmt.Sprintf("Test%sParsing", test.name)
		t.Run(testName, func(t *testing.T) {
			if test.got != test.want {
				t.Errorf("got: %s want: %s", test.got, test.want)
			}
		})
	}
}

func Test5XXMetadataResponse(t *testing.T) {
	errorResponder := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprint(w, "502 Bad Gateway")
	}

	provider := ecsMetadataProvider{}
	server := httptest.NewServer(http.HandlerFunc(errorResponder))
	defer server.Close()

	t.Setenv(metadataEnvVarName, server.URL)

	meta, err := provider.getMetadata()

	if !meta.IsEmpty {
		t.Error("expecting empty response when error status is received")
	}

	if err == nil {
		t.Error("expecting error to be populated when 5XX status code is received")
	}
}

func Test4XXMetadataResponse(t *testing.T) {
	errorResponder := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}

	provider := ecsMetadataProvider{}
	server := httptest.NewServer(http.HandlerFunc(errorResponder))
	defer server.Close()

	t.Setenv(metadataEnvVarName, server.URL)

	meta, err := provider.getMetadata()

	if !meta.IsEmpty {
		t.Error("expecting empty response when error status is received")
	}

	if err == nil {
		t.Error("expecting error to be populated when 5XX status code is received")
	}
}

func TestMetadataConnectionRefused(t *testing.T) {
	// Set to localhost but a port that's not open, in order to force a connection error
	t.Setenv(metadataEnvVarName, "http://localhost:8228")

	provider := ecsMetadataProvider{}

	meta, err := provider.getMetadata()

	if !meta.IsEmpty {
		t.Error("expecting empty response on connection error")
	}

	if err == nil {
		t.Error("expecting error to be populated on connection error")
	}
}
