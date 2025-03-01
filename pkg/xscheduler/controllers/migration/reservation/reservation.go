package reservation

import (
	schedulingv1alpha1 "github.com/fleezesd/xscheduler/apis/scheduling/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ Object = &Reservation{}

type Reservation struct {
	*schedulingv1alpha1.Reservation
}

func NewReservation(reservation *schedulingv1alpha1.Reservation) Object {
	return &Reservation{Reservation: reservation}
}

func (o *Reservation) String() string {
	return o.Reservation.Name
}

func (o *Reservation) OriginObject() client.Object {
	return o.Reservation
}

func (o *Reservation) GetReservationConditions() []schedulingv1alpha1.ReservationCondition {
	return o.Status.Conditions
}

func (o *Reservation) QueryPreemptedPodRefs() []corev1.ObjectReference {
	return nil
}

func (o *Reservation) GetBoundPod() *corev1.ObjectReference {
	if len(o.Status.CurrentOwners) == 0 {
		return nil
	}
	return &o.Status.CurrentOwners[0]
}

func (o *Reservation) GetReservationOwners() []schedulingv1alpha1.ReservationOwner {
	return o.Spec.Owners
}

func (o *Reservation) GetScheduledNodeName() string {
	return o.Status.NodeName
}

func (o *Reservation) GetPhase() schedulingv1alpha1.ReservationPhase {
	return o.Status.Phase
}

func (o *Reservation) NeedPreemption() bool {
	return false
}
