package options

import (
	"github.com/fleezesd/xscheduler/pkg/apis/config/v1alpha1"
	clientset "k8s.io/client-go/kubernetes"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"
	"k8s.io/component-base/metrics"
)

func newDefaultComponentConfig() {
	_ = v1alpha1.XschedulerConfiguration{}

}

type Options struct {
	Metrics *metrics.Options
	Logs    *logs.Options

	// Flags stores the parsed CLI flags
	Flags *cliflag.NamedFlagSets

	Client clientset.Interface
}

// todo: init usage options
func NewOptions() *Options {
	o := &Options{}
	return o
}
