// Package aws provides node discovery for Amazon AWS.
package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ecs"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Provider struct{}

func (p *Provider) Help() string {
	return `Amazon AWS:

    provider:          "aws"
    region:            The AWS region. Default to region of instance.
    tag_key:           The tag key to filter on
    tag_value:         The tag value to filter on
    addr_type:         "private_v4", "public_v4" or "public_v6". Defaults to "private_v4".
    access_key_id:     The AWS access key to use
    secret_access_key: The AWS secret access key to use
    service:           The AWS service to filter. "ec2" or "ecs". Defaults to "ec2".
    endpoint:          The endpoint URL of EC2 to use. If not set the AWS client will set
                       this value, which defaults to the ec2 public dns for the specified
                       region

    The only required IAM permission is 'ec2:DescribeInstances'. If the Consul agent is
    running on AWS instance it is recommended you use an IAM role, otherwise it is
    recommended you make a dedicated IAM user and access key used only for auto-joining.
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "aws" {
		return nil, fmt.Errorf("discover-aws: invalid provider " + args["provider"])
	}

	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}

	region := args["region"]
	tagKey := args["tag_key"]
	tagValue := args["tag_value"]
	addrType := args["addr_type"]
	accessKey := args["access_key_id"]
	secretKey := args["secret_access_key"]
	sessionToken := args["session_token"]
	service := args["service"]
	endpoint := args["endpoint"]

	if service != "ec2" && service != "ecs" {
		l.Printf("[INFO] discover-aws: Service type %s is not supported. Valid values are {ec2,ecs}. Falling back to 'ec2'", service)
		service = "ec2"
	} else if service == "ecs" && addrType != "public_v4" {
		l.Printf("[INFO] discover-aws: Address Type %s is not supported for ECS. Valid values are {private_v4}. Falling back to 'private_v4'", addrType)
		addrType = "private_v4"
	}

	if addrType != "private_v4" && addrType != "public_v4" && addrType != "public_v6" {
		l.Printf("[INFO] discover-aws: Address type %s is not supported. Valid values are {private_v4,public_v4,public_v6}. Falling back to 'private_v4'", addrType)
		addrType = "private_v4"
	}

	if addrType == "" {
		l.Printf("[DEBUG] discover-aws: Address type not provided. Using 'private_v4'")
		addrType = "private_v4"
	}

	l.Printf("[DEBUG] discover-aws: Using region=%s tag_key=%s tag_value=%s addr_type=%s", region, tagKey, tagValue, addrType)
	if accessKey == "" && secretKey == "" {
		l.Printf("[DEBUG] discover-aws: No static credentials")
		l.Printf("[DEBUG] discover-aws: Using environment variables, shared credentials or instance role")
	} else {
		l.Printf("[DEBUG] discover-aws: Static credentials provided")
	}

	if region == "" {
		_, ecsEnabled := os.LookupEnv("ECS_CONTAINER_METADATA_URI_V4")
		if ecsEnabled != false {
			l.Printf("[INFO] discover-aws: Region not provided. Looking up region in ecs metadata...")
			region = os.Getenv("AWS_REGION")
		} else {
			l.Printf("[INFO] discover-aws: Region not provided. Looking up region in ec2 metadata...")
			ec2meta := ec2metadata.New(session.New())
			identity, err := ec2meta.GetInstanceIdentityDocument()
			if err != nil {
				return nil, fmt.Errorf("discover-aws: GetInstanceIdentityDocument failed: %s", err)
			}
			region = identity.Region
		}
	}
	l.Printf("[INFO] discover-aws: Region is %s", region)

	l.Printf("[DEBUG] discover-aws: Creating session...")
	config := aws.Config{
		Region: &region,
		Credentials: credentials.NewChainCredentials(
			[]credentials.Provider{
				&credentials.StaticProvider{
					Value: credentials.Value{
						AccessKeyID:     accessKey,
						SecretAccessKey: secretKey,
						SessionToken:    sessionToken,
					},
				},
				&credentials.EnvProvider{},
				&credentials.SharedCredentialsProvider{},
				defaults.RemoteCredProvider(*(defaults.Config()), defaults.Handlers()),
			}),
	}
	if endpoint != "" {
		l.Printf("[INFO] discover-aws: Endpoint is %s", endpoint)
		config.Endpoint = &endpoint
	}

	// Split here for ec2 vs ecs decision tree
	if service == "ecs" {
		svc := ecs.New(session.New(), &config)

		log.Printf("[INFO] discover-aws: Filter ECS instances with %s=%s", tagKey, tagValue)
		clusterArns := getEcsClusters(svc, aws.Int64(100), aws.String(""))

		var taskIps []string
		for _, clusterArn := range clusterArns {
			taskArns := getEcsTasks(svc, clusterArn, aws.Int64(100), aws.String(""))
			log.Printf("[DEBUG] discover-aws: Found %d ECS Tasks", len(taskArns))

			if len(taskArns) > 0 {
				taskIps = append(taskIps, getEcsTaskIps(svc, clusterArn, taskArns, &tagKey, &tagValue, &addrType)...)
				log.Printf("[DEBUG] discover-aws: Found %d ECS IPs", len(taskArns))
			}
		}
		log.Printf("[DEBUG] discover-aws: Discovered ECS Task IPs: %v", taskIps)
		return taskIps, nil
	}

	// When not using ECS continue with the default EC2 search

	svc := ec2.New(session.New(), &config)

	l.Printf("[INFO] discover-aws: Filter instances with %s=%s", tagKey, tagValue)
	resp, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("tag:" + tagKey),
				Values: []*string{aws.String(tagValue)},
			},
			&ec2.Filter{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("discover-aws: DescribeInstancesInput failed: %s", err)
	}

	l.Printf("[DEBUG] discover-aws: Found %d reservations", len(resp.Reservations))
	var addrs []string
	for _, r := range resp.Reservations {
		l.Printf("[DEBUG] discover-aws: Reservation %s has %d instances", *r.ReservationId, len(r.Instances))
		for _, inst := range r.Instances {
			id := *inst.InstanceId
			l.Printf("[DEBUG] discover-aws: Found instance %s", id)

			switch addrType {
			case "public_v6":
				l.Printf("[DEBUG] discover-aws: Instance %s has %d network interfaces", id, len(inst.NetworkInterfaces))

				for _, networkinterface := range inst.NetworkInterfaces {
					l.Printf("[DEBUG] discover-aws: Checking NetworInterfaceId %s on Instance %s", *networkinterface.NetworkInterfaceId, id)
					// Check if instance got any ipv6
					if networkinterface.Ipv6Addresses == nil {
						l.Printf("[DEBUG] discover-aws: Instance %s has no IPv6 on NetworkInterfaceId %s", id, *networkinterface.NetworkInterfaceId)
						continue
					}
					for _, ipv6address := range networkinterface.Ipv6Addresses {
						l.Printf("[INFO] discover-aws: Instance %s has IPv6 %s on NetworkInterfaceId %s", id, *ipv6address.Ipv6Address, *networkinterface.NetworkInterfaceId)
						addrs = append(addrs, *ipv6address.Ipv6Address)
					}
				}

			case "public_v4":
				if inst.PublicIpAddress == nil {
					l.Printf("[DEBUG] discover-aws: Instance %s has no public IPv4", id)
					continue
				}

				l.Printf("[INFO] discover-aws: Instance %s has public ip %s", id, *inst.PublicIpAddress)
				addrs = append(addrs, *inst.PublicIpAddress)

			default:
				// EC2-Classic don't have the PrivateIpAddress field
				if inst.PrivateIpAddress == nil {
					l.Printf("[DEBUG] discover-aws: Instance %s has no private ip", id)
					continue
				}

				l.Printf("[INFO] discover-aws: Instance %s has private ip %s", id, *inst.PrivateIpAddress)
				addrs = append(addrs, *inst.PrivateIpAddress)
			}
		}
	}

	l.Printf("[DEBUG] discover-aws: Found ip addresses: %v", addrs)
	return addrs, nil
}

func getEcsClusters(svc *ecs.ECS, maxResults *int64, nextToken *string) []*string {
	// List up to maxResults cluster arns
	clusterOutput, err := svc.ListClusters(&ecs.ListClustersInput{
		MaxResults: maxResults,
		NextToken:  nextToken,
	})

	if err != nil {
		log.Printf("discover-aws: ListClusters failed: %s", err)
	}

	clusterArns := clusterOutput.ClusterArns
	log.Printf("[INFO] discover-aws: Retrieved %d ECS Clusters", len(clusterArns))

	// If we have a next token string recursively append the resulting cluster arns
	if clusterOutput.NextToken != nil {
		clusterArns = append(clusterArns, getEcsClusters(svc, maxResults, clusterOutput.NextToken)...)
	}

	return clusterArns
}

func getEcsTasks(svc *ecs.ECS, clusterArn *string, maxResults *int64, nextToken *string) []*string {
	var taskArns []*string
	// List up to maxResults task arns
	taskOutput, err := svc.ListTasks(&ecs.ListTasksInput{
		Cluster:       clusterArn,
		DesiredStatus: aws.String("RUNNING"),
		MaxResults:    maxResults,
		NextToken:     nextToken,
	})

	if err != nil {
		log.Printf("discover-aws: ListTasks failed: %s", err)
	}

	taskArns = append(taskArns, taskOutput.TaskArns...)
	log.Printf("[INFO] discover-aws: Retrieved %d TaskArns", len(taskArns))

	if taskOutput.NextToken != nil {
		taskArns = append(taskArns, getEcsTasks(svc, clusterArn, maxResults, taskOutput.NextToken)...)
	}
	return taskArns
}

func getEcsTaskIps(svc *ecs.ECS, clusterArn *string, taskArns []*string, tagKey *string, tagValue *string, addrType *string) []string {
	// Describe all the tasks listed for this cluster
	taskDescriptions, err := svc.DescribeTasks(&ecs.DescribeTasksInput{
		Cluster: clusterArn,
		Include: []*string{aws.String(ecs.TaskFieldTags)},
		Tasks:   taskArns,
	})

	if err != nil {
		log.Printf("discover-aws: DescribeTasks failed: %s", err)
	}

	taskRequestFailures := taskDescriptions.Failures
	tasks := taskDescriptions.Tasks
	log.Printf("[INFO] discover-aws: Retrieved %d Task Descriptions and %d Failures", len(tasks), len(taskRequestFailures))

	// Filter tasks by Tag and Connectivity Status
	var ipList []string
	for _, taskDescription := range tasks {

		//log.Printf("[DEBUG] discover-aws: Parsing Task: %s", taskDescription)
		for _, tag := range taskDescription.Tags {
			//log.Printf("[DEBUG] discover-aws: Tag: %s", tag)

			if *tag.Key == *tagKey && *tag.Value == *tagValue {
				log.Printf("[DEBUG] discover-aws: Tag Match: %s : %s, desiredStatus: %s", *tag.Key, *tag.Value, *taskDescription.DesiredStatus)

				if *taskDescription.DesiredStatus == *aws.String("RUNNING") {
					log.Printf("[INFO] discover-aws: Found Running Instance: %s", *taskDescription.TaskArn)
					ip := getIpFromTaskDescription(taskDescription)

					if ip != nil {
						log.Printf("[DEBUG] discover-aws: Found Private IP: %s", *ip)
						ipList = append(ipList, *ip)
					}

				}

			}

		}
	}

	log.Printf("[INFO] discover-aws: Retrieved %d IPs from %d Tasks", len(ipList), len(taskArns))
	return ipList
}

func getIpFromTaskDescription(taskDesc *ecs.Task) *string {
	log.Printf("[DEBUG] discover-aws: Searching %d attachments for IPs", len(taskDesc.Attachments))
	for _, attachment := range taskDesc.Attachments {

		log.Printf("[DEBUG] discover-aws: Searching %d attachment details for IPs", len(attachment.Details))
		for _, detail := range attachment.Details {

			if *detail.Name == *aws.String("privateIPv4Address") {
				log.Printf("[DEBUG] discover-aws: Parsing Private IPv4: %s", *detail.Value)
				return detail.Value
			}

		}

	}
	return nil
}
