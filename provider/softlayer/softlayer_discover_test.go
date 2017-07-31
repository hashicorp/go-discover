package softlayer_test

import (
	"log"
	"os"
	"reflect"
	"testing"

	discover "github.com/hashicorp/go-discover"

	_ "github.com/hashicorp/go-discover/provider/softlayer"
)

func TestAddrs(t *testing.T) {
	if got, want := discover.ProviderNames(), []string{"softlayer"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got providers %v want %v", got, want)
	}

	cfg := discover.Config{
		"provider":   "softlayer",
		"username":   os.Getenv("SL_USERNAME"),
		"api_key":    os.Getenv("SL_API_KEY"),
		"datacenter": "dal06",
		"tag_value":  "consul-server",
	}
	if cfg["username"] == "" || cfg["api_key"] == "" {
		t.Skip("SoftLayer credentials missing")
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
