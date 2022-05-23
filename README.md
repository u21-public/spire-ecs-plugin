# ECS Workload Attestation Plugin for SPIRE

This plugin provides workload attestation for ECS, particularly tasks running on [Fargate](https://aws.amazon.com/fargate/). Since organizations running on Fargate don't have access to the underlying EC2 host,
nor the ability to customize any of the services running there the SPIRE agent must be run in-container. 
This environment is rather different from most of the ones SPIRE is designed around, and existing plugins don't fit this model very well.

The spire-ecs-plugin allows agents to attest workloads using native ECS attributes such as task names, docker images, and AWS labels among others. 
It does so by inspecting the `ECS_CONTAINER_METADATA_URI_V4` endpoint. 

## Configuration
There is no configurable `plugin_data` stanza for this plugin, since it solely uses the AWS provided workload metadata.

Sample configuration:
```
    WorkloadAttestor "ecs" {
        plugin_data {}
        plugin_cmd = "/opt/spire/ecs-attestor"
        plugin_checksum = "your_release_checksum"
    }
```

## Workload Selectors
A subset of the fields from the [ECS metadata v4 endpoint](https://docs.aws.amazon.com/AmazonECS/latest/userguide/task-metadata-endpoint-v4-fargate.html) are available as selectors:

|Selector  |Example                 |Description                                      |
|----------|------------------------|-------------------------------------------------|
|`ecs:name`|`ecs:name:application`  |The container name within the ECS task definition.|
|`ecs:image`|`ecs:image:docker.mycompany.com/images/app:latest`|The full image path + tag of the current image.|
|`ecs:label`|`ecs:label:stage:development`|All AWS resource labels applied to this ECS service, as key:value pairs.|
|`ecs:cluster`|`ecs:cluster:aws:arn:someaccount:someregion:somecluster`|The ARN of the cluster in which this task is running.|
|`ecs:arn`|`ecs:arn:aws:arn:someaccount:someregion:sometask`|The ARN of this ECS task.|
|`ecs:family`|`ecs:family:myservice`|The name of the ECS service associated with this task.|
|`ecs:revision`|`ecs:revision:11`|The revision number of the ECS task definition for this container.|
|`ecs:availabilityzone`|`ecs:availabilityzone:us-east-1d`|The AZ where this task is running.|

## Installing
```
go build -o ecs-attestor github.com/u21-public/spire-ecs-plugin/cmd/agent
```
