// Package aliyun provides node discovery for Aliyun.
package aliyun

import (
	"fmt"
	"io"
	"log"
	"os"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v6/client"
	"github.com/alibabacloud-go/tea/tea"
)

const (
	// ECSDefaultEndpoint is the default API endpoint of ECS services
	ECSDefaultEndpoint = "ecs.cn-hangzhou.aliyuncs.com"
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
		l = log.New(io.Discard, "", 0)
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

	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyID),
		AccessKeySecret: tea.String(accessKeySecret),
	}
	endpoint := os.Getenv("ECS_ENDPOINT")
	if endpoint == "" {
		endpoint = ECSDefaultEndpoint
	}
	config.SetEndpoint(endpoint)

	if p.userAgent != "" {
		config.SetUserAgent(p.userAgent)
	}

	svc, err := ecs20140526.NewClient(config)

	if err != nil {
		return nil, fmt.Errorf("discover-aliyun: NewClient failed: %s", err)
	}

	l.Printf("[INFO] discover-aliyun: Filter instances with %s=%s", tagKey, tagValue)
	tag0 := &ecs20140526.DescribeInstancesRequestTag{
		Key:   tea.String(tagKey),
		Value: tea.String(tagValue),
	}
	describeInstancesRequest := &ecs20140526.DescribeInstancesRequest{
		RegionId: tea.String(region),
		Status:   tea.String("Running"),
		Tag:      []*ecs20140526.DescribeInstancesRequestTag{tag0},
	}
	resp, err := svc.DescribeInstances(describeInstancesRequest)
	content := resp.Body

	if err != nil {
		return nil, fmt.Errorf("discover-aliyun: DescribeInstances failed: %s", err)
	}

	l.Printf("[DEBUG] discover-aliyun: Found total %d instances", tea.Int32Value(content.TotalCount))

	var addrs []string
	for _, instanceAttributesType := range content.Instances.Instance {
		networkType := tea.StringValue(instanceAttributesType.InstanceNetworkType)
		switch networkType {
		case "classic":
			for _, ipAddress := range instanceAttributesType.InnerIpAddress.IpAddress {
				l.Printf("[DEBUG] discover-aliyun: Instance %s has innner ip %s ", tea.StringValue(instanceAttributesType.InstanceId), tea.StringValue(ipAddress))
				addrs = append(addrs, tea.StringValue(ipAddress))
			}
		case "vpc":
			for _, ipAddress := range instanceAttributesType.VpcAttributes.PrivateIpAddress.IpAddress {
				l.Printf("[DEBUG] discover-aliyun: Instance %s has private ip %s ", tea.StringValue(instanceAttributesType.InstanceId), tea.StringValue(ipAddress))
				addrs = append(addrs, tea.StringValue(ipAddress))
			}
		}
	}

	l.Printf("[DEBUG] discover-aliyun: Found ip addresses: %v", addrs)
	return addrs, nil
}
