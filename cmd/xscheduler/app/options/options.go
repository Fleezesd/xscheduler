package options

import (
	"context"
	"fmt"
	"net"
	"time"

	xschedulerappconfig "github.com/fleezesd/xscheduler/cmd/xscheduler/app/config"
	xschedulerconfig "github.com/fleezesd/xscheduler/pkg/xscheduler/apis/config"
	xschedulerscheme "github.com/fleezesd/xscheduler/pkg/xscheduler/apis/config/scheme"
	"github.com/fleezesd/xscheduler/pkg/xscheduler/apis/config/v1alpha1"
	"github.com/fleezesd/xscheduler/pkg/xscheduler/apis/config/validation"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic/dynamicinformer"
	netutils "k8s.io/utils/net"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	apiserveroptions "k8s.io/apiserver/pkg/server/options"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/client-go/dynamic"
	clientset "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/events"
	"k8s.io/client-go/tools/leaderelection"
	cliflag "k8s.io/component-base/cli/flag"
	componentbaseconfig "k8s.io/component-base/config"
	componentbaseoptions "k8s.io/component-base/config/options"
	"k8s.io/component-base/logs"
	logsapi "k8s.io/component-base/logs/api/v1"
	"k8s.io/component-base/metrics"
	zpagesfeatures "k8s.io/component-base/zpages/features"
	"k8s.io/component-base/zpages/flagz"
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
	Master        string

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
	fs.StringVar(&o.Master, "master", o.Master, "The address of the Kubernetes API server (overrides any value in kubeconfig)")

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

// Validate validates all the required options
func (o *Options) Validate() []error {
	var errs []error
	if err := validation.ValidateXschedulerConfiguration(o.ComponentConfig); err != nil {
		errs = append(errs, err.Errors()...)
	}
	errs = append(errs, o.SecureServing.Validate()...)
	// todo: check if or not need insecure server validate
	errs = append(errs, o.Metrics.Validate()...)
	return errs
}

// Config returns a scheduler config object
func (o *Options) Config(ctx context.Context) (*xschedulerappconfig.Config, error) {
	// logger := klog.FromContext(ctx)
	if o.SecureServing != nil {
		if err := o.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{netutils.ParseIPSloppy("127.0.0.1")}); err != nil {
			return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
		}
	}

	c := &xschedulerappconfig.Config{}
	if err := o.ApplyTo(c); err != nil {
		return nil, err
	}

	// Prepare kube clients.
	client, eventClient, err := createClients(c.KubeConfig)
	if err != nil {
		return nil, err
	}
	c.EventBroadcaster = events.NewEventBroadcasterAdapter(eventClient)

	var leaderElectionConfig *leaderelection.LeaderElectionConfig
	// todo: make leader election if enabled later

	mgrKubeConfig := *c.KubeConfig
	mgrKubeConfig.ContentType = ""
	mgrKubeConfig.AcceptContentTypes = ""
	mgr, err := ctrl.NewManager(&mgrKubeConfig, ctrl.Options{
		Scheme:                 scheme,
		Metrics:                metricsserver.Options{BindAddress: "0"},
		HealthProbeBindAddress: "0",
		LeaderElection:         false,
		Cache: cache.Options{
			SyncPeriod: nil,
		},
		// todo: make mgr new client
	})

	c.Manager = mgr
	c.Client = client
	// todo: make xscheduler informer factory
	// c.InformerFactory = informers.NewSharedInformerFactory(mgr, 0)
	dynClient := dynamic.NewForConfigOrDie(c.KubeConfig)
	c.DynInformerFactory = dynamicinformer.NewFilteredDynamicSharedInformerFactory(dynClient, 0, corev1.NamespaceAll, nil)
	c.LeaderElection = leaderElectionConfig

	return c, nil
}

