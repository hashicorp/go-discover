package azure_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
)

func TestAddrs(t *testing.T) {
	cfg := discover.Config{
		"provider":          "azure",
		"tag_name":          "type",
		"tag_value":         "Foundation",
		"subscription_id":   os.Getenv("ARM_SUBSCRIPTION_ID"),
		"tenant_id":         os.Getenv("ARM_TENANT_ID"),
		"client_id":         os.Getenv("ARM_CLIENT_ID"),
		"secret_access_key": os.Getenv("ARM_CLIENT_SECRET"),
		"environment":       os.Getenv("ARM_ENVIRONMENT"),
	}

	if cfg["subscription_id"] == "" || cfg["client_id"] == "" || cfg["secret_access_key"] == "" || cfg["tenant_id"] == "" {
		t.Skip("Azure credentials missing")
	}

	if cfg["environment"] == "" {
		t.Log("Environments other than Public not supported at the moment")
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
