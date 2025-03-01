package reservation

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func GetReservationNamespacedName(reservationRef *corev1.ObjectReference) types.NamespacedName {
	return types.NamespacedName{
		Namespace: reservationRef.Namespace,
		Name:      reservationRef.Name,
	}
}
