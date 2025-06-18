// Package awsecs provides node discovery for Amazon ECS AWS.
package awsecs

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type Provider struct{}

type ec2InstanceIdContainerARNMap struct {
	EC2InstanceId        string
	ContainerInstanceARN string
}

func (p *Provider) Help() string {
	return `Amazon AWS ECS:

	provider:          "awsecs"
	region:            The AWS region. Default to region of instance.
	addr_type:         "private_v4", "public_v4" or "public_v6". Defaults to "private_v4".
	service_port:	   The port that the container exposes for the service
	cluster_name:      The name of the cluster where the service is deployed
	service_name:      The name of the service
	container_name:    The name of the running contasiner within the service
	access_key_id:     The AWS access key to use
	secret_access_key: The AWS secret access key to use
`
}

// Addrs attempts to find ECS containers matching the arguments, and return their address and container port
//
// Given the following scenario:
// A service running a container exposing port 80 is deployed 3 times on 2 container instances:
//  - Instance A: IP address 10.0.0.1 container 80 exposed dynamically on 32001
//  - Instance A: IP address 10.0.0.1 container 80 exposed dynamically on 32002
//  - Instance B: IP address 10.0.0.2 container 80 exposed dynamically on 32001
// This will return
// [
//    10.0.0.1:32001,
//    10.0.0.1:32002,
//    10.0.0.2:32001,
// ]
func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "awsecs" {
		return nil, fmt.Errorf("discover-awsecs: invalid provider " + args["provider"])
	}

	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}

	region := args["region"]
	addrType := args["addr_type"]
	clusterName := args["cluster_name"]
	serviceName := args["service_name"]
	servicePort := args["service_port"]
	containerName := args["container_name"]
	accessKey := args["access_key_id"]
	secretKey := args["secret_access_key"]

	if addrType != "private_v4" && addrType != "public_v4" && addrType != "public_v6" {
		l.Printf("[INFO] discover-awsecs: Address type %s is not supported. Valid values are {private_v4,public_v4,public_v6}. Falling back to 'private_v4'", addrType)
		addrType = "private_v4"
	}

	if clusterName == "" || serviceName == "" || servicePort == "" || containerName == "" {
		return nil, fmt.Errorf("discover-awsecs: cluster_name, service_name, service_port and container_name are all required")
	}

	// Service port should be an int64
	servicePortInt, err := strconv.ParseInt(servicePort, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("discover-awsecs: service_port must be an integer")
	}

	log.Printf("[DEBUG] discover-awsecs: Using region=%s cluster_name=%s service_name=%s container_name=%s addr_type=%s", region, clusterName, serviceName, containerName, addrType)
	if accessKey == "" && secretKey == "" {
		log.Printf("[DEBUG] discover-awsecs: No static credentials")
		log.Printf("[DEBUG] discover-awsecs: Using environment variables, shared credentials or instance role")
	} else {
		log.Printf("[DEBUG] discover-awsecs: Static credentials provided")
	}

	if region == "" {
		l.Printf("[INFO] discover-aws: Region not provided. Looking up region in metadata...")
		ec2meta := ec2metadata.New(session.New())
		identity, err := ec2meta.GetInstanceIdentityDocument()
		if err != nil {
			return nil, fmt.Errorf("discover-awsecs: GetInstanceIdentityDocument failed: %s", err)
		}
		region = identity.Region
	}
	l.Printf("[INFO] discover-awsecs: Region is %s", region)

	l.Printf("[DEBUG] discover-awsecs: Creating session...")
	sess, conf := p.getAWSSessionAndConfig(region, accessKey, secretKey)

	l.Printf("[DEBUG] discover-awsecs: Creating ECS service...")
	ecssvc := ecs.New(sess, conf)

	l.Printf("[DEBUG] discover-awsecs: Creating EC2 service...")
	ec2svc := ec2.New(sess, conf)

	// Get container instances from AWS
	containerInstances, err := p.describeContainerInstances(ecssvc, clusterName, l)
	if err != nil {
		return nil, err
	}

	// Map container instances ARN to EC2 instance ID
	containerArnEc2IDMap := p.mapContainerInstances(containerInstances)

	// Get EC2 instance data for the container instances
	ec2Instances, err := p.describeEC2Instances(ec2svc, containerArnEc2IDMap, l)
	if err != nil {
		return nil, err
	}

	// Parse the EC2 instances and associate their IP address with the container ARN
	containerAddrs := p.getContainerAddresses(ec2Instances, containerArnEc2IDMap, addrType, l)

	// Get a list of matching tasks
	tasks, err := p.describeTasks(ecssvc, clusterName, serviceName, l)
	if err != nil {
		return nil, err
	}

	var addrs []string

	// This will add an address to the response list if (a) the container name matches, (b) the container
	// exposes the defined port and (c) the host as an address type that matches addr_type
	for _, t := range tasks {
		for _, c := range t.Containers {
			if *c.Name != containerName {
				continue
			}
			for _, nb := range c.NetworkBindings {
				if *nb.ContainerPort == servicePortInt && *nb.Protocol == "tcp" { // default consul HTTP port
					if containerAddrs[*t.ContainerInstanceArn] == "" {
						l.Printf("[INFO] discover-awsecs: Task %s is running but container instance %s has no %s IP address", *t.TaskArn, *t.ContainerInstanceArn, addrType)
						continue
					}
					addr := fmt.Sprintf("%s:%d", containerAddrs[*t.ContainerInstanceArn], *nb.HostPort)
					if addrType == "public_v6" {
						// IPv6 should be wrapped before adding port
						pi := strings.LastIndex(addr, ":")
						addr = fmt.Sprintf("[%s]:%s", addr[:pi], addr[pi+1:])
					}
					addrs = append(addrs, addr)
				}
			}
		}
	}

	l.Printf("[DEBUG] discover-aws: Found ip addresses: %v", addrs)
	return addrs, nil
}

