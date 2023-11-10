package aws_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/aws"
)

func TestAddrs(t *testing.T) {
	args := discover.Config{
		"provider":          "aws",
		"region":            os.Getenv("AWS_REGION"),
		"tag_key":           "consul",
		"tag_value":         "server",
		"access_key_id":     os.Getenv("AWS_ACCESS_KEY_ID"),
		"secret_access_key": os.Getenv("AWS_SECRET_ACCESS_KEY"),
		"session_token":     os.Getenv("AWS_SESSION_TOKEN"),
	}

	if args["region"] == "" || args["access_key_id"] == "" || args["secret_access_key"] == "" {
		t.Skip("AWS credentials or region missing")
	}

	p := &aws.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}
func TestAddrsEndpoint(t *testing.T) {
	args := discover.Config{
		"provider":          "aws",
		"region":            os.Getenv("AWS_REGION"),
		"tag_key":           "consul",
		"tag_value":         "server",
		"access_key_id":     os.Getenv("AWS_ACCESS_KEY_ID"),
		"secret_access_key": os.Getenv("AWS_SECRET_ACCESS_KEY"),
		"endpoint":          os.Getenv("AWS_EC2_METADATA_SERVICE_ENDPOINT"),
	}

	if args["region"] == "" || args["access_key_id"] == "" || args["secret_access_key"] == "" {
		t.Skip("AWS credentials or region missing")
	}

	p := &aws.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}
func TestAddrsECS(t *testing.T) {
	args := discover.Config{
		"provider":          "aws",
		"service":           "ecs",
		"region":            os.Getenv("AWS_REGION"),
		"tag_key":           "consul",
		"tag_value":         "server",
		"access_key_id":     os.Getenv("AWS_ACCESS_KEY_ID"),
		"secret_access_key": os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}

	if args["region"] == "" || args["access_key_id"] == "" || args["secret_access_key"] == "" {
		t.Skip("AWS credentials or region missing")
	}

	p := &aws.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 3 {
		t.Fatalf("bad: %v", addrs)
	}
}

func TestAddrsECSFilterClusterName(t *testing.T) {
	args := discover.Config{
		"provider":          "aws",
		"service":           "ecs",
		"ecs_cluster":       "go-discover-2",
		"region":            os.Getenv("AWS_REGION"),
		"tag_key":           "consul",
		"tag_value":         "server",
		"access_key_id":     os.Getenv("AWS_ACCESS_KEY_ID"),
		"secret_access_key": os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}

	if args["region"] == "" || args["access_key_id"] == "" || args["secret_access_key"] == "" {
		t.Skip("AWS credentials or region missing")
	}

	p := &aws.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}

func TestAddrsECSFilterFuzzyClusterName(t *testing.T) {
	args := discover.Config{
		"provider":          "aws",
		"service":           "ecs",
		"ecs_cluster":       "go-discover-1",
		"region":            os.Getenv("AWS_REGION"),
		"tag_key":           "consul",
		"tag_value":         "server",
		"access_key_id":     os.Getenv("AWS_ACCESS_KEY_ID"),
		"secret_access_key": os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}

	if args["region"] == "" || args["access_key_id"] == "" || args["secret_access_key"] == "" {
		t.Skip("AWS credentials or region missing")
	}

	p := &aws.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 1 {
		t.Fatalf("bad: %v", addrs)
	}
}

func TestAddrsECSFilterTaskFamily(t *testing.T) {
	args := discover.Config{
		"provider":          "aws",
		"service":           "ecs",
		"ecs_family":        "go-discover-familia",
		"region":            os.Getenv("AWS_REGION"),
		"tag_key":           "consul",
		"tag_value":         "server",
		"access_key_id":     os.Getenv("AWS_ACCESS_KEY_ID"),
		"secret_access_key": os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}

	if args["region"] == "" || args["access_key_id"] == "" || args["secret_access_key"] == "" {
		t.Skip("AWS credentials or region missing")
	}

	p := &aws.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}

}
