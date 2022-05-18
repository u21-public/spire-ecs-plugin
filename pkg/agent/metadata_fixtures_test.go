package agent

import (
	"fmt"
	"net/http"
)

func staticMetadataResponder(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/task" {
		fmt.Fprint(w, staticTaskResponse)
	} else {
		fmt.Fprint(w, staticMetaResponse)
	}
}

const staticMetaResponse string = `
{
    "DockerId": "958f1db5294a45a6a5f1b2445d4d7362-524788293",
    "Name": "application",
    "DockerName": "application",
    "Image": "awsid.dkr.ecr.us-east-1.amazonaws.com/myapp:latest",
    "ImageID": "sha256:9c6cca6e6abcb4c02fd33940501b06565edd348273a76e5a5a2a83c4ca418a11",
    "Labels": {
		"environment": "development",
		"owning_team": "dev-team-a"
    },
    "DesiredStatus": "RUNNING",
    "KnownStatus": "RUNNING",
    "Limits": {
        "CPU": 2
    },
    "CreatedAt": "2022-05-10T19:09:25.679023799Z",
    "StartedAt": "2022-05-10T19:09:25.679023799Z",
    "Type": "NORMAL",
    "Networks": [
        {
            "NetworkMode": "awsvpc",
            "IPv4Addresses": [
                "1.1.1.1"
            ],
            "AttachmentIndex": 0,
            "MACAddress": "0e:4a:13:36:84:cf",
            "IPv4SubnetCIDRBlock": "1.1.1.0/24",
            "DomainNameServers": [
                "1.1.1.2"
            ],
            "DomainNameSearchList": [
                "us-east-1.compute.internal"
            ],
            "PrivateDNSName": "ip-1-1-1-1.us-east-1.compute.internal",
            "SubnetGatewayIpv4Address": "1.1.1.1/24"
        }
    ],
    "ContainerARN": "arn:aws:ecs:us-east-1:awsid:container/my-cluster/958f1db5294a45a6a5f1b2445d4d7362/7095f5f9-6a0a-425c-bcd0-ed8bdea0c802",
    "LogOptions": {
        "awslogs-group": "myloggroup",
        "awslogs-region": "us-east-1",
        "awslogs-stream": "service/application/958f1db5294a45a6a5f1b2445d4d7362"
    },
    "LogDriver": "awslogs"
}
`

const staticTaskResponse string = `
{
    "Cluster": "arn:aws:ecs:us-east-1:awsid:cluster/my-cluster",
    "TaskARN": "arn:aws:ecs:us-east-1:awsid:task/my-cluster/958f1db5294a45a6a5f1b2445d4d7362",
    "Family": "service-a",
    "Revision": "11",
    "DesiredStatus": "RUNNING",
    "KnownStatus": "RUNNING",
    "Limits": {
        "CPU": 0.5,
        "Memory": 1024
    },
    "PullStartedAt": "2022-05-10T19:08:58.424906316Z",
    "PullStoppedAt": "2022-05-10T19:09:20.630175811Z",
    "AvailabilityZone": "us-east-1d",
    "Containers": [
        {
            "DockerId": "958f1db5294a45a6a5f1b2445d4d7362-524788293",
            "Name": "application",
            "DockerName": "application",
            "Image": "awsid.dkr.ecr.us-east-1.amazonaws.com/zero-trust:latest",
            "ImageID": "sha256:9c6cca6e6abcb4c02fd33940501b06565edd348273a76e5a5a2a83c4ca418a11",
            "Labels": {
                "com.amazonaws.ecs.cluster": "arn:aws:ecs:us-east-1:awsid:cluster/my-cluster",
                "com.amazonaws.ecs.container-name": "application",
                "com.amazonaws.ecs.task-arn": "arn:aws:ecs:us-east-1:awsid:task/my-cluster/958f1db5294a45a6a5f1b2445d4d7362",
                "com.amazonaws.ecs.task-definition-family": "service-a",
                "com.amazonaws.ecs.task-definition-version": "11"
            },
            "DesiredStatus": "RUNNING",
            "KnownStatus": "RUNNING",
            "Limits": {
                "CPU": 2
            },
            "CreatedAt": "2022-05-10T19:09:25.679023799Z",
            "StartedAt": "2022-05-10T19:09:25.679023799Z",
            "Type": "NORMAL",
            "Networks": [
                {
                    "NetworkMode": "awsvpc",
                    "IPv4Addresses": [
                        "1.1.1.1"
                    ],
                    "AttachmentIndex": 0,
                    "MACAddress": "0e:4a:13:36:84:cf",
                    "IPv4SubnetCIDRBlock": "1.1.1.0/24",
                    "DomainNameServers": [
                        "1.1.1.2"
                    ],
                    "DomainNameSearchList": [
                        "us-east-1.compute.internal"
                    ],
                    "PrivateDNSName": "ip-1-1-1-1.us-east-1.compute.internal",
                    "SubnetGatewayIpv4Address": "1.1.1.1/24"
                }
            ],
            "ContainerARN": "arn:aws:ecs:us-east-1:awsid:container/my-cluster/958f1db5294a45a6a5f1b2445d4d7362/7095f5f9-6a0a-425c-bcd0-ed8bdea0c802",
            "LogOptions": {
                "awslogs-group": "myloggroup",
                "awslogs-region": "us-east-1",
                "awslogs-stream": "myservice/application/958f1db5294a45a6a5f1b2445d4d7362"
            },
            "LogDriver": "awslogs"
        }
    ],
    "LaunchType": "FARGATE",
    "ClockDrift": {
        "ClockErrorBound": 0.403273,
        "ReferenceTimestamp": "2022-05-10T19:33:11Z",
        "ClockSynchronizationStatus": "SYNCHRONIZED"
    }
}
`
