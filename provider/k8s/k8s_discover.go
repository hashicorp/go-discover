// Package k8s provides pod discovery for Kubernetes.
package k8s

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/go-homedir"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	// Register all known auth mechanisms since we might be authenticating
	// from anywhere.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

const (
	// AnnotationKeyPort is the annotation name of the field that specifies
	// the port name or number to append to the address.
	AnnotationKeyPort = "hashicorp.com/consul-auto-join-port"
)

type Provider struct{}

func (p *Provider) Help() string {
	return `Kubernetes (K8S):

    provider:         "k8s"
    kubeconfig:       Path to the kubeconfig file.
    namespace:        Namespace to search for pods (defaults to "default").
    label_selector:   Label selector value to filter pods.
    field_selector:   Field selector value to filter pods.

    The kubeconfig file value will be searched in the following locations:

     1. Use path from "kubeconfig" option if provided.
     2. Use path from KUBECONFIG environment variable.
     3. Use default path of $HOME/.kube/config
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "k8s" {
		return nil, fmt.Errorf("discover-k8s: invalid provider " + args["provider"])
	}

	// Get the kubeconfig path. If it is not set, set the default path.
	kubeconfig := args["kubeconfig"]
	if kubeconfig == "" {
		dir, err := homedir.Dir()
		if err != nil {
			return nil, fmt.Errorf("discover-k8s: error retrieving home directory: %s", err)
		}

		kubeconfig = filepath.Join(dir, ".kube", "config")
	}

	// Get the kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("discover-k8s: error loading kubeconfig: %s", err)
	}

	// Initialize the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("discover-k8s: error initializing k8s client: %s", err)
	}

	namespace := args["namespace"]
	if namespace == "" {
		namespace = "default"
	}

	// List all the pods based on the filters we requested
	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{
		LabelSelector: args["label_selector"],
		FieldSelector: args["field_selector"],
	})
	if err != nil {
		return nil, fmt.Errorf("discover-k8s: error listing pods: %s", err)
	}

	// Parse out the addresses from the pods
	var addrs []string
PodLoop:
	for _, pod := range pods.Items {
		if pod.Status.Phase != corev1.PodRunning {
			l.Printf("[DEBUG] discover-k8s: ignoring pod %q, not running: %q",
				pod.Name, pod.Status.Phase)
			continue
		}

		// If there is a Ready condition available, we need that to be true.
		// If no ready condition is set, then we accept this pod regardless.
		for _, condition := range pod.Status.Conditions {
			if condition.Type == corev1.PodReady && condition.Status != corev1.ConditionTrue {
				l.Printf("[DEBUG] discover-k8s: ignoring pod %q, not ready state", pod.Name)
				continue PodLoop
			}
		}

		// Get the IP address that we will join.
		addr := pod.Status.PodIP
		if addr == "" {
			// This can be empty according to the API docs, so we protect that.
			l.Printf("[DEBUG] discover-k8s: ignoring pod %q, PodIP is empty", pod.Name)
			continue
		}

		// We only use the port if it is specified as an annotation. The
		// annotation value can be a name or a number.
		if v := pod.Annotations[AnnotationKeyPort]; v != "" {
			port, err := podPort(&pod, v)
			if err != nil {
				l.Printf("[DEBUG] discover-k8s: ignoring pod %q, error retrieving port: %s",
					pod.Name, err)
				continue
			}

			addr = fmt.Sprintf("%s:%d", addr, port)
		}

		addrs = append(addrs, addr)
	}

	return addrs, nil
}

// podPort extracts the proper port for the address from the given pod
// for a non-empty annotation.
//
// Pre-condition: annotation is non-empty
func podPort(pod *corev1.Pod, annotation string) (int32, error) {
	// First look for a matching port matching the value of the annotation.
	for _, container := range pod.Spec.Containers {
		for _, portDef := range container.Ports {
			if portDef.Name == annotation {
				return portDef.ContainerPort, nil
			}
		}
	}

	// Otherwise assume that the port is a numeric value.
	v, err := strconv.ParseInt(annotation, 0, 32)
	return int32(v), err
}
