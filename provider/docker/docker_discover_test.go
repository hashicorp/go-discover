package docker_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/docker/docker/api"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/docker"
)

func createContainer(t *testing.T, cli *client.Client, ctx context.Context, labels map[string]string) func() {

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:  "alpine",
		Tty:    true,
		Labels: labels,
	}, nil, nil, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		t.Fatal(err)
	}

	return func() {
		timeout := 10 * time.Second
		if err := cli.ContainerStop(ctx, resp.ID, &timeout); err != nil {
			t.Fatal(err)
		}

		if err := cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestAddrs(t *testing.T) {
	if os.Getenv("DOCKER_ENV") == "" {
		t.Skip("Skipping Docker test in non-docker env")
	}

	ctx := context.Background()

	version := api.DefaultVersion
	if v := os.Getenv("DOCKER_VERSION"); v != "" {
		version = v
	}

	cli, err := client.NewClientWithOpts(client.WithVersion(version))
	if err != nil {
		t.Fatal(err)
	}

	if _, err := cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{}); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		Name        string
		Services    []map[string]string // Containers Labels
		Label_Key   string
		Label_Value string
		Expected    int
	}{
		{
			"One container (One match)",
			[]map[string]string{
				map[string]string{"test1": "test2"},
			},
			"test1",
			"test2",
			1,
		},
		{
			"Two containers (One match)",
			[]map[string]string{
				map[string]string{"test1": "test2"},
				map[string]string{"test3": "test4"},
			},
			"test1",
			"test2",
			1,
		},
		{
			"Two containers (Two matches)",
			[]map[string]string{
				map[string]string{"test1": "test2"},
				map[string]string{"test1": "test2"},
			},
			"test1",
			"test2",
			2,
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			for _, labels := range tt.Services {
				close := createContainer(t, cli, ctx, labels)
				defer close()
			}

			args := discover.Config{
				"provider":    "docker",
				"label_key":   tt.Label_Key,
				"label_value": tt.Label_Value,
			}

			l := log.New(os.Stderr, "", log.LstdFlags)
			p := &docker.Provider{}

			addrs, err := p.Addrs(args, l)
			if err != nil {
				t.Fatal(err)
			}

			if len(addrs) != tt.Expected {
				t.Fatalf("bad: %v", addrs)
			}
		})
	}
}
