// Package ucloud provides node discovery for UCloud.
package ucloud

import (
	"fmt"
	"github.com/ahmetb/go-linq"
	"github.com/ucloud/ucloud-sdk-go/services/uhost"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Provider struct {
}

func (p *Provider) Help() string {
	return `UCloud:

    provider:                    "ucloud"
    region(required):            The UCloud region (env UCLOUD_REGION)
    tag(required):               The tag value to filter on
	project_id(required):        The UCloud project id (env UCLOUD_PROJECT_ID)
    access_key_id(required):     The UCloud public key to use (env UCLOUD_PUBLIC_KEY)
    access_key_secret(required): The UCloud private key to use (env UCLOUD_PRIVATE_KEY)
	zone:                        The UCloud zone
	vpc_id:                      Target instance's vpc id
	subnet_id:                   Target instnace's subnet id
	ip_type:                     "private"/"bgp" (for mainland China)/"international" (for international), default to "private"
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "ucloud" {
		return nil, fmt.Errorf("discover-ucloud: invalid provider " + args["provider"])
	}
	l = discardIfNil(l)
	region, err := requiredConfig(args, "region", "UCLOUD_REGION", l)
	if err != nil {
		return nil, err
	}
	projectId, err := requiredConfig(args, "project_id", "UCLOUD_PROJECT_ID", l)
	if err != nil {
		return nil, err
	}
	zone := args["zone"]
	tag, err := requiredConfig(args, "tag", "", l)
	if err != nil {
		return nil, err
	}
	accessKeyID := argsOrEnv(args, "access_key_id", "UCLOUD_PUBLIC_KEY", discardLogger)
	accessKeySecret := argsOrEnv(args, "access_key_secret", "UCLOUD_PRIVATE_KEY", discardLogger)
	vpcID := args["vpc_id"]
	subnetID := args["subnet_id"]
	ipType := args["ip_type"]
	if ipType == "" {
		ipType = "Private"
	}
	if !ipTypeEqual(ipType, "private") && !ipTypeEqual(ipType, "bgp") && !ipTypeEqual(ipType, "international") {
		l.Printf("[DEBUG] discover-ucloud: invalid ip_type:%s", ipType)
		return nil, fmt.Errorf("invalid ip_type:%s", ipType)
	}
	ipType = rewriteIpTypeForFilter(ipType)
	l.Printf("[DEBUG] discover-ucloud: Using region=%s zone=%s project_id=%s vpc_id=%s subnet_id=%s tag=%s ", region, zone, projectId, vpcID, subnetID, tag)
	cfg := newConfig(projectId, region, zone)
	credential := newCredential(accessKeyID, accessKeySecret)
	req := &uhost.DescribeUHostInstanceRequest{
		CommonBase: request.CommonBase{
			Region:    stringConfig(region),
			ProjectId: stringConfig(projectId),
			Zone:      stringConfig(zone),
		},
		Tag:      ucloud.String(tag),
		Limit:    ucloud.Int(100),
		VPCId:    stringConfig(vpcID),
		SubnetId: stringConfig(subnetID),
	}
	client := uhost.NewClient(cfg, credential)
	response, err := client.DescribeUHostInstance(req)
	if err != nil {
		return nil, fmt.Errorf("discover-ucloud: DescribeUHostInstance failed: %s", err)
	}
	var addrs []string

	linq.From(response.UHostSet).Where(runningHost).Select(func(i interface{}) interface{} {
		instance := i.(uhost.UHostInstanceSet)
		ip := linq.From(instance.IPSet).FirstWith(filterIPBy(ipType)).(uhost.UHostIPSet).IP
		l.Printf("[DEBUG] discover-ucloud: Instance %s has innner ip %s ", instance.UHostId, ip)
		return ip
	}).ToSlice(&addrs)

	l.Printf("[DEBUG] discover-ucloud: Found %d running instances", len(addrs))
	l.Printf("[DEBUG] discover-ucloud: Found ip addresses: %v", addrs)
	return addrs, nil
}

func rewriteIpTypeForFilter(ipType string) string {
	if ipTypeEqual(ipType, "international") {
		ipType = "internation"
	}
	return ipType
}

func discardIfNil(l *log.Logger) *log.Logger {
	if l == nil {
		l = discardLogger
	}
	return l
}

func stringConfig(value string) *string {
	if value == "" {
		return nil
	}
	return ucloud.String(value)
}

func newCredential(accessKeyID string, accessKeySecret string) *auth.Credential {
	credential := auth.NewCredential()
	credential.PublicKey = accessKeyID
	credential.PrivateKey = accessKeySecret
	return &credential
}

func newConfig(projectID string, region string, zone string) *ucloud.Config {
	cfg := ucloud.NewConfig()
	cfg.ProjectId = projectID
	cfg.Region = region
	cfg.Zone = zone
	return &cfg
}

func requiredConfig(args map[string]string, key string, env string, l *log.Logger) (string, error) {
	value := argsOrEnv(args, key, env, l)
	if value == "" {
		l.Printf("[DEBUG] discover-ucloud: %s not provided", strings.Title(key))
		return "", fmt.Errorf("discover-ucloud: invalid %s:%s", key, value)
	}
	l.Printf("[INFO] discover-ucloud: %s is %s", strings.Title(key), value)
	return value, nil
}

func argsOrEnv(args map[string]string, key, env string, l *log.Logger) string {
	if value, ok := args[key]; ok {
		l.Printf("[INFO] discover-ucloud: %s is %s", strings.Title(key), value)
		return value
	}
	if env == "" {
		return ""
	}
	value := GetFromEnv(env)
	l.Printf("[INFO] discover-ucloud from env %s: %s is %s", env, strings.Title(key), value)
	return value
}

var GetFromEnv = func(env string) string {
	return os.Getenv(env)
}

var runningHost = func(i interface{}) bool {
	return i.(uhost.UHostInstanceSet).State == "Running"
}

func filterIPBy(ipType string) func(ipSet interface{}) bool {
	return func(ipSet interface{}) bool {
		return ipTypeEqual(ipSet.(uhost.UHostIPSet).Type, ipType)
	}
}

func ipTypeEqual(ipType, expected string) bool {
	return strings.EqualFold(ipType, expected)
}

var discardLogger = log.New(ioutil.Discard, "", 0)
