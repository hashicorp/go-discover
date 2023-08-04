// Package awsecs provides node discovery for Amazon ECS AWS.
package awsecs

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// describeTasks gets a list of ECS instances which match the service name and returns data about the tasks
//
// List functionality provided by listTasks()
func (p *Provider) describeTasks(svc *ecs.ECS, clusterName string, serviceName string, l *log.Logger) ([]*ecs.Task, error) {

	taskArns, err := p.listTasks(svc, clusterName, serviceName, l)
	if err != nil {
		return nil, err
	}

	l.Printf("[INFO] discover-awsecs: Describe tasks")
	resp, err := svc.DescribeTasks(&ecs.DescribeTasksInput{
		Cluster: aws.String(clusterName),
		Tasks:   taskArns,
	})

	if err != nil {
		return nil, fmt.Errorf("discover-awsecs: DescribeTasks failed: %s", err)
	}

	return resp.Tasks, nil
}

// describeTasks gets a list of ECS instances which match the service name
func (p *Provider) listTasks(svc *ecs.ECS, clusterName string, serviceName string, l *log.Logger) ([]*string, error) {

	l.Printf("[INFO] discover-awsecs: List tasks of service %s in cluster %s", serviceName, clusterName)
	resp, err := svc.ListTasks(&ecs.ListTasksInput{
		Cluster:     aws.String(clusterName),
		ServiceName: aws.String(serviceName),
	})
	if err != nil {
		return nil, fmt.Errorf("discover-awsecs: ListTasks failed: %s", err)
	}
	l.Printf("[DEBUG] discover-awsecs: Found tasks")

	return resp.TaskArns, nil
}

// describeContainerInstances gets a list of ECS instances which match the service name and returns data about the instance
//
// List functionality provided by listContainerInstances()
func (p *Provider) describeContainerInstances(svc *ecs.ECS, clusterName string, l *log.Logger) ([]*ecs.ContainerInstance, error) {

	list, err := p.listContainerInstances(svc, clusterName, l)
	if err != nil {
		return nil, err
	}

	l.Printf("[INFO] discover-awsecs: Describe container instances in cluster %s", clusterName)
	resp, err := svc.DescribeContainerInstances(&ecs.DescribeContainerInstancesInput{
		Cluster:            aws.String(clusterName),
		ContainerInstances: list,
	})
	if err != nil {
		return nil, fmt.Errorf("discover-awsecs: DescribeContainerInstances failed: %s", err)
	}

	return resp.ContainerInstances, nil

}

func (p *Provider) listContainerInstances(svc *ecs.ECS, clusterName string, l *log.Logger) ([]*string, error) {

	l.Printf("[INFO] discover-awsecs: List container instances in cluster %s", clusterName)
	resp, err := svc.ListContainerInstances(&ecs.ListContainerInstancesInput{
		Cluster: aws.String(clusterName),
	})
	if err != nil {
		return nil, fmt.Errorf("discover-awsecs: ListContainerInstances failed: %s", err)
	}

	return resp.ContainerInstanceArns, nil
}
