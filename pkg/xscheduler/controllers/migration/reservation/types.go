package reservation

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	schedulingv1alpha1 "github.com/fleezesd/xscheduler/apis/scheduling/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	DefaultCreator = "xscheduler"
	LabelCreatedBy = "app.kubernetes.io/created-by"
)

var NewInterpreter = newInterpreter

type Interpreter interface {
	GetReservationType() client.Object

	Preemption()

	CreateReservation(ctx context.Context, job *schedulingv1alpha1.PodMigrationJob) (Object, error)

	GetReservation(ctx context.Context, reservationRef *corev1.ObjectReference) (Object, error)

	DeleteReservation(ctx context.Context, reservationRef *corev1.ObjectReference) error
}

type Preemption interface {
	Preempt(ctx context.Context, job *schedulingv1alpha1.PodMigrationJob,
		reservation Object) (bool, reconcile.Result, error)
}

type Object interface {
	metav1.Object
	runtime.Object
	String() string
	OriginObject() client.Object
	GetReservationConditions() []schedulingv1alpha1.ReservationCondition
	QueryPreemptedPodRefs() []corev1.ObjectReference
	GetBoundPod() *corev1.ObjectReference
	GetReservationOwners() []schedulingv1alpha1.ReservationOwner
	GetScheduledNodeName() string
	GetPhase() schedulingv1alpha1.ReservationPhase
	NeedPreemption() bool
}
