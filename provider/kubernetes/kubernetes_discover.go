package kubernetes

import (
	"fmt"
	"io/ioutil"
	"log"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Provider struct{}

func (p *Provider) Help() string {
	return `Kubernetes:

    provider:     "kubernetes"
    namespace:    The Kubernetes namespace to filter on
    label_key:    The Kubernetes label to filter on
    label_value:  The Kubernetes label value to filter on
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "kubernetes" {
		return nil, fmt.Errorf("discover-kubernetes: invalid provider " + args["provider"])
	}

	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}

	namespace := args["namespace"]
	labelKey := args["label_key"]
	labelValue := args["label_value"]

	if namespace == "" {
		l.Printf("[DEBUG] discover-kubernetes: Namespace type not provided. Using 'default'")
		namespace = "default"
	}

	l.Printf("[DEBUG] discover-kubernetes: Using namespace=%s label_key=%s label_value=%s", namespace, labelKey, labelValue)

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("discover-kubernetes: %s", err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("discover-kubernetes: %s", err)
	}

	l.Printf("[INFO] discover-aws: Select instances with %s=%s", labelKey, labelValue)
	var podIpAddresses []string
	labelsSet := labels.Set(map[string]string{labelKey: labelValue})
	pods, err := client.Core().Pods(namespace).List(meta_v1.ListOptions{LabelSelector: labelsSet.AsSelector().String()})
	if err != nil {
		return nil, fmt.Errorf("discover-kubernetes: listing pods with label %s=%s in %s: %v", labelKey, labelValue, namespace, err)
	}
	for _, v := range pods.Items {
		podIpAddresses = append(podIpAddresses, v.Status.PodIP)
	}

	return podIpAddresses, nil
}
