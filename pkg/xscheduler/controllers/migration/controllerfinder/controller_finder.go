package controllerfinder

import (
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
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
	finder.scaleNamespacer = scaleclient.New(restClient, finder.mapper, dynamic.LegacyAPIPathResolverFunc, scaleKindResolver)
	return finder, nil
}
