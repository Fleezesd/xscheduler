package config

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MigrationControllerArgs holds arguments used to configure the MigrationController
type MigrationControllerArgs struct {
	metav1.TypeMeta

	// DryRun means only execute the entire migration logic except create Reservation or Delete Pod
	// Default is false
	DryRun bool

	// MaxConcurrentReconciles is the maximum number of concurrent Reconciles which can be run. Defaults to 1.
	MaxConcurrentReconciles int32

	// EvictFailedBarePods allows pods without ownerReferences and in failed phase to be evicted.
	EvictFailedBarePods bool

	// EvictLocalStoragePods allows pods using local storage to be evicted.
	EvictLocalStoragePods bool

	// EvictSystemCriticalPods allows eviction of pods of any priority (including Kubernetes system pods)
	EvictSystemCriticalPods bool

	// IgnorePVCPods prevents pods with PVCs from being evicted.
	IgnorePvcPods bool

	// PriorityThreshold filtering only pods under the threshold can be evicted
	PriorityThreshold *PriorityThreshold

	// LabelSelector sets whether to apply label filtering when evicting.
	// Any pod matching the label selector is considered evictable.
	LabelSelector *metav1.LabelSelector

	// Namespaces carries a list of included/excluded namespaces
	Namespaces *Namespaces

	// NodeFit if enabled, it will check whether the candidate Pods have suitable nodes,
	// including NodeAffinity, TaintTolerance, and whether resources are sufficient.
	NodeFit bool

	// NodeSelector for a set of nodes to operate over
	NodeSelector string

	// MaxMigratingGlobally represents the maximum number of pods that can be migrating during migrate globally.
	MaxMigratingGlobally *int32

	// MaxMigratingPerNode represents the maximum number of pods that can be migrating during migrate per node.
	MaxMigratingPerNode *int32

	// MaxMigratingPerNamespace represents the maximum number of pods that can be migrating during migrate per namespace.
	MaxMigratingPerNamespace *int32

	// MaxMigratingPerWorkload represents the maximum number of pods that can be migrating during migrate per workload.
	// Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%).
	MaxMigratingPerWorkload *intstr.IntOrString

	// MaxUnavailablePerWorkload represents he maximum number of pods that can be unavailable during migrate per workload.
	// The unavailable state includes NotRunning/NotReady/Migrating/Evicting
	// Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%).
	MaxUnavailablePerWorkload *intstr.IntOrString

	// SkipCheckExpectedReplicas if enabled, it will allow eviction expectedReplicas equals maxUnavailable or maxMigrating.
	// Default is false
	SkipCheckExpectedReplicas *bool

	// ObjectLimiters control the frequency of migration/eviction to make it smoother,
	// and also protect Pods of the same class from being evicted frequently.
	// e.g. limiting the frequency of Pods of the same workload being evicted.
	// The default is to set the MigrationLimitObjectWorkload limiter.
	ObjectLimiters ObjectLimiterMap

	// DefaultJobMode represents the default operating mode of the PodMigrationJob
	// Default is PodMigrationJobModeReservationFirst
	DefaultJobMode string

	// DefaultJobTTL represents the default TTL of the PodMigrationJob
	// Default is 5 minute
	DefaultJobTTL metav1.Duration

	// EvictQPS controls the number of evict per second
	EvictQPS *Float64OrString
	// EvictBurst is the maximum number of tokens
	EvictBurst int32
	// EvictionPolicy represents how to delete Pod, support "Delete" and "Eviction" and "SoftEviction", default value is "Eviction"
	EvictionPolicy string
	// DefaultDeleteOptions defines options when deleting migrated pods and preempted pods through the method specified by EvictionPolicy
	DefaultDeleteOptions *metav1.DeleteOptions

	// SchedulerNames defines options to assign schedulers that can handle reservation if pmj.mode is ReservationFirst, koord-scheduler by default.
	SchedulerNames []string

	// ArbitrationArgs defines the control parameters of the Arbitration Mechanism.
	ArbitrationArgs *ArbitrationArgs
}

type PriorityThreshold struct {
	Value *int32
	Name  string
}

// Namespaces carries a list of included/excluded namespaces
// for which a given strategy is applicable
type Namespaces struct {
	Include []string
	Exclude []string
}

type MigrationLimitObjectType string

const (
	MigrationLimitObjectWorkload  MigrationLimitObjectType = "workload"
	MigrationLimitObjectNamespace MigrationLimitObjectType = "namespace"
)

type ObjectLimiterMap map[MigrationLimitObjectType]MigrationObjectLimiter

// MigrationObjectLimiter means that if the specified dimension has multiple migrations within the configured time period
// and exceeds the configured threshold, it will be limited.
type MigrationObjectLimiter struct {
	// Duration indicates the time window of the desired limit.
	Duration metav1.Duration
	// MaxMigrating indicates the maximum number of migrations/evictions allowed within the window time.
	// If configured as nil or 0, the maximum number will be calculated according to MaxMigratingPerWorkload.
	MaxMigrating *intstr.IntOrString
	// Burst indicates the limiter allows bursts of up to 'burst' to exceed within the time window.
	Burst int
}

// ArbitrationArgs holds arguments used to configure the Arbitration Mechanism.
type ArbitrationArgs struct {
	// Enabled defines if Arbitration Mechanism should be enabled.
	// Default is true
	Enabled bool

	// Interval defines the running interval (ms) of the Arbitration Mechanism.
	// Default is 500 ms
	Interval *metav1.Duration
}
