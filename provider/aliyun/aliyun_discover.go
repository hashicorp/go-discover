// Package aliyun provides node discovery for Aliyun.
package aliyun

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

type Provider struct {
	userAgent string
}

func (p *Provider) SetUserAgent(s string) {
	p.userAgent = s
}

func (p *Provider) Help() string {
	return `Aliyun(Alibaba Cloud):

    provider:          "aliyun"
    region:            The Aliyun region.
    tag_key:           The tag key to filter on
    tag_value:         The tag value to filter on
    access_key_id:     The Aliyun access key to use
    access_key_secret: The Aliyun access key secret to use

	The required RAM permission is 'ecs:DescribeInstances'.
	It is recommended you make a dedicated key used only for auto-joining.
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "aliyun" {
		return nil, fmt.Errorf("discover-aliyun: invalid provider " + args["provider"])
	}

	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}

	region := args["region"]
	tagKey := args["tag_key"]
	tagValue := args["tag_value"]
	accessKeyID := args["access_key_id"]
	accessKeySecret := args["access_key_secret"]

	log.Printf("[DEBUG] discover-aliyun: Using region=%s tag_key=%s tag_value=%s", region, tagKey, tagValue)
	if accessKeyID == "" && accessKeySecret == "" {
		log.Printf("[DEBUG] discover-aliyun: No static credentials")
	} else {
		log.Printf("[DEBUG] discover-aliyun: Static credentials provided")
	}

	if region == "" {
		l.Printf("[DEBUG] discover-aliyun: Region not provided")
		return nil, fmt.Errorf("discover-aliyun: invalid region")
	}
	l.Printf("[INFO] discover-aliyun: Region is %s", region)

	svc, err := ecs.NewClientWithAccessKey(region, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("discover-aliyun: NewClientWithAccessKey failed: %s", err)
	}

	l.Printf("[INFO] discover-aliyun: Filter instances with %s=%s", tagKey, tagValue)
	request := &ecs.DescribeInstancesRequest{
		Status: "Running",
		Tag: &[]ecs.DescribeInstancesTag{ecs.DescribeInstancesTag{
			Key:   tagKey,
			Value: tagValue,
		}},
	}
	if p.userAgent != "" {
		request.AppendUserAgent("go-discover", p.userAgent)
	}

	resp, err := svc.DescribeInstances(request)
	if err != nil {
		return nil, fmt.Errorf("discover-aliyun: DescribeInstances failed: %s", err)
	}

	l.Printf("[DEBUG] discover-aliyun: Found total %d instances", resp.TotalCount)

	var addrs []string
	for _, instanceAttributesType := range resp.Instances.Instance {
		switch instanceAttributesType.InstanceNetworkType {
		case "classic":
			for _, ipAddress := range instanceAttributesType.InnerIpAddress.IpAddress {
				l.Printf("[DEBUG] discover-aliyun: Instance %s has innner ip %s ", instanceAttributesType.InstanceId, ipAddress)
				addrs = append(addrs, ipAddress)
			}
		case "vpc":
			for _, ipAddress := range instanceAttributesType.VpcAttributes.PrivateIpAddress.IpAddress {
				l.Printf("[DEBUG] discover-aliyun: Instance %s has private ip %s ", instanceAttributesType.InstanceId, ipAddress)
				addrs = append(addrs, ipAddress)
			}
		}
	}

	l.Printf("[DEBUG] discover-aliyun: Found ip addresses: %v", addrs)
	return addrs, nil
}
