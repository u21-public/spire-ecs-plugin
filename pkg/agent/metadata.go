package agent

/*
ECS has a couple of useful metadata endpoints for workload attestation
with the base URL specified by the environment variable ECS_CONTAINER_METADATA_URI_V4:

- The root, which contains information about the container name, docker image used, and AWS labels applied
- ECS_CONTAINER_METADATA_URI_V4/task, which provides details about the ECS cluster, task ARN,
  task defition revision, and task family (ECS service name)
*/

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

const metadataEnvVarName = "ECS_CONTAINER_METADATA_URI_V4"

type metadataRootResponse struct {
	Name   string            `json:"Name"`
	Image  string            `json:"Image"`
	Labels map[string]string `json:"Labels"`
}

type metadataTaskResponse struct {
	Cluster          string `json:"Cluster"`
	TaskARN          string `json:"TaskARN"`
	Family           string `json:"Family"`
	Revision         string `json:"Revision"`
	AvailabilityZone string `json:"AvailabilityZone"`
}

type metadataResponse struct {
	IsEmpty bool
	Root    metadataRootResponse
	Task    metadataTaskResponse
}

type metadataProvider interface {
	getMetadata() (metadataResponse, error)
}

type ecsMetadataProvider struct{}

func (e *ecsMetadataProvider) getMetadata() (metadataResponse, error) {
	metadata_url := os.Getenv(metadataEnvVarName)
	if metadata_url == "" {
		return metadataResponse{IsEmpty: true}, errors.New("metadata url environment variable not found")
	}

	client := http.Client{}

	root_resp, err := client.Get(metadata_url)
	if err != nil {
		return metadataResponse{IsEmpty: true}, err
	}
	defer root_resp.Body.Close()

	task_resp, err := client.Get(metadata_url + "/task")
	if err != nil {
		return metadataResponse{IsEmpty: true}, err
	}
	defer task_resp.Body.Close()

	var root_data metadataRootResponse
	var task_data metadataTaskResponse

	err = json.NewDecoder(root_resp.Body).Decode(&root_data)
	if err != nil {
		return metadataResponse{IsEmpty: true}, err
	}

	err = json.NewDecoder(task_resp.Body).Decode(&task_data)
	if err != nil {
		return metadataResponse{IsEmpty: true}, err
	}

	return metadataResponse{
		IsEmpty: false,
		Root:    root_data,
		Task:    task_data,
	}, nil
}
