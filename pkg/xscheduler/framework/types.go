package framework

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	clientset "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/events"
)

type Handle interface {
	PluginsRunner

	// ClientSet returns a kubernetes clientSet
	ClientSet() clientset.Interface

	// KubeConfig returns the raw kube config.
	KubeConfig() *restclient.Config

	// EventRecorder returns an event recorder.
	EventRecorder() events.EventRecorder

	// Evictor
	Evictor() Evictor

	GetPodsAssignedToNodeFunc() GetPodsAssignedToNodeFunc

	SharedInformerFactory() informers.SharedInformerFactory
}

type Evictor interface {
	// Filter checks if a pod can be evicted
	Filter(pod *corev1.Pod) bool

	// PreEvictionFilter checks if pod can be evicted right before eviction
	PreEvictionFilter(pod *corev1.Pod) bool

	// Evict evicts a pod (no pre-check performed)
	Evict(ctx context.Context, pod *corev1.Pod, evictOption EvictOptions) bool
}

type EvictOptions struct {
	// PluginName represents the initiator of the eviction operation
	PluginName string
	// Reason allows for passing details about the specific eviction for logging.
	Reason string
	// DeleteOptions holds the arguments used to delete
	DeleteOptions *metav1.DeleteOptions
}

// GetPodsAssignedToNodeFunc is a function which accept a node name and a pod filter function
// as input and returns the pods that assigned to the node.
type GetPodsAssignedToNodeFunc func(string, FilterFunc) ([]*corev1.Pod, error)

// FilterFunc is a filter for a pod.
type FilterFunc func(*corev1.Pod) bool
