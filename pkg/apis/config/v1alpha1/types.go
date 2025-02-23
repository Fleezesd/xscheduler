/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/component-base/config/v1alpha1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DeschedulerConfiguration configures a descheduler
type XschedulerConfiguration struct {
	metav1.TypeMeta

	// LeaderElection defines the configuration of leader election client.
	LeaderElection *v1alpha1.LeaderElectionConfiguration `json:"leaderElection,omitempty"`

	// ClientConnection specifies the kubeconfig file and client connection
	// settings for the proxy server to use when communicating with the apiserver.
	ClientConnection v1alpha1.ClientConnectionConfiguration `json:"clientConnection,omitempty"`

	// DebuggingConfiguration holds configuration for Debugging related features.
	v1alpha1.DebuggingConfiguration `json:",inline"`

	// HealthzBindAddress is the IP address and port for the health check server to serve on.
	HealthzBindAddress *string `json:"healthzBindAddress,omitempty"`

	// MetricsBindAddress is the IP address and port for the metrics server to serve on.
	MetricsBindAddress *string `json:"metricsBindAddress,omitempty"`

	// Time interval for xscheduer to run
	XschedulingInterval metav1.Duration `json:"xschedulingInterval,omitempty"`

	// Dry run
	DryRun bool `json:"dryRun,omitempty"`

	// Profiles
	Profiles []XschedulerProfile `json:"profiles,omitempty"`

	// NodeSelector for a set of nodes to operate over
	NodeSelector *metav1.LabelSelector `json:"nodeSelector,omitempty"`

	// MaxNoOfPodsToEvictPerNode restricts maximum of pods to be evicted per node.
	MaxNoOfPodsToEvictPerNode *uint `json:"maxNoOfPodsToEvictPerNode,omitempty"`

	// MaxNoOfPodsToEvictPerNamespace restricts maximum of pods to be evicted per namespace.
	MaxNoOfPodsToEvictPerNamespace *uint `json:"maxNoOfPodsToEvictPerNamespace,omitempty"`

	// MaxNoOfPodsToTotal restricts maximum of pods to be evicted total.
	MaxNoOfPodsToEvictTotal *uint `json:"maxNoOfPodsToEvictTotal,omitempty"`
}

// XschedulerProfile is a xscheduling profile.
type XschedulerProfile struct {
	Name string `json:"name,omitempty"`
	// todo: add plugin
}
