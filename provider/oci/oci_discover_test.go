package oci_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/oci"
)

var tests = []struct {
	name      string
	config    discover.Config
	addrCount int
}{
	{
		"freeform",
		discover.Config{
			"provider"              : "oci",
			"tenancy_ocid"          : os.Getenv("OCI_TENANCY_OCID"),
			"user_ocid"             : os.Getenv("OCI_USER_OCID"),
			"region"                : os.Getenv("OCI_REGION"),
			"key_fingerprint"       : os.Getenv("OCI_KEY_FINGERPRINT"),
			"private_key"           : os.Getenv("OCI_PRIVATE_KEY"),
			"private_key_passphrase": os.Getenv("OCI_PRIVATE_KEY_PASSPHRASE"),
			"tag_key"               : "discover",
			"tag_value"             : "me",
		},
		1,
	},
	{
		"defined",
		discover.Config{
			"provider"              : "oci",
			"tenancy_ocid"          : os.Getenv("OCI_TENANCY_OCID"),
			"user_ocid"             : os.Getenv("OCI_USER_OCID"),
			"region"                : os.Getenv("OCI_REGION"),
			"key_fingerprint"       : os.Getenv("OCI_KEY_FINGERPRINT"),
			"private_key"           : os.Getenv("OCI_PRIVATE_KEY"),
			"private_key_passphrase": os.Getenv("OCI_PRIVATE_KEY_PASSPHRASE"),
			"tag_namespace"         : "defined",
			"tag_key"               : "discover",
			"tag_value"             : "me",
		},
		2,
	},
	{
		"freePartial",
		discover.Config{
			"provider"              : "oci",
			"tenancy_ocid"          : os.Getenv("OCI_TENANCY_OCID"),
			"user_ocid"             : os.Getenv("OCI_USER_OCID"),
			"region"                : os.Getenv("OCI_REGION"),
			"key_fingerprint"       : os.Getenv("OCI_KEY_FINGERPRINT"),
			"private_key"           : os.Getenv("OCI_PRIVATE_KEY"),
			"private_key_passphrase": os.Getenv("OCI_PRIVATE_KEY_PASSPHRASE"),
			"tag_key"               : "discover",
		},
		1,
	},
	{
		"definedPartial",
		discover.Config{
			"provider"              : "oci",
			"tenancy_ocid"          : os.Getenv("OCI_TENANCY_OCID"),
			"user_ocid"             : os.Getenv("OCI_USER_OCID"),
			"region"                : os.Getenv("OCI_REGION"),
			"key_fingerprint"       : os.Getenv("OCI_KEY_FINGERPRINT"),
			"private_key"           : os.Getenv("OCI_PRIVATE_KEY"),
			"private_key_passphrase": os.Getenv("OCI_PRIVATE_KEY_PASSPHRASE"),
			"tag_namespace"         : "defined",
			"tag_key"               : "discover",
		},
		2,
	},
	{
		"definedPublic",
		discover.Config{
			"provider"              : "oci",
			"tenancy_ocid"          : os.Getenv("OCI_TENANCY_OCID"),
			"user_ocid"             : os.Getenv("OCI_USER_OCID"),
			"region"                : os.Getenv("OCI_REGION"),
			"key_fingerprint"       : os.Getenv("OCI_KEY_FINGERPRINT"),
			"private_key"           : os.Getenv("OCI_PRIVATE_KEY"),
			"private_key_passphrase": os.Getenv("OCI_PRIVATE_KEY_PASSPHRASE"),
			"tag_namespace"         : "defined",
			"tag_key"               : "discover",
			"tag_value"             : "me",
			"addr_type"             : "public",
		},
		1,
	},
}

func TestAddrs(t *testing.T) {
	if os.Getenv("OCI_TENANCY_OCID")    == "" ||
		 os.Getenv("OCI_USER_OCID")       == "" ||
		 os.Getenv("OCI_REGION")          == "" ||
		 os.Getenv("OCI_KEY_FINGERPRINT") == "" ||
		 os.Getenv("OCI_PRIVATE_KEY")     == "" {
			t.Skip("OCI credentials missing.")
	}
	
	p := &oci.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	for _, test := range tests {
		l.Printf("[INFO] Begin Test: %s", test.name)
		addrs, err := p.Addrs(test.config, l)
		if err != nil {
			t.Error(err)
		}
		if len(addrs) != test.addrCount {
			t.Errorf("bad: %v", addrs)
		}
	}
}