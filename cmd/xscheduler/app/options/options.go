package options

import (
	clientset "k8s.io/client-go/kubernetes"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"
	"k8s.io/component-base/metrics"
)

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