// getContainerAddresses iterates over a list of EC2 instances to return IP addresses
//
// This method iterates over a given list of EC2 instances and inspects the network properties, looking for a matching
// address defined by `addrType`. If a container has an IP address matching the type, then a lookup is done on the passed
// instanceContainerMap map, which is a simple map of EC2 instance IDs to ECS container instance ARNs.
//
// The returned map is a definition of: [ContainerInstanceARN = IPAddress (of type `addrType`)]
//
// Main logic copied from github.com/hashicorp/go-discover/provider/aws/aws_discover.go
func (p *Provider) getContainerAddresses(reservations []*ec2.Reservation, instanceContainerMap []ec2InstanceIdContainerARNMap, addrType string, l *log.Logger) map[string]string {

	instancesAddresses := map[string]string{}

	for _, r := range reservations {
		l.Printf("[DEBUG] discover-awsecs: Reservation %s has %d instances", *r.ReservationId, len(r.Instances))
		for _, inst := range r.Instances {
			id := *inst.InstanceId
			l.Printf("[DEBUG] discover-awsecs: Found instance %s", id)

			switch addrType {
			case "public_v6":
				l.Printf("[DEBUG] discover-awsecs: Instance %s has %d network interfaces", id, len(inst.NetworkInterfaces))

				for _, networkinterface := range inst.NetworkInterfaces {
					l.Printf("[DEBUG] discover-awsecs: Checking NetworInterfaceId %s on Instance %s", *networkinterface.NetworkInterfaceId, id)
					// Check if instance got any ipv6
					if networkinterface.Ipv6Addresses == nil {
						l.Printf("[DEBUG] discover-awsecs: Instance %s has no IPv6 on NetworkInterfaceId %s", id, *networkinterface.NetworkInterfaceId)
						continue
					}
					for _, ipv6address := range networkinterface.Ipv6Addresses {
						l.Printf("[INFO] discover-awsecs: Instance %s has IPv6 %s on NetworkInterfaceId %s", id, *ipv6address.Ipv6Address, *networkinterface.NetworkInterfaceId)
						instancesAddresses[p.getContainerARNFromEc2InstanceId(*inst.InstanceId, instanceContainerMap)] = *ipv6address.Ipv6Address
					}
				}

			case "public_v4":
				if inst.PublicIpAddress == nil {
					l.Printf("[DEBUG] discover-awsecs: Instance %s has no public IPv4", id)
					continue
				}

				l.Printf("[INFO] discover-awsecs: Instance %s has public ip %s", id, *inst.PublicIpAddress)
				instancesAddresses[p.getContainerARNFromEc2InstanceId(*inst.InstanceId, instanceContainerMap)] = *inst.PublicIpAddress

			default:
				// EC2-Classic don't have the PrivateIpAddress field
				if inst.PrivateIpAddress == nil {
					l.Printf("[DEBUG] discover-awsecs: Instance %s has no private ip", id)
					continue
				}

				l.Printf("[INFO] discover-awsecs: Instance %s has private ip %s", id, *inst.PrivateIpAddress)
				instancesAddresses[p.getContainerARNFromEc2InstanceId(*inst.InstanceId, instanceContainerMap)] = *inst.PrivateIpAddress
			}
		}
	}

	return instancesAddresses
}

func (p *Provider) getContainerARNFromEc2InstanceId(ec2InstanceID string, instanceContainerMap []ec2InstanceIdContainerARNMap) string {
	for _, v := range instanceContainerMap {
		if v.EC2InstanceId == ec2InstanceID {
			return v.ContainerInstanceARN
		}
	}
	return ""
}

// mapContainerInstances iterattes over container instances and creates a correlation map between the container instance
// ARN and the EC2 instance ID
func (p *Provider) mapContainerInstances(ci []*ecs.ContainerInstance) []ec2InstanceIdContainerARNMap {
	containerArnEc2IDMap := []ec2InstanceIdContainerARNMap{}

	for _, ci := range ci {
		containerArnEc2IDMap = append(containerArnEc2IDMap, ec2InstanceIdContainerARNMap{
			EC2InstanceId:        *ci.Ec2InstanceId,
			ContainerInstanceARN: *ci.ContainerInstanceArn,
		})
	}

	return containerArnEc2IDMap
}

// mapContainerInstances iterattes over container instances and creates a correlation map between the container instance
// ARN and the EC2 instance ID
func (p *Provider) getAWSSessionAndConfig(region, accessKey, secretKey string) (*session.Session, *aws.Config) {
	sess := session.New()
	conf := &aws.Config{
		Region: &region,
		Credentials: credentials.NewChainCredentials(
			[]credentials.Provider{
				&credentials.StaticProvider{
					Value: credentials.Value{
						AccessKeyID:     accessKey,
						SecretAccessKey: secretKey,
					},
				},
				&credentials.EnvProvider{},
				&credentials.SharedCredentialsProvider{},
				defaults.RemoteCredProvider(*(defaults.Config()), defaults.Handlers()),
			},
		),
	}

	return sess, conf
}
