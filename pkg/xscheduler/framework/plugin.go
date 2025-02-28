package framework

import (
	"context"

	corev1 "k8s.io/api/core/v1"
)

type PluginsRunner interface {
	RunXschedulerPlugin(ctx context.Context, node []*corev1.Node) *Status
	RunBalancePlugins(ctx context.Context, node []*corev1.Node) *Status
}

type Status struct {
	Err error
}

type Plugin interface {
	Name() string
}

type XschedulePlugin interface {
	Plugin
	Xscheduler(ctx context.Context, nodes []*corev1.Node) *Status
}

type BalancePlugin interface {
	Plugin
	Balance(ctx context.Context, nodes []*corev1.Node) *Status
}

type FilterPlugin interface {
	Plugin
	// Filter checks if a pod can be evicted
	Filter(pod *corev1.Pod) bool
	// PreEvictionFilter checks if pod can be evicted right before eviction
	PreEvictionFilter(pod *corev1.Pod) bool
}

type EvictPlugin interface {
	Plugin
	// Evict evicts a pod (no pre-check performed)
	Evict(ctx context.Context, pod *corev1.Pod, evictOptions EvictOptions) bool
}
