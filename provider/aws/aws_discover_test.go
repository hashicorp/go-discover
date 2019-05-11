package aws_test

import (
	"log"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/aws"
)

func TestAddrs(t *testing.T) {
	args := discover.Config{
		"provider":          "aws",
		"region":            os.Getenv("AWS_REGION"),
		"tag_filters":       "consul=server",
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

func TestCreateTagFilterMap(t *testing.T) {
	tagFilters := "type=server,environment=dev,service=consul"

	answer := aws.TagFilterMap{
		"type":        "server",
		"environment": "dev",
		"service":     "consul",
	}

	check := aws.CreateTagFilterMap(tagFilters)

	if !cmp.Equal(check, answer) {

		t.Fatalf("The result of %v does not match %v", check, answer)
	}

}
