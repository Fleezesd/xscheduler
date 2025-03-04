package controllerfinder

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodControllerFinder is a function type that maps a pod to a list of
// controllers and their scale.
type PodControllerFinder func(ref ControllerReference, namespace string) (*ScaleAndSelector, error)

func (r *ControllerFinder) GetPodsForRef(ownerReference *metav1.OwnerReference, namespace string, labelSelector *metav1.LabelSelector, active bool) ([]*corev1.Pod, int32, error) {
	// todo: make pod ref
	return nil, 0, nil
}
