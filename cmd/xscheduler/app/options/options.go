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
	// 2. todo: internal convert & default codegen
	if err := xschedulerscheme.Scheme.Convert(&versionedCfg, &cfg, nil); err != nil {
		return nil, err
	}
	return &cfg, nil
}

type Options struct {
	ComponentConfig *xschedulerconfig.XschedulerConfiguration

	// SecureServing is the main context that defines what certificates to use for serving.
	SecureServing  *apiserveroptions.SecureServingOptionsWithLoopback
	Authentication *apiserveroptions.DelegatingAuthenticationOptions
	Authorization  *apiserveroptions.DelegatingAuthorizationOptions
	Metrics        *metrics.Options
	Logs           *logs.Options
	LeaderElection *componentbaseconfig.LeaderElectionConfiguration

	// ConfigFile is the location of the scheduler server's configuration file.s
	ConfigFile string
	// WriteConfigTo is the path where the default configuration will be written.
	WriteConfigTo string

	// Flags stores the parsed CLI flags
	Flags  *cliflag.NamedFlagSets
	Client clientset.Interface
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
		Authentication:  apiserveroptions.NewDelegatingAuthenticationOptions(),
		Authorization:   apiserveroptions.NewDelegatingAuthorizationOptions(),
		LeaderElection: &componentbaseconfig.LeaderElectionConfiguration{
			LeaderElect:       true,
			LeaseDuration:     metav1.Duration{Duration: 15 * time.Second},
			RenewDeadline:     metav1.Duration{Duration: 10 * time.Second},
			RetryPeriod:       metav1.Duration{Duration: 2 * time.Second},
			ResourceLock:      "leases",
			ResourceName:      "xscheduler",
			ResourceNamespace: "xscheduler-system",
		},
		Metrics: metrics.NewOptions(),
		Logs:    logs.NewOptions(),
	}

	o.Authentication.TolerateInClusterLookupFailure = true
	o.Authentication.RemoteKubeConfigFileOptional = true
	o.Authorization.RemoteKubeConfigFileOptional = true

	// Set the PairName but leave certificate directory blank to generate in-memory by default
	o.SecureServing.ServerCert.CertDirectory = ""
	o.SecureServing.ServerCert.PairName = "xscheduler"
	o.SecureServing.BindPort = xschedulerconfig.DefaultXschedulerPort
	o.initFlags()
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
	o.Authentication.AddFlags(nfs.FlagSet("authentication"))
	o.Authorization.AddFlags(nfs.FlagSet("authorization"))
	// leader election options
	componentbaseoptions.BindLeaderElectionFlags(o.LeaderElection, nfs.FlagSet("leader election"))
	// todo: feature gate flags
	o.Metrics.AddFlags(nfs.FlagSet("metrics"))
	logsapi.AddFlags(o.Logs, nfs.FlagSet("logs"))

	o.Flags = &nfs
}
