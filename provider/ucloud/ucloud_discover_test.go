package ucloud_test

import (
	"github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/ucloud"
	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"strings"
	"testing"
)

var l = log.New(os.Stderr, "", log.LstdFlags)

func TestAddrs_Incorrect_Provider_Name(t *testing.T) {
	args := discover.Config{
		"provider": "aws",
	}
	p := new(ucloud.Provider)
	addrs, err := p.Addrs(args, l)
	assert.Nil(t, addrs)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "provider"))
}

func TestAddrs_Required_Config_Missing(t *testing.T) {
	p := new(ucloud.Provider)
	stub := gostub.Stub(&ucloud.GetFromEnv, func(_ string) string {
		return ""
	})
	defer stub.Reset()
	inputs := []struct {
		Config         discover.Config
		ExpectedSubstr string
	}{
		{
			discover.Config{
				"provider": "ucloud",
			},
			"region",
		},
		{
			discover.Config{
				"provider": "ucloud",
				"region":   "cn-sh2",
			},
			"project_id",
		},
		{
			discover.Config{
				"provider":   "ucloud",
				"region":     "cn-sh2",
				"project_id": "this is a project id",
			},
			"tag",
		},
	}
	for _, input := range inputs {
		t.Run(strings.Title(input.ExpectedSubstr), func(t *testing.T) {
			addrs, err := p.Addrs(input.Config, l)
			assert.Nil(t, addrs)
			assert.NotNil(t, err)
			assert.True(t, strings.Contains(err.Error(), input.ExpectedSubstr))
		})
	}
}

func TestAddrs(t *testing.T) {
	p := new(ucloud.Provider)
	stub := gostub.Stub(&ucloud.GetFromEnv, func(_ string) string {
		panic("Should Not Read Env Here")
	})
	defer stub.Reset()

	if !allRequiredEnvPresent() {
		t.Skip("missing required env")
	}

	config := discover.Config{
		"provider":    "ucloud",
		"region":      os.Getenv("UCLOUD_REGION"),
		"project_id":  os.Getenv("UCLOUD_PROJECT_ID"),
		"tag":         "UCloud",
		"public_key":  os.Getenv("UCLOUD_PUBLIC_KEY"),
		"private_key": os.Getenv("UCLOUD_PRIVATE_KEY"),
	}

	addrs, err := p.Addrs(config, l)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(addrs))
}

func TestAddrsWithEnv(t *testing.T) {
	p := new(ucloud.Provider)

	if !allRequiredEnvPresent() {
		t.Skip("missing required env")
	}

	config := discover.Config{
		"provider": "ucloud",
		"tag":      "UCloud",
	}

	addrs, err := p.Addrs(config, l)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(addrs))
}

func TestAddrs_With_Optional_Configs(t *testing.T) {
	p := new(ucloud.Provider)
	if !allRequiredEnvPresent() {
		t.Skip("missing required env")
	}

	vpcID := os.Getenv("TESTADDRS_UCLOUD_VPC_ID")
	subnetID := os.Getenv("TESTADDRS_UCLOUD_SUBNET_ID")
	zone := os.Getenv("TESTADDRS_UCLOUD_ZONE")
	if vpcID == "" || subnetID == "" || zone == "" {
		t.Skip("missing TESTADDRS_UCLOUD_VPC_ID/TESTADDRS_UCLOUD_SUBNET_ID/TESTADDRS_UCLOUD_ZONE Environments")
	}

	config := discover.Config{
		"provider":  "ucloud",
		"tag":       "UCloud",
		"vpc_id":    vpcID,
		"subnet_id": subnetID,
	}

	addrs, err := p.Addrs(config, l)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(addrs))
}

func TestAddrs_Get_BGP_IP(t *testing.T) {
	p := new(ucloud.Provider)
	stub := gostub.Stub(&ucloud.GetFromEnv, func(_ string) string {
		panic("Should Not Read Env Here")
	})
	defer stub.Reset()

	if !allRequiredEnvPresent() {
		t.Skip("missing required env")
	}

	config := discover.Config{
		"provider":    "ucloud",
		"region":      os.Getenv("UCLOUD_REGION"),
		"project_id":  os.Getenv("UCLOUD_PROJECT_ID"),
		"tag":         "UCloud",
		"public_key":  os.Getenv("UCLOUD_PUBLIC_KEY"),
		"private_key": os.Getenv("UCLOUD_PRIVATE_KEY"),
		"ip_type":     "BGP",
	}

	addrs, err := p.Addrs(config, l)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(addrs))
	for _, addr := range addrs {
		assert.False(t, isPrivateIP(addr))
	}
}

func TestAddrs_Invalid_IP_Type(t *testing.T) {
	p := new(ucloud.Provider)
	stub := gostub.Stub(&ucloud.GetFromEnv, func(_ string) string {
		panic("Should Not Read Env Here")
	})
	defer stub.Reset()

	if !allRequiredEnvPresent() {
		t.Skip("missing required env")
	}

	config := discover.Config{
		"provider":    "ucloud",
		"region":      os.Getenv("UCLOUD_REGION"),
		"project_id":  os.Getenv("UCLOUD_PROJECT_ID"),
		"tag":         "UCloud",
		"public_key":  os.Getenv("UCLOUD_PUBLIC_KEY"),
		"private_key": os.Getenv("UCLOUD_PRIVATE_KEY"),
		"ip_type":     "invalid",
	}

	addrs, err := p.Addrs(config, l)
	assert.Nil(t, addrs)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "ip_type"))
}

func isPrivateIP(addr string) bool {
	return strings.HasPrefix(addr, "10.")
}

func allRequiredEnvPresent() bool {
	return !(os.Getenv("UCLOUD_REGION") == "" || os.Getenv("UCLOUD_PROJECT_ID") == "" || os.Getenv("UCLOUD_PUBLIC_KEY") == "" || os.Getenv("UCLOUD_PRIVATE_KEY") == "")
}
