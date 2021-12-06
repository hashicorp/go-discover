// Package ksyun provides node discovery for Kingsoft Cloud.
package ksyun

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/KscSDK/ksc-sdk-go/ksc/utils"
	"github.com/KscSDK/ksc-sdk-go/service/kec"

	"github.com/mitchellh/mapstructure"
)

type Provider struct {
	userAgent string
}

func (p *Provider) SetUserAgent(s string) {
	p.userAgent = s
}

func (p *Provider) Help() string {
	return `Ksyun:

    provider:          "ksyun"
    region:            The ksyun region to filter on
    tag_key:           The tag key to filter on
    tag_value:         The tag value to filter on
    access_key_id:     The ksyun access key to use
    access_key_secret: The ksyun access key secret to use
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "ksyun" {
		return nil, fmt.Errorf("discover-ksyun: invalid provider " + args["provider"])
	}

	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}

	region := args["region"]
	tagKey := args["tag_key"]
	tagValue := args["tag_value"]
	accessKeyID := args["access_key_id"]
	accessKeySecret := args["access_key_secret"]

	l.Printf("[DEBUG] discover-ksyun: Using region=%s tag_key=%s tag_value=%s", region, tagKey, tagValue)
	if accessKeyID == "" && accessKeySecret == "" {
		l.Printf("[DEBUG] discover-ksyun: No credentials provided")
		return nil, fmt.Errorf("discover-ksyun: invalid access key and secret key")
	}

	if region == "" {
		l.Printf("[DEBUG] discover-ksyun: Region not provided")
		return nil, fmt.Errorf("discover-ksyun: invalid region")
	}
	l.Printf("[INFO] discover-ksyun: Region is %s", region)

	client := ksc.NewClient(accessKeyID, accessKeySecret)
	svc := kec.SdkNew(client, &ksc.Config{Region: &region}, &utils.UrlInfo{
		UseSSL: true,
	})

	// For detail, see https://docs.ksyun.com/documents/816
	filters := make(map[string]interface{})
	if tagKey != "" {
		filters["Filter.1.Name"] = tagKey
		filters["Filter.1.Value.1"] = tagValue
	}

	respRaw, err := svc.DescribeInstances(&filters)
	if err != nil {
		return nil, fmt.Errorf("discover-ksyun: DescribeInstances: %s", err)
	}

	var resp Reponse
	err = mapstructure.Decode(respRaw, &resp)
	if err != nil {
		return nil, fmt.Errorf("discover-ksyun: Parsing response failed: %s", err)
	}
	l.Printf("[DEBUG] discover-ksyun: Found total %d instances", resp.InstanceCount)

	var addrs []string
	for _, instance := range resp.InstanceSet {
		l.Printf("[DEBUG] discover-ksyun: Instance %s has innner ip %s ", instance.InstanceId, instance.PrivateIpAddress)
		addrs = append(addrs, instance.PrivateIpAddress)
	}

	return addrs, nil
}

// Processes the returned results of DescribeInstances.
// The original result is of map type, which is complex to handle,
// and is simplified by using Response.
type Reponse struct {
	InstanceCount int        `mapstructure:"InstanceCount"`
	InstanceSet   []Instance `mapstructure:"InstancesSet"`
}

// Instance returns a part of the result.
// Anything else not related to go-discovery is ignored.
type Instance struct {
	InstanceId       string `mapstructure:"InstanceId"`
	PrivateIpAddress string `mapstructure:"PrivateIpAddress"`
}
