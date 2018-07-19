// Package awsecs provides node discovery for Amazon ECS AWS.
package awsecs

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// describeInstances returns data about the provded EC2 instance IDs
func (p *Provider) describeEC2Instances(svc *ec2.EC2, instanceContainerMap []ec2InstanceIdContainerARNMap, l *log.Logger) ([]*ec2.Reservation, error) {

	i := []*string{}
	for _, v := range instanceContainerMap {
		i = append(i, aws.String(v.EC2InstanceId))
	}

	l.Printf("[INFO] discover-awsecs: Describe EC2 instances")
	resp, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: i,
	})
	if err != nil {
		return nil, fmt.Errorf("discover-awsecs: DescribeInstances failed: %s", err)
	}

	return resp.Reservations, nil
}
