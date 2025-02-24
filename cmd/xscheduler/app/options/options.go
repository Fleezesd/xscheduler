package options

import (
	"time"

	xschedulerconfig "github.com/fleezesd/xscheduler/pkg/apis/config"
	xschedulerscheme "github.com/fleezesd/xscheduler/pkg/apis/config/scheme"
	"github.com/fleezesd/xscheduler/pkg/apis/config/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	apiserveroptions "k8s.io/apiserver/pkg/server/options"
	clientset "k8s.io/client-go/kubernetes"
	cliflag "k8s.io/component-base/cli/flag"
	componentbaseconfig "k8s.io/component-base/config"
	componentbaseoptions "k8s.io/component-base/config/options"
	"k8s.io/component-base/logs"
	logsapi "k8s.io/component-base/logs/api/v1"
	"k8s.io/component-base/metrics"
	"k8s.io/klog/v2"
)

func newDefaultComponentConfig() (*xschedulerconfig.XschedulerConfiguration, error) {
	versionedCfg := v1alpha1.XschedulerConfiguration{}
	xschedulerscheme.Scheme.Default(&versionedCfg)
	cfg := xschedulerconfig.XschedulerConfiguration{}
	// todo: internal convert & default codegen
	if err := xschedulerscheme.Scheme.Convert(&versionedCfg, &cfg, nil); err != nil {
		return nil, err
	}
	return &cfg, nil
}

type Options struct {
	ComponentConfig *xschedulerconfig.XschedulerConfiguration

	// SecureServing is the main context that defines what certificates to use for serving.
	SecureServing           *apiserveroptions.SecureServingOptionsWithLoopback
	CombinedInsecureServing CombinedInsecureServingOptions
	Metrics                 *metrics.Options
	Logs                    *logs.Options

	// ConfigFile is the location of the scheduler server's configuration file.s
	ConfigFile string
	// WriteConfigTo is the path where the default configuration will be written.
	WriteConfigTo string

	// Flags stores the parsed CLI flags
	Flags *cliflag.NamedFlagSets

	LeaderElection *componentbaseconfig.LeaderElectionConfiguration
	Client         clientset.Interface
}

// todo: init usage options
func NewOptions() *Options {
	cfg, err := newDefaultComponentConfig()
	if err != nil {
		klog.Fatalf("Failed to new default component config: %v", err)
	}
	o := &Options{
		ComponentConfig: cfg,
		SecureServing:   apiserveroptions.NewSecureServingOptions().WithLoopback(),
		CombinedInsecureServing: CombinedInsecureServingOptions{
			Healthz: &apiserveroptions.SecureServingOptions{
				BindNetwork: "tcp",
			},
			Metrics: &apiserveroptions.SecureServingOptions{
				BindNetwork: "tcp",
			},
		},
		Metrics: metrics.NewOptions(),
		Logs:    logs.NewOptions(),
		LeaderElection: &componentbaseconfig.LeaderElectionConfiguration{
			LeaderElect: true,
			LeaseDuration: metav1.Duration{
				Duration: 15 * time.Second,
			},
			RenewDeadline: metav1.Duration{
				Duration: 10 * time.Second,
			},
			RetryPeriod: metav1.Duration{
				Duration: 2 * time.Second,
			},
			ResourceLock:      "leases",
			ResourceName:      "fleezesd-xscheduler",
			ResourceNamespace: "xscheduler-system",
		},
	}
	o.SecureServing.BindPort = xschedulerconfig.DefaultXschedulerPort
	return o
}

// initFlags
func (o *Options) initFlags() {
	if o.Flags != nil {
		return
	}

	nfs := cliflag.NamedFlagSets{}
	fs := nfs.FlagSet("misc")
	fs.StringVar(&o.ConfigFile, "config", o.ConfigFile, "The path to the configuration file.")
	fs.StringVar(&o.WriteConfigTo, "write-config-to", o.WriteConfigTo, "If set, write the configuration values to this file and exit.")

	o.SecureServing.AddFlags(nfs.FlagSet("secure serving"))
	o.CombinedInsecureServing.AddFlags(nfs.FlagSet("insecure serving"))
	// leader election options
	componentbaseoptions.BindLeaderElectionFlags(o.LeaderElection, nfs.FlagSet("leader election"))
	// todo: feature gate flags
	o.Metrics.AddFlags(nfs.FlagSet("metrics"))
	logsapi.AddFlags(o.Logs, nfs.FlagSet("logs"))
	o.Flags = &nfs
}
