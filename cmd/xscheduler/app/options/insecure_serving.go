package options

import (
	xschedulerconfig "github.com/fleezesd/xscheduler/pkg/apis/config"
	"github.com/samber/lo"
	"github.com/spf13/pflag"
	apiserveroptions "k8s.io/apiserver/pkg/server/options"
)

// CombinedInsecureServingOptions sets up to two insecure listeners for healthz and metrics. The flags
// override the ComponentConfig and DeprecatedInsecureServingOptions values for both.
type CombinedInsecureServingOptions struct {
	Healthz *apiserveroptions.SecureServingOptions
	Metrics *apiserveroptions.SecureServingOptions

	// overrides
	BindPort    int
	BindAddress string
}

func (o *CombinedInsecureServingOptions) AddFlags(fs *pflag.FlagSet) {
	if lo.IsNil(o) {
		return
	}
	fs.StringVar(&o.BindAddress, "address", "0.0.0.0", "the IP address on which to listen for the --port port (set to 0.0.0.0 or :: for listening in all interfaces and IP families). See --bind-address instead. This parameter is ignored if a config file is specified in --config.")
	fs.IntVar(&o.BindPort, "port", xschedulerconfig.DefaultInsecureXschedulerPort, "the port on which to serve HTTP insecurely without authentication and authorization. If 0, don't serve plain HTTP at all. See --secure-port instead. This parameter is ignored if a config file is specified in --config.")
}
