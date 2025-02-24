package config

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/component-base/config"
)

const (
	DefaultXschedulerPort         = 10258
	DefaultInsecureXschedulerPort = 10251
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DeschedulerConfiguration configures a descheduler
type XschedulerConfiguration struct {
	metav1.TypeMeta

	// LeaderElection defines the configuration of leader election client.
	LeaderElection config.LeaderElectionConfiguration

	// ClientConnection specifies the kubeconfig file and client connection
	// settings for the proxy server to use when communicating with the apiserver.
	ClientConnection config.ClientConnectionConfiguration

	// DebuggingConfiguration holds configuration for Debugging related features.
	config.DebuggingConfiguration

	// HealthzBindAddress is the IP address and port for the health check server to serve on.
	HealthzBindAddress string

	// MetricsBindAddress is the IP address and port for the metrics server to serve on.
	MetricsBindAddress string

	// Time interval for xscheduer to run
	XschedulingInterval metav1.Duration

	// Dry run
	DryRun bool

	// Profiles
	Profiles []XschedulerProfile

	// NodeSelector for a set of nodes to operate over
	NodeSelector *metav1.LabelSelector

	// MaxNoOfPodsToEvictPerNode restricts maximum of pods to be evicted per node.
	MaxNoOfPodsToEvictPerNode *uint

	// MaxNoOfPodsToEvictPerNamespace restricts maximum of pods to be evicted per namespace.
	MaxNoOfPodsToEvictPerNamespace *uint

	// MaxNoOfPodsToTotal restricts maximum of pods to be evicted total.
	MaxNoOfPodsToEvictTotal *uint
}

// XschedulerProfile is a xscheduling profile.
type XschedulerProfile struct {
	Name string `json:"name,omitempty"`
	// todo: add plugin
}
