package kubernetes_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/kubernetes"
)

func TestAddrs(t *testing.T) {
	args := discover.Config{
		"provider":  "kubernetes",
		"namespace": os.Getenv("AWS_REGION"),
		"service":   "consul-role",
	}

	if args["namespace"] == "" || args["service"] == "" {
		t.Skip("Kubernetes namespace or service missing")
	}

	p := &kubernetes.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 3 {
		t.Fatalf("bad: %v", addrs)
	}
}
