package azure_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/azure"
)

var _ discover.Provider = (*azure.Provider)(nil)
var _ discover.ProviderWithUserAgent = (*azure.Provider)(nil)

func TestTagAddrsWithEnv(t *testing.T) {
	args := discover.Config{
		"provider":  "azure",
		"tag_name":  "consul",
		"tag_value": "server",
	}

	if os.Getenv("ARM_SUBSCRIPTION_ID") == "" || os.Getenv("ARM_CLIENT_ID") == "" || os.Getenv("ARM_CLIENT_SECRET") == "" || os.Getenv("ARM_TENANT_ID") == "" {
		t.Skip("Azure Enviornmental credentials missing")
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
	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}
func TestTagAddrs(t *testing.T) {
	args := discover.Config{
		"provider":          "azure",
		"tag_name":          "consul",
		"tag_value":         "server",
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
	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}

func TestVmScaleSetAddrs(t *testing.T) {
	args := discover.Config{
		"provider":          "azure",
		"resource_group":    "go-discover-azure-vmss-dev",
		"vm_scale_set":      "go-discover-azure-vmss-01-scale-set",
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
