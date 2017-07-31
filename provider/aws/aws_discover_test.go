package aws_test

import (
	"log"
	"os"
	"reflect"
	"testing"

	discover "github.com/hashicorp/go-discover"

	_ "github.com/hashicorp/go-discover/provider/aws"
)

func TestAddrs(t *testing.T) {
	if got, want := discover.ProviderNames(), []string{"aws"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got providers %v want %v", got, want)
	}

	cfg := discover.Config{
		"provider":          "aws",
		"region":            os.Getenv("AWS_REGION"),
		"tag_key":           "consul-role",
		"tag_value":         "server",
		"access_key_id":     os.Getenv("AWS_ACCESS_KEY_ID"),
		"secret_access_key": os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}

	if cfg["region"] == "" || cfg["access_key_id"] == "" || cfg["secret_access_key"] == "" {
		t.Skip("AWS credentials or region missing")
	}

	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := discover.Addrs(cfg.String(), l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 3 {
		t.Fatalf("bad: %v", addrs)
	}
}
