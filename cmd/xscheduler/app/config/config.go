package config

import (
	"time"

	xscheduerconfig "github.com/fleezesd/xscheduler/pkg/apis/config"
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
	// 3.todo: Flagz is the Reader interface to get flags for flagz page.
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
	kubeConfig         *restClient.Config
	InformerFactory    informers.SharedInformerFactory
	DynInformerFactory dynamicinformer.DynamicSharedInformerFactory

	//nolint:staticcheck // SA1019 this deprecated field still needs to be used for now. It will be removed once the migration is done.
	EventBroadcaster events.EventBroadcasterAdapter
	LeaderElection   *leaderelection.LeaderElectionConfig

	// PodMaxInUnschedulablePodsDuration is the maximum time a pod can stay in unschedulablePods.
	PodMaxInUnschedulablePodsDuration time.Duration

	// 4. todo: componentGlobalRegistry
}

type completedConfig struct {
	*Config
}

type CompletedConfig struct {
	*completedConfig
}

func (c *Config) Complete() CompletedConfig {
	cc := completedConfig{c}
	return CompletedConfig{&cc}
}
