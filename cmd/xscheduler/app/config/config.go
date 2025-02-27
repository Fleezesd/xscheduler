package config

import (
	"time"

	xscheduerconfig "github.com/fleezesd/xscheduler/pkg/xscheduler/apis/config"
	apiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	clientset "k8s.io/client-go/kubernetes"
	restClient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/events"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/component-base/zpages/flagz"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Config has all the context to run a Scheduler
type Config struct {
	// Flagz is the Reader interface to get flags for flagz page.
	Flagz flagz.Reader

	// ComponentConfig is the scheduler server's configuration object.
	ComponentConfig xscheduerconfig.XschedulerConfiguration

	// LoopbackClientConfig is a config for a privileged loopback connection
	LoopbackClientConfig *restClient.Config

	SecureServing  *apiserver.SecureServingInfo
	Authentication apiserver.AuthenticationInfo
	Authorization  apiserver.AuthorizationInfo

	Manager            ctrl.Manager
	Client             clientset.Interface
	KubeConfig         *restClient.Config
	InformerFactory    informers.SharedInformerFactory
	DynInformerFactory dynamicinformer.DynamicSharedInformerFactory

	//nolint:staticcheck // SA1019 this deprecated field still needs to be used for now. It will be removed once the migration is done.
	EventBroadcaster events.EventBroadcasterAdapter
	LeaderElection   *leaderelection.LeaderElectionConfig

	// PodMaxInUnschedulablePodsDuration is the maximum time a pod can stay in unschedulablePods.
	PodMaxInUnschedulablePodsDuration time.Duration

	// todo: componentGlobalRegistry
}

type completedConfig struct {
	*Config
}

// CompletedConfig same as Config, just to swap private object.
type CompletedConfig struct {
	*completedConfig
}

func (c *Config) Complete() CompletedConfig {
	cc := completedConfig{c}

	// AuthorizeClientBearerToken wraps the authenticator and authorizer in loopback authentication logic
	// if the loopback client config is specified AND it has a bearer token. Note that if either authn or
	// authz is nil, this function won't add a token authenticator or authorizer.
	apiserver.AuthorizeClientBearerToken(c.LoopbackClientConfig, &c.Authentication, &c.Authorization)

	return CompletedConfig{&cc}
}
