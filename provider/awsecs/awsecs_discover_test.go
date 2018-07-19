package awsecs_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/awsecs"
)

func TestAddrs(t *testing.T) {
	args := discover.Config{
		"provider":          "awsecs",
		"region":            os.Getenv("AWS_REGION"),
		"cluster_name":      "go-discover-cluster",
		"service_port":      "80",
		"service_name":      "test-ecs-service",
		"container_name":    "nginx",
		"access_key_id":     os.Getenv("AWS_ACCESS_KEY_ID"),
		"secret_access_key": os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}

	if args["region"] == "" || args["access_key_id"] == "" || args["secret_access_key"] == "" {
		t.Skip("AWS credentials or region missing")
	}

	p := &awsecs.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}

}
