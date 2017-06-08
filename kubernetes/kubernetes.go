package kubernetes

import (
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
)

// AllNamespaces returns the string which represents all namespaces
func AllNamespaces() string {
	return v1.NamespaceAll
}

// DefaultConfigFile returns the default kubeconfig file path
func DefaultConfigFile() string {
	return clientcmd.RecommendedHomeFile
}

// DefaultNamespace returns the default namespace
func DefaultNamespace() string {
	return v1.NamespaceDefault
}

// IsPodFailed returns whether the given Pod finished unsuccessfully or not
func IsPodFailed(pod v1.Pod) bool {
	return pod.Status.Phase == v1.PodFailed
}

// IsPodSucceeded returns whether the given Pod finished successfully or not
func IsPodSucceeded(pod v1.Pod) bool {
	return pod.Status.Phase == v1.PodSucceeded
}
