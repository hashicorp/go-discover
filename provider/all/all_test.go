package all_test

import (
	"reflect"
	"testing"

	discover "github.com/hashicorp/go-discover"
)

func TestAddrs(t *testing.T) {
	names := []string{"aws", "azure", "gce", "softlayer"}
	if got, want := discover.ProviderNames(), names; !reflect.DeepEqual(got, want) {
		t.Fatalf("got providers %v want %v", got, want)
	}
}
