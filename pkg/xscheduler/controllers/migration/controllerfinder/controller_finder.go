package controllerfinder

import (
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	scaleclient "k8s.io/client-go/scale"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type ScaleAndSelector struct {
	ControllerReference
	// controller.spec.Replicas
	Scale int32
	// kruise statefulSet.spec.ReservedOrdinals
	ReserveOrdinals []int
	// controller.spec.Selector
	Selector *metav1.LabelSelector
	// metadata
	Metadata metav1.ObjectMeta
}

type ControllerReference struct {
	// API version of the referent.
	APIVersion string `json:"apiVersion" protobuf:"bytes,5,opt,name=apiVersion"`
	// Kind of the referent.
	Kind string `json:"kind" protobuf:"bytes,1,opt,name=kind"`
	// Name of the referent.
	Name string `json:"name" protobuf:"bytes,3,opt,name=name"`
	// UID of the referent.
	UID types.UID `json:"uid" protobuf:"bytes,4,opt,name=uid,casttype=k8s.io/apimachinery/pkg/types.UID"`
}

type ControllerFinder struct {
	client.Client

	mapper          meta.RESTMapper
	scaleNamespacer scaleclient.ScalesGetter
	discoveryClient discovery.DiscoveryInterface
}

type Interface interface {
	GetPodsForRef(ownerReference *metav1.OwnerReference, ns string, labelSelector *metav1.LabelSelector, active bool) ([]*corev1.Pod, int32, error)
	GetExpectedScaleForPod(pods *corev1.Pod) (int32, error)
	ListPodsByWorkloads(workloadUIDs []types.UID, ns string, labelSelector *metav1.LabelSelector, active bool) ([]*corev1.Pod, error)
}

func NewControllerFinder(mgr manager.Manager) (*ControllerFinder, error) {
	finder := &ControllerFinder{
		Client: mgr.GetClient(),
		mapper: mgr.GetRESTMapper(),
	}
	cfg := mgr.GetConfig()
	if lo.IsNil(cfg.GroupVersion) {
		cfg.GroupVersion = &schema.GroupVersion{}
	}
	codecs := serializer.NewCodecFactory(mgr.GetScheme())
	// NegotiatedSerializer is used for obtaining encoders and decoders for multiple supported media types.
	cfg.NegotiatedSerializer = codecs.WithoutConversion()
	restClient, err := rest.RESTClientFor(cfg)
	if err != nil {
		return nil, err
	}
	k8sClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	finder.discoveryClient = k8sClient.Discovery()
	scaleKindResolver := scaleclient.NewDiscoveryScaleKindResolver(finder.discoveryClient)
	// NewDiscoveryScaleKindResolver creates a new ScaleKindResolver which uses information from the given
	// disovery client to resolve the correct Scale GroupVersionKind for different resources.
	finder.scaleNamespacer = scaleclient.New(restClient, finder.mapper, dynamic.LegacyAPIPathResolverFunc, scaleKindResolver)
	return finder, nil
}

func (r *ControllerFinder) GetExpectedScaleForPod(pod *corev1.Pod) (int32, error) {
	if lo.IsNil(pod) {
		return 0, nil
	}
	ref := metav1.GetControllerOf(pod)
	if lo.IsNil(ref) {
		return 0, nil
	}
	workload, err := r.GetScaleAndSelectorForRef(ref.APIVersion, ref.Kind, pod.Namespace, ref.Name, ref.UID)
	if err != nil && !errors.IsNotFound(err) {
		return 0, err
	}
	if workload != nil && workload.Metadata.DeletionTimestamp.IsZero() {
		return workload.Scale, nil
	}
	return 0, nil
}

func (r *ControllerFinder) GetScaleAndSelectorForRef(apiVersion, kind, namespace,
	name string, uid types.UID) (*ScaleAndSelector, error) {
	targetRef := ControllerReference{
		APIVersion: apiVersion,
		Kind:       kind,
		Name:       name,
		UID:        uid,
	}

	for _, finder := range r.Finders() {
		scale, err := finder(targetRef, namespace)
		if scale != nil || err != nil {
			return scale, err
		}
	}
	return nil, nil
}

func (r *ControllerFinder) Finders() []PodControllerFinder {
	// todo: configure pod controller finder
	return []PodControllerFinder{}
}
