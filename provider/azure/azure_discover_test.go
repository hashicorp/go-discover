package azure_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/azure"
)

func TestAddrs(t *testing.T) {
	args := discover.Config{
		"provider":          "azure",
		"tag_name":          "type",
		"tag_value":         "Foundation",
		"subscription_id":   os.Getenv("ARM_SUBSCRIPTION_ID"),
		"tenant_id":         os.Getenv("ARM_TENANT_ID"),
		"client_id":         os.Getenv("ARM_CLIENT_ID"),
		"secret_access_key": os.Getenv("ARM_CLIENT_SECRET"),
		"environment":       os.Getenv("ARM_ENVIRONMENT"),
	}

	if args["subscription_id"] == "" || args["client_id"] == "" || args["secret_access_key"] == "" || args["tenant_id"] == "" {
		t.Skip("Azure credentials missing")
	}

	if args["environment"] == "" {
		t.Log("Environments other than Public not supported at the moment")
	}

	p := &azure.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 3 {
		t.Fatalf("bad: %v", addrs)
	}
}
