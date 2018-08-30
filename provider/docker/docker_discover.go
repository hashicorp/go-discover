// Package docker provides node discovery for Docker.
package docker

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/docker/docker/api"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Provider struct{}

func (p *Provider) Help() string {
	return `Docker:

    provider:		"docker"
    label_key:		The tag key to filter on
    label_value:	The tag value to filter on
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "docker" {
		return nil, fmt.Errorf("discover-docker: invalid provider " + args["provider"])
	}

	labelKey := args["label_key"]
	labelValue := args["label_value"]

	ctx := context.Background()

	version := api.DefaultVersion
	if v := os.Getenv("DOCKER_VERSION"); v != "" {
		version = v
	}

	cli, err := client.NewClientWithOpts(client.WithVersion(version))
	if err != nil {
		return nil, fmt.Errorf("discover-docker: %s", err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, fmt.Errorf("discover-docker: %s", err)
	}

	var addrs []string
	for _, container := range containers {
		key, ok := container.Labels[labelKey]
		if ok && key == labelValue {
			for _, network := range container.NetworkSettings.Networks {
				addrs = append(addrs, network.IPAddress)
			}
		}
	}

	return addrs, nil
}
