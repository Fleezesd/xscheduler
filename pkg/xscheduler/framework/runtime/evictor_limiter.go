package runtime

import (
	corev1 "k8s.io/api/core/v1"
)

type EvictionLimiter interface {
	AllowEvict(pod *corev1.Pod) bool
	Done(pod *corev1.Pod)
	Reset()
	NodeLimitExceeded(node *corev1.Node) bool
	TotalEvicted() uint
}