// ApplyTo applies the scheduler options to the given scheduler
func (o *Options) ApplyTo(c *xschedulerappconfig.Config) error {
	if len(o.ConfigFile) == 0 {
		// If the --config arg is not specified, honor as well as leader election CLI args.
		o.ApplyLeaderElectionTo(o.ComponentConfig)
		c.ComponentConfig = *o.ComponentConfig
	} else {
		cfg, err := loadConfigFromFile(o.ConfigFile)
		if err != nil {
			return err
		}
		// If the --config arg is specified, honor the leader election CLI args only.
		o.ApplyLeaderElectionTo(cfg)

		// validate configuration from file
		if err := validation.ValidateXschedulerConfiguration(cfg); err != nil {
			return err
		}
		c.ComponentConfig = *cfg
	}
	kubeConfig, err := createKubeConfig(c.ComponentConfig.ClientConnection, o.Master)
	if err != nil {
		return err
	}
	c.KubeConfig = kubeConfig

	if err := o.SecureServing.ApplyTo(&c.SecureServing, &c.LoopbackClientConfig); err != nil {
		return err
	}
	// authentication & authorization
	if o.SecureServing != nil && (o.SecureServing.BindPort != 0 || o.SecureServing.Listener != nil) {
		if err := o.Authentication.ApplyTo(&c.Authentication, c.SecureServing, nil); err != nil {
			return err
		}
		if err := o.Authorization.ApplyTo(&c.Authorization); err != nil {
			return err
		}
	}
	o.Metrics.Apply()

	// flagz allow to view the current values of all command line parameters at runtime
	if utilfeature.DefaultFeatureGate.Enabled(zpagesfeatures.ComponentFlagz) {
		if o.Flags != nil {
			c.Flagz = flagz.NamedFlagSetsReader{FlagSets: *o.Flags}
		}
	}
	return nil
}

func (o *Options) ApplyLeaderElectionTo(cfg *xschedulerconfig.XschedulerConfiguration) {
	// non cli named flags
	if lo.IsNil(o.Flags) {
		return
	}
	leaderElection := o.Flags.FlagSet("leader election")
	// Changed returns true if the flag was explicitly set during Parse() and false otherwise
	if leaderElection.Changed("leader-elect") {
		cfg.LeaderElection.LeaderElect = o.LeaderElection.LeaderElect
	}
	if leaderElection.Changed("leader-elect-lease-duration") {
		cfg.LeaderElection.LeaseDuration = o.LeaderElection.LeaseDuration
	}
	if leaderElection.Changed("leader-elect-renew-deadline") {
		cfg.LeaderElection.RenewDeadline = o.LeaderElection.RenewDeadline
	}
	if leaderElection.Changed("leader-elect-retry-period") {
		cfg.LeaderElection.RetryPeriod = o.LeaderElection.RetryPeriod
	}
	if leaderElection.Changed("leader-elect-resource-lock") {
		cfg.LeaderElection.ResourceLock = o.LeaderElection.ResourceLock
	}
	if leaderElection.Changed("leader-elect-resource-name") {
		cfg.LeaderElection.ResourceName = o.LeaderElection.ResourceName
	}
	if leaderElection.Changed("leader-elect-resource-namespace") {
		cfg.LeaderElection.ResourceNamespace = o.LeaderElection.ResourceNamespace
	}
	o.ComponentConfig = cfg
}

// createKubeConfig creates a kubeConfig from the given config and masterOverride.
func createKubeConfig(config componentbaseconfig.ClientConnectionConfiguration, masterOverride string) (*restclient.Config, error) {
	kubeConfig, err := clientcmd.BuildConfigFromFlags(masterOverride, config.Kubeconfig)
	if err != nil {
		return nil, err
	}
	// DisableCompression bypasses automatic GZip compression requests to the server.
	kubeConfig.DisableCompression = true
	kubeConfig.AcceptContentTypes = config.AcceptContentTypes
	kubeConfig.ContentType = config.ContentType
	kubeConfig.QPS = config.QPS
	kubeConfig.Burst = int(config.Burst)

	return kubeConfig, nil
}

func createClients(kubeConfig *restclient.Config) (clientset.Interface, clientset.Interface, error) {
	client, err := clientset.NewForConfig(restclient.AddUserAgent(kubeConfig, "xscheduler"))
	if err != nil {
		return nil, nil, err
	}
	eventClient, err := clientset.NewForConfig(kubeConfig)
	if err != nil {
		return nil, nil, err
	}
	return client, eventClient, nil
}
