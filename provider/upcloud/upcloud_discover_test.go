package upcloud_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/upcloud"
	"github.com/stretchr/testify/assert"
)

func TestAddrs(t *testing.T) {
	tests := []struct {
		name   string
		config discover.Config
	}{
		{
			name: "tag",
			config: discover.Config{
				"provider": "upcloud",
				"tag":      "vault-dev",
				"username": os.Getenv("UPCLOUD_API_USERNAME"),
				"password": os.Getenv("UPCLOUD_API_PASSWORD"),
			},
		},
		{
			name: "title_match",
			config: discover.Config{
				"provider":    "upcloud",
				"title_match": `^terraform\.test[0-9]+`,
				"username":    os.Getenv("UPCLOUD_API_USERNAME"),
				"password":    os.Getenv("UPCLOUD_API_PASSWORD"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.config["username"] == "" || tc.config["password"] == "" {
				t.Skip("UpCloud credentials missing")
			}

			p := upcloud.Provider{}
			l := log.New(os.Stderr, "", log.LstdFlags)
			addrs, err := p.Addrs(tc.config, l)
			assert.NoError(t, err)
			assert.NotEmpty(t, addrs)
			t.Logf("UpCloud found the following IP Addresses: %v", addrs)
		})
	}
}
